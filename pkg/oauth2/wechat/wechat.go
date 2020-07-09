package wechat

import (
	"encoding/json"
	"word/pkg/oauth2/internal"
	jsoniter "github.com/json-iterator/go"
	"math"
	"time"
)

type tokenJSON struct {
	AccessToken  string         `json:"access_token"`
	TokenType    string         `json:"token_type"`
	RefreshToken string         `json:"refresh_token"`
	ExpiresIn    expirationTime `json:"expires_in"`
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

var _ = internal.SetTokenParser("api.weixin.qq.com", func(body string) *internal.Token {
	var tj tokenJSON
	if err := jsoniter.UnmarshalFromString(body, &tj); err != nil {
		return nil
	}
	token := &internal.Token{
		AccessToken:  tj.AccessToken,
		TokenType:    tj.TokenType,
		RefreshToken: tj.RefreshToken,
		Expiry:       tj.expiry(),
		Raw:          make(map[string]interface{}),
	}
	_ = jsoniter.UnmarshalFromString(body, &token.Raw)
	return token
})
