package sms

import (
	"errors"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/dysmsapi"
	jsoniter "github.com/json-iterator/go"
	"sync"
	"word/pkg/app"
)

var (
	single sync.Once
	//config struct {
	//	AppKey   string `mapstructure:"app_key" toml:"app_key"`
	//	SmsAPI   string `mapstructure:"sms_api" toml:"sms_api"`
	//	User     string `mapstructure:"user" toml:"sms_api"`
	//	Password string `mapstructure:"password" toml:"password"`
	//}
	aliConfig struct {
		AppKey    string `mapstructure:"app_key"    toml:"app_key" env:"SMS_APP_KEY"`
		AppSecret string `mapstructure:"app_secret" toml:"app_secret" env:"SMS_APP_SECRET"`
		Sign      string `mapstructure:"sign" toml:"sign" env:"SMS_SIGN`
		Area      string `mapstructure:"area" toml:"area" env:"SMS_AREA`
	}
)

// Send 发送短信
//func Send(ctx *gin.Context, mobile, content string) error {
//	single.Do(func() {
//		_ = app.Config().Sub("sms").Unmarshal(&config)
//	})
//	var (
//		formData = url.Values{
//			"userId":      {config.User},
//			"md5password": {config.Password},
//			"content":     {content},
//			"mobile":      {mobile},
//		}
//		request, _ = http.NewRequestWithContext(ctx, http.MethodPost, config.SmsAPI, strings.NewReader(formData.Encode()))
//	)
//	request.Header.Add("content-type", "application/x-www-form-urlencoded")
//	resp, err := http.DefaultClient.Do(request)
//	if err != nil {
//		app.Logger().Debug(resp)
//		app.Logger().WithField("log_type", "pkg.sms.sms").Warn(err)
//		return err
//	}
//
//	return nil
//}

// Send 发送短信
// mobile 手机号
// code  验证码
// msgType 消息模板编号
func Send(mobile, code, msgType string) error {
	single.Do(func() {
		_ = app.InitConfig("ali_sms", &aliConfig)
	})
	var (
		accessKey    = aliConfig.AppKey
		accessSecret = aliConfig.AppSecret
		sign         = aliConfig.Sign
		area         = aliConfig.Area
	)
	client, err := dysmsapi.NewClientWithAccessKey(area, accessKey, accessSecret)
	if err != nil {
		return err
	}
	request := dysmsapi.CreateSendSmsRequest()
	request.Scheme = "https"
	request.PhoneNumbers = mobile
	request.SignName = sign
	request.TemplateCode = msgType
	codeStr, _ := jsoniter.Marshal(map[string]string{"code": code})
	request.TemplateParam = string(codeStr)
	res, err := client.SendSms(request)
	if err != nil {
		return err
	}
	if res.Code == "OK" {
		return nil
	} else {
		return errors.New(res.Message)
	}
}
