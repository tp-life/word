package sessions

import (
	"word/pkg/app"
	"github.com/gin-contrib/sessions"
	redisSession "github.com/gin-contrib/sessions/redis"
	"github.com/gin-gonic/gin"
	"strconv"
)

type config struct {
	Key          string `mapstructure:"key" toml:"key" env:"SESSION_KEY`
	Name         string `mapstructure:"name" toml:"name" env:"SESSION_NAME`
	Domain       string `mapstructure:"domain" toml:"domain" env:"SESSION_DOMAIN`
	Addr         string `mapstructure:"addr" toml:"addr" env:"SESSION_URL"`
	Password     string `mapstructure:"password" toml:"password" env:"SESSION_PASS`
	Db           int    `mapstructure:"db" toml:"db" env:"SESSION_DB`
	PoolSize     int    `mapstructure:"pool_size" toml:"pool_size" env:"SESSION_POOL_SIZE`
	MinIdleConns int    `mapstructure:"min_idle_conns" toml:"min_idle_conns" env:"SESSION_MIN_CONN`
}

var conf config

// Inject 启动session服务, 在自定义的路由代码中调用, 传入 *gin.Engine 对象
func Inject(engine *gin.Engine) gin.IRoutes {
	err := app.InitConfig("sessions",&conf)
	if err != nil {
		app.Logger().Fatalln("unable to decode sessions config", err)
	}
	store, err := redisSession.NewStoreWithDB(conf.PoolSize, "tcp", conf.Addr, conf.Password, strconv.Itoa(conf.Db), []byte(conf.Key))
	if err != nil {
		app.Logger().WithField("log_type", "pkg.sessions.sessions").Error(err)
		return engine
	}

	store.Options(sessions.Options{MaxAge: 3600 * 24, Path: "/", Domain: conf.Domain, HttpOnly: true})
	return engine.Use(sessions.Sessions(conf.Name, store))
}

// Get 获取指定session
func Get(c *gin.Context, key string) string {
	sess := sessions.Default(c)
	val := sess.Get(key)
	if val != nil {
		return val.(string)
	}
	return ""
}

// Set 设置session
func Set(c *gin.Context, key, val string) {
	sess := sessions.Default(c)
	sess.Set(key, val)
	_ = sess.Save()
}

// Del 删除指定session
func Del(c *gin.Context, key string) {
	sess := sessions.Default(c)
	sess.Delete(key)
	_ = sess.Save()
}
