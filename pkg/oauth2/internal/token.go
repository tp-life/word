// Copyright 2014 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package internal

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"math"
	"mime"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"sync"
	"time"

	"golang.org/x/net/context/ctxhttp"
)

type (
	// TokenHost token地址
	TokenHost string
	// TokenParser 解析方法
	TokenParser func(body string) *Token
)

// TokenResponseParser 指定TOKEN解析方法
var TokenResponseParser = make(map[TokenHost]TokenParser)

// SetTokenParser 指定TOKEN解析方法
func SetTokenParser(uri TokenHost, parser TokenParser) map[TokenHost]TokenParser {
	TokenResponseParser[uri] = parser
	return TokenResponseParser
}

// Token represents the credentials used to authorize
// the requests to access protected resources on the OAuth 2.0
// provider's backend.
//
// This type is a mirror of oauth2.Token and exists to break
// an otherwise-circular dependency. Other internal packages
// should convert this Token into an oauth2.Token before use.
type Token struct {
	// AccessToken is the token that authorizes and authenticates
	// the requests.
	AccessToken string

	// TokenType is the type of token.
	// The Type method returns either this or "Bearer", the default.
	TokenType string

	// RefreshToken is a token that's used by the application
	// (as opposed to the user) to refresh the access token
	// if it expires.
	RefreshToken string

	// Expiry is the optional expiration time of the access token.
	//
	// If zero, TokenSource implementations will reuse the same
	// token forever and RefreshToken or equivalent
	// mechanisms for that TokenSource will not be used.
	Expiry time.Time

	// Raw optionally contains extra metadata from the server
	// when updating a token.
	Raw interface{}
}

// tokenJSON is the struct representing the HTTP response from OAuth2
// providers returning a token in JSON form.
type tokenJSON struct {
	AccessToken  string         `json:"access_token"`
	TokenType    string         `json:"token_type"`
	RefreshToken string         `json:"refresh_token"`
	ExpiresIn    expirationTime `json:"expires_in"` // at least PayPal returns string, while most return number
}

func (e *tokenJSON) expiry() (t time.Time) {
	if v := e.ExpiresIn; v != 0 {
		return time.Now().Add(time.Duration(v) * time.Second)
	}
	return
}

type expirationTime int32

func (e *expirationTime) UnmarshalJSON(b []byte) error {
	if len(b) == 0 || string(b) == "null" {
		return nil
	}
	var n json.Number
	err := json.Unmarshal(b, &n)
	if err != nil {
		return err
	}
	i, err := n.Int64()
	if err != nil {
		return err
	}
	if i > math.MaxInt32 {
		i = math.MaxInt32
	}
	*e = expirationTime(i)
	return nil
}

// AuthStyle 认证类型
type AuthStyle int

const (
	// AuthStyleUnknown 默认的
	AuthStyleUnknown AuthStyle = 0
	// AuthStyleInParams 放入url参数中
	AuthStyleInParams AuthStyle = 1
	// AuthStyleInHeader 放入 Header 中
	AuthStyleInHeader AuthStyle = 2
)

// authStyleCache is the set of tokenURLs we've successfully used via
// RetrieveToken and which style auth we ended up using.
// It's called a cache, but it doesn't (yet?) shrink. It's expected that
// the set of OAuth2 servers a program contacts over time is fixed and
// small.
var authStyleCache struct {
	sync.Mutex
	m map[string]AuthStyle // keyed by tokenURL
}

// ResetAuthCache resets the global authentication style cache used
// for AuthStyleUnknown token requests.
func ResetAuthCache() {
	authStyleCache.Lock()
	defer authStyleCache.Unlock()
	authStyleCache.m = nil
}

// lookupAuthStyle reports which auth style we last used with tokenURL
// when calling RetrieveToken and whether we have ever done so.
func lookupAuthStyle(tokenURL string) (style AuthStyle, ok bool) {
	authStyleCache.Lock()
	defer authStyleCache.Unlock()
	style, ok = authStyleCache.m[tokenURL]
	return
}

// setAuthStyle adds an entry to authStyleCache, documented above.
func setAuthStyle(tokenURL string, v AuthStyle) {
	authStyleCache.Lock()
	defer authStyleCache.Unlock()
	if authStyleCache.m == nil {
		authStyleCache.m = make(map[string]AuthStyle)
	}
	authStyleCache.m[tokenURL] = v
}

