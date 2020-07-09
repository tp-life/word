package app

import (
	"github.com/robfig/cron/v3"
	"github.com/sirupsen/logrus"
	"time"
)

// ILogger 用于将 logrus 转为 log.Logger
type ILogger struct {
	*logrus.Logger
}

// Info 实现来自 cron 的logger接口
func (logger *ILogger) Info(msg string, keysAndValues ...interface{}) {
	logger.Logger.WithField("kvs", keysAndValues).WithField("log_type", "cron").Info(msg)
}

// Error 实现来自 cron 的logger接口
func (logger *ILogger) Error(err error, msg string, keysAndValues ...interface{}) {
	logger.Logger.WithField("kvs", keysAndValues).WithField("error", err).WithField("log_type", "cron").Error(msg)
}

// Timer 计划任务
type Timer struct {
	*cron.Cron
}

// NewTimer 开启新的定时器
//  AddJob 和 AddFunc spec 参数说明
//  ┌───────────── minute (0 - 59)
//  │ ┌───────────── hour (0 - 23)
//  │ │ ┌───────────── day of the month (1 - 31)
//  │ │ │ ┌───────────── month (1 - 12)
//  │ │ │ │ ┌───────────── day of the week (0 - 6) (Sunday to Saturday;
//  │ │ │ │ │
//  * * * * * command to execute
//  AddJob 和 AddFunc 的 spec 配置更多信息参阅 wiki
func NewTimer() *Timer {
	return &Timer{cron.New(
		cron.WithLocation(time.Local),
		cron.WithLogger(&ILogger{Logger()}),
		cron.WithChain(cron.Recover(&ILogger{Logger()})),
	)}
}
