package qq

import (
	"encoding/json"
	"word/pkg/app"
	"word/pkg/oauth2/internal"
	"math"
	"net/url"
	"strconv"
	"time"
)

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

var _ = internal.SetTokenParser("graph.qq.com", func(body string) *internal.Token {
	result, err := url.ParseQuery(body)
	if err != nil {
		app.Logger().WithField("body", body).Warn(err)
		return nil
	}
	var expiry, _ = strconv.Atoi(result.Get("expires_in"))
	token := &internal.Token{
		AccessToken:  result.Get("access_token"),
		TokenType:    result.Get("token_type"),
		RefreshToken: result.Get("refresh_token"),
		Expiry:       time.Now().Add(time.Duration(expiry) * time.Second),
	}

	e := result.Get("expires_in")
	if e == "" {
		e = result.Get("expires")
	}
	expires, _ := strconv.Atoi(e)
	if expires != 0 {
		token.Expiry = time.Now().Add(time.Duration(expires) * time.Second)
	}

	if token.AccessToken == "" {
		return token
	}
	return token
})
