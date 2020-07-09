package password

import (
	"golang.org/x/crypto/bcrypt"
	"word/pkg/app"
)

type passConfig struct {
	PassToken string `mapstructure:"password_token" env:"PASS_TOKEN`
}

var (
	passConf      passConfig
	passwordToken = loadPassConfig()
)

func loadPassConfig() string {
	app.InitConfig("application", &passConf)
	return passConf.PassToken
}

// Hash 密码hash
func Hash(password string) string {
	bytes, err := bcrypt.GenerateFromPassword([]byte(passwordToken+password), bcrypt.DefaultCost)
	if err != nil {
		app.Logger().WithField("log_type", "pkg.password.password").Error(err)
		return ""
	}

	return string(bytes)
}

// Verify 密码hash验证
func Verify(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(passwordToken+password))
	return err == nil
}
