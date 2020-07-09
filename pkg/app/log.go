package app

import (
	"github.com/gin-gonic/gin"
	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	"github.com/medivh-jay/lfshook"
	"github.com/sirupsen/logrus"
	"time"
)

type logConfig struct {
	// 日志保存路径
	Path string `mapstructure:"path" env:"LOG_PATH`
}

var (
	// terminal 默认终端log输出
	terminal = logrus.New()
	logger   = terminal
	config   logConfig
	// 日志文件保留7日
	maxAge = rotatelogs.WithMaxAge(time.Duration(7*86400) * time.Second)
	// 每小时
	rotationTime = rotatelogs.WithRotationTime(time.Hour)
)

// 启动日志服务
func startLogger() {
	err := InitConfig("log", &config)
	if err != nil {
		logger.Fatalln("unable to decode log config, ", err)
	}

	terminal.SetFormatter(&logrus.TextFormatter{ForceColors: true})
	terminal.SetReportCaller(gin.Mode() != gin.ReleaseMode)
	terminal.SetNoLock()

	infoWriter, _ := rotatelogs.New(config.Path+"/info_%Y%m%d.log", maxAge, rotationTime)
	errorWriter, _ := rotatelogs.New(config.Path+"/error_%Y%m%d.log", maxAge, rotationTime)
	debugWriter, _ := rotatelogs.New(config.Path+"/debug_%Y%m%d.log", maxAge, rotationTime)
	warnWriter, _ := rotatelogs.New(config.Path+"/warn_%Y%m%d.log", maxAge, rotationTime)

	if gin.Mode() != gin.ReleaseMode {
		terminal.Level = logrus.DebugLevel
	} else {
		terminal.Level = logrus.WarnLevel
	}

	terminal.AddHook(lfshook.NewHook(
		lfshook.WriterMap{logrus.InfoLevel: infoWriter, logrus.ErrorLevel: errorWriter, logrus.DebugLevel: debugWriter, logrus.WarnLevel: warnWriter}, &logrus.JSONFormatter{}))

}

// Logger 对外日志操作
func Logger() *logrus.Logger {
	return logger
}