// newTokenRequest returns a new *http.Request to retrieve a new token
// from tokenURL using the provided clientID, clientSecret, and POST
// body parameters.
//
// inParams is whether the clientID & clientSecret should be encoded
// as the POST body. An 'inParams' value of true means to send it in
// the POST body (along with any values in v); false means to send it
// in the Authorization header.
func newTokenRequest(tokenURL, clientID, clientSecret string, v url.Values, authStyle AuthStyle) (*http.Request, error) {
	if authStyle == AuthStyleInParams {
		v = cloneURLValues(v)
		if clientID != "" {
			v.Set("client_id", clientID)
		}
		if clientSecret != "" {
			v.Set("client_secret", clientSecret)
		}
	}
	req, err := http.NewRequest("POST", tokenURL, strings.NewReader(v.Encode()))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	if authStyle == AuthStyleInHeader {
		req.SetBasicAuth(url.QueryEscape(clientID), url.QueryEscape(clientSecret))
	}
	return req, nil
}

func cloneURLValues(v url.Values) url.Values {
	v2 := make(url.Values, len(v))
	for k, vv := range v {
		v2[k] = append([]string(nil), vv...)
	}
	return v2
}

// RetrieveToken 获取 token
func RetrieveToken(ctx context.Context, clientID, clientSecret, tokenURL string, v url.Values, authStyle AuthStyle) (*Token, error) {
	needsAuthStyleProbe := authStyle == 0
	if needsAuthStyleProbe {
		if style, ok := lookupAuthStyle(tokenURL); ok {
			authStyle = style
			needsAuthStyleProbe = false
		} else {
			authStyle = AuthStyleInHeader
		}
	}
	req, err := newTokenRequest(tokenURL, clientID, clientSecret, v, authStyle)
	if err != nil {
		return nil, err
	}
	token, err := doTokenRoundTrip(ctx, req)
	if err != nil && needsAuthStyleProbe {
		authStyle = AuthStyleInParams
		req, _ = newTokenRequest(tokenURL, clientID, clientSecret, v, authStyle)
		token, err = doTokenRoundTrip(ctx, req)
	}
	if needsAuthStyleProbe && err == nil {
		setAuthStyle(tokenURL, authStyle)
	}

	if token != nil && token.RefreshToken == "" {
		token.RefreshToken = v.Get("refresh_token")
	}
	return token, err
}

func doTokenRoundTrip(ctx context.Context, req *http.Request) (*Token, error) {
	r, err := ctxhttp.Do(ctx, ContextClient(ctx), req)
	if err != nil {
		return nil, err
	}
	body, err := ioutil.ReadAll(io.LimitReader(r.Body, 1<<20))
	r.Body.Close()
	if err != nil {
		return nil, fmt.Errorf("oauth2: cannot fetch token: %v", err)
	}
	if code := r.StatusCode; code < 200 || code > 299 {
		return nil, &RetrieveError{
			Response: r,
			Body:     body,
		}
	}

	var token *Token

	if parser, ok := TokenResponseParser[TokenHost(req.Host)]; ok {
		token = parser(string(body))
	} else {
		content, _, _ := mime.ParseMediaType(r.Header.Get("Content-Type"))
		switch content {
		case "application/x-www-form-urlencoded", "text/plain":
			vals, err := url.ParseQuery(string(body))
			if err != nil {
				return nil, err
			}
			token = &Token{
				AccessToken:  vals.Get("access_token"),
				TokenType:    vals.Get("token_type"),
				RefreshToken: vals.Get("refresh_token"),
				Raw:          vals,
			}
			e := vals.Get("expires_in")
			expires, _ := strconv.Atoi(e)
			if expires != 0 {
				token.Expiry = time.Now().Add(time.Duration(expires) * time.Second)
			}
		default:
			var tj tokenJSON
			if err = json.Unmarshal(body, &tj); err != nil {
				return nil, err
			}
			token = &Token{
				AccessToken:  tj.AccessToken,
				TokenType:    tj.TokenType,
				RefreshToken: tj.RefreshToken,
				Expiry:       tj.expiry(),
				Raw:          make(map[string]interface{}),
			}
			_ = json.Unmarshal(body, &token.Raw) // no error checks for optional fields
		}
	}

	if token.AccessToken == "" {
		return nil, errors.New("oauth2: server response missing access_token")
	}
	return token, nil
}

// RetrieveError 令牌获取的错误信息
type RetrieveError struct {
	Response *http.Response
	Body     []byte
}

func (r *RetrieveError) Error() string {
	return fmt.Sprintf("oauth2: cannot fetch token: %v\nResponse: %s", r.Response.Status, r.Body)
}
