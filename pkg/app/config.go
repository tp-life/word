package app

import (
	"fmt"
	"github.com/caarlos0/env/v6"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"path/filepath"
	"strings"
	"sync"
)

var (
	loadConfig sync.Once
	baseConfig *Configuration
)

// Configuration 配置驱动
type Configuration struct {
	// 配置文件路径
	configPath string
	config     *viper.Viper
}

//
func InitConfig(sub string, target interface{}, handler ...func() error ) (err error) {
	if len(handler) == 0  && sub != ""{
		if err = Config().Sub(sub).Unmarshal(target); err != nil {
			logger.Error("获取环境变量失败")
		}

	} else {
		handler[0]()
	}
	if err = env.Parse(target); err != nil {
		logger.Error("获取环境变量失败")
	}
	return err
}

// Config 配置读取
func Config() *viper.Viper {
	loadConfig.Do(func() {
		baseConfig = new(Configuration)
		baseConfig.configPath = fmt.Sprintf("%s/configs/%s", Root(), gin.Mode())
		baseConfig.config = viper.New()
		baseConfig.config.AddConfigPath(baseConfig.configPath)
		baseConfig.config.SetConfigType("yaml")
		baseConfig.config.SetConfigName("application")
		err := baseConfig.config.ReadInConfig()
		if err != nil {
			Logger().Fatalln(err)
		}
	})
	return baseConfig.config
}

// NewConfig 加载其他配置文件, 你能指定一个单一的文件名, 必须包含后缀, 将会解析后缀为文件类型
func NewConfig(filename string) *viper.Viper {
	var ext = filepath.Ext(filename)
	var configuration *Configuration
	configuration = new(Configuration)
	configuration.configPath = fmt.Sprintf("%s/configs/%s", Root(), gin.Mode())
	configuration.config = viper.New()
	configuration.config.AddConfigPath(configuration.configPath)
	configuration.config.SetConfigType(ext[1:])
	configuration.config.SetConfigName(strings.Replace(filename, ext, "", -1))
	err := configuration.config.ReadInConfig()
	if err != nil {
		Logger().Fatalln(err)
	}
	return configuration.config
}
