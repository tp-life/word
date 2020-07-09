package captcha

import (
	"fmt"
	"image/color"
	"sync"
	"time"
	"word/pkg/app"

	"github.com/go-redis/redis"
	"github.com/mojocn/base64Captcha"
)

var (
	captcha = newCaptcha()
	store   = new(customizeRdsStore)
	conf    config
	source  = "23456789qwertyuipkjhgfdsazxcvbnm"
)

func newCaptcha() *base64Captcha.Captcha {
	//var driver = base64Captcha.NewDriverDigit(80, 240, 6, 0.7, 80)
	var driver = base64Captcha.NewDriverString(
		80,
		240,
		30,
		base64Captcha.OptionShowHollowLine|base64Captcha.OptionShowSineLine,
		6,
		source,
		&color.RGBA{R: 254, G: 254, B: 254, A: 254},
		[]string{"wqy-microhei.ttc"})
	return base64Captcha.NewCaptcha(driver, store)
}

type (
	customizeRdsStore struct {
		redisClient redis.UniversalClient
		sync.Once
	}

	config struct {
		Addr         string `mapstructure:"addr" toml:"addr" env:"CAPTCHA_URL"`
		Password     string `mapstructure:"password" toml:"password" env:"CAPTCHA_PASS"`
		Db           int    `mapstructure:"db" toml:"db" env:"CAPTCHA_DB"`
		PoolSize     int    `mapstructure:"pool_size" toml:"pool_size" env:"CAPTCHA_POOL_SIZE"`
		MinIdleConns int    `mapstructure:"min_idle_conns" toml:"min_idle_conns" env:"CAPTCHA_MIN_CONN"`
	}
)

func (s *customizeRdsStore) Verify(id, answer string, clear bool) bool {
	return s.Get(id, clear) == answer
}

// GenerateWithDriver 自定义生成驱动
//  content 为随机使用的字节集合
func GenerateWithDriver(driver base64Captcha.Driver) (id, b64s string, err error) {
	return base64Captcha.NewCaptcha(driver, store).Generate()
}

// Generate 生成验证码
func Generate() (id, b64s string, err error) {
	return captcha.Generate()
}

// Verify 验证验证码是否有效
func Verify(id, value string) bool {
	return captcha.Verify(id, value, true)
}

// GetCaptchaValue 获取验证码内容
func GetCaptchaValue(id string) string {
	return captcha.Store.Get(id, false)
}

func (s *customizeRdsStore) lazyLoad() {
	s.Once.Do(func() {
		_ = app.InitConfig("captcha", &conf)
		store.redisClient = redis.NewUniversalClient(&redis.UniversalOptions{
			Addrs:        []string{conf.Addr},
			Password:     conf.Password,
			DB:           conf.Db,
			PoolSize:     conf.PoolSize,
			MinIdleConns: conf.MinIdleConns,
		})
	})
}

func (s *customizeRdsStore) Set(id string, value string) {
	s.lazyLoad()
	err := s.redisClient.Set(fmt.Sprintf("catpcha_id_%s", id), value, time.Minute*10).Err()
	if err != nil {
		app.Logger().WithField("log_type", "pkg.Captcha.captcha").Error(err)
	}
}

func (s *customizeRdsStore) Get(id string, clear bool) string {
	s.lazyLoad()
	val, err := s.redisClient.Get(fmt.Sprintf("catpcha_id_%s", id)).Result()
	if err != nil {
		app.Logger().WithField("log_type", "pkg.Captcha.captcha").Error(err)
		return ""
	}
	if clear {
		err := s.redisClient.Del(fmt.Sprintf("catpcha_id_%s", id)).Err()
		if err != nil {
			app.Logger().WithField("log_type", "pkg.Captcha.captcha").Error(err)
			return ""
		}
	}
	return val
}
