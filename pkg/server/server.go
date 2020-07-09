package server

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/http"
	"net/http/httputil"
	"os"
	"strings"
	"time"
	"word/pkg/app"
	"word/pkg/middlewares"
	"word/pkg/validator"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/medivh-jay/daemon"
	"github.com/sirupsen/logrus"
)

// Register 注册服务
func Register(service daemon.Worker) {
	daemon.GetCommand().AddWorker(daemon.NewProcess(service))
}

// Run 运行所有服务
func Run() {
	if rs := daemon.Run(); rs != nil {
		log.Fatalln(rs)
	}
}

// Handler 路由
type Handler struct {
	router *gin.Engine
}

// Register 注册路由
func (handler *Handler) Register(handlers ...func(router gin.IRouter)) {
	for _, f := range handlers {
		f(handler.router)
	}
}

// Group 路由分组
func (*Handler) Group(relativePath string, handlers ...func(router gin.IRouter)) func(router gin.IRouter) {
	return func(router gin.IRouter) {
		var group = router.Group(relativePath)
		for _, f := range handlers {
			f(group)
		}
	}
}

// NewHandler 初始化 handler
func NewHandler() *Handler {
	var handler = &Handler{
		router: gin.New(),
	}
	//handler.router.Use(logger, recovery)
	handler.router.Use(logger)
	//if gin.Mode() != gin.ReleaseMode {
	handler.router.Use(middlewares.CORS)
	//}
	return handler
}

// HTTPServer 统一服务
type HTTPServer struct {
	server       *http.Server
	name         string
	service      func(handler *Handler)
	handler      *Handler
	tasks        func(timer *app.Timer)
	dependencies []func()
	Port         string `mapstructure:"port" env:"PORT" envDefault:"8080"`
}

// PidSavePath pid save path
func (httpServer *HTTPServer) PidSavePath() string {
	if gin.Mode() == gin.ReleaseMode {
		return "/app/data/run/word/"
	}
	return fmt.Sprintf("%s/", app.Root())
}

// Name pid filename , and its command name
func (httpServer *HTTPServer) Name() string {
	return httpServer.name
}

// Start start http server
func (httpServer *HTTPServer) Start() {
	for _, dependency := range httpServer.dependencies {
		dependency()
	}

	// 将 gin 的验证器替换为 v9 版本
	binding.Validator = new(validator.Validator)
	app.InitConfig("", httpServer, func() error {
		err := app.Config().Sub("application").Sub("services").Sub(httpServer.name).Unmarshal(httpServer)
		if err != nil {
			app.Logger().Fatalln("unable to decode application config", err)
		}
		return err
	})

	var h = NewHandler()
	httpServer.service(h)
	httpServer.handler = h

	var timer = app.NewTimer()
	httpServer.tasks(timer)
	go timer.Run()
	defer timer.Stop()

	httpServer.server = &http.Server{Handler: httpServer.handler.router, Addr: ":" + httpServer.Port}
	err := httpServer.server.ListenAndServe()
	if err == http.ErrServerClosed {
		app.Logger().Info(fmt.Sprintf("service [%s] closed at %s", httpServer.name, time.Now().Format(time.RFC3339)))
	} else {
		app.Logger().Error(fmt.Sprintf("service [%s] error: %v", httpServer.name, err))
	}
}

// Stop before stop http server operation, usually gracefully shut down the server
func (httpServer *HTTPServer) Stop() error {
	if httpServer.server == nil {
		return nil
	}
	err := httpServer.server.Shutdown(context.Background())
	app.Logger().Info(fmt.Sprintf("service [%s] stoped", httpServer.name))
	return err
}

// Restart restart http server, usually gracefully shut down the server and then execute start
func (httpServer *HTTPServer) Restart() error {
	app.Logger().Info(fmt.Sprintf("service [%s] restarting", httpServer.name))
	return httpServer.Stop()
}

// NewService 服务
func NewService(name string, tasks func(timer *app.Timer), handler func(handler *Handler), dependencies ...func()) *HTTPServer {
	return &HTTPServer{name: name, service: handler, tasks: tasks, dependencies: dependencies}
}

// 自定义的GIN日志处理中间件
// 可能在终端输出时看起来比较的难看
func logger(ctx *gin.Context) {
	start := time.Now()
	path := ctx.Request.URL.Path
	raw := ctx.Request.URL.RawQuery

	ctx.Next()

	if raw != "" {
		path = path + "?" + raw
	}

	var params = make(logrus.Fields)
	params["latency"] = time.Now().Sub(start)
	params["url"] = ctx.Request.URL.String()
	params["method"] = ctx.Request.Method
	params["status"] = ctx.Writer.Status()
	params["body_size"] = ctx.Writer.Size()
	params["client_ip"] = ctx.ClientIP()
	params["user_agent"] = ctx.Request.UserAgent()
	params["log_type"] = "pkg.server.server"
	params["keys"] = ctx.Keys
	app.Logger().WithFields(params).Info()
}

func recovery(ctx *gin.Context) {
	defer func() {
		if err := recover(); err != nil {
			var brokenPipe bool
			if ne, ok := err.(*net.OpError); ok {
				if se, ok := ne.Err.(*os.SyscallError); ok {
					if strings.Contains(strings.ToLower(se.Error()), "broken pipe") || strings.Contains(strings.ToLower(se.Error()), "connection reset by peer") {
						brokenPipe = true
					}
				}
			}
			stack := app.Stack(3)
			httpRequest, _ := httputil.DumpRequest(ctx.Request, false)

			if gin.Mode() != gin.ReleaseMode {
				app.Logger().WithField("log_type", "pkg.server.server").Error(string(httpRequest))
				var errors = make([]logrus.Fields, 0)
				for i := 0; i < len(stack); i++ {
					errors = append(errors, logrus.Fields{
						"func":   stack[i]["func"],
						"source": stack[i]["source"],
						"file":   fmt.Sprintf("%s:%d", stack[i]["file"], stack[i]["line"]),
					})
				}
				app.Logger().WithField("log_type", "pkg.server.server").WithField("stack", errors).Error(err)
				ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"stack": errors, "message": err})
			} else {
				app.Logger().WithField("log_type", "pkg.server.server").
					WithField("stack", stack).WithField("request", string(httpRequest)).Error()
			}

			if brokenPipe {
				_ = ctx.Error(err.(error))
				ctx.Abort()
			} else {
				ctx.AbortWithStatus(http.StatusInternalServerError)
			}
		}
	}()
	ctx.Next()
}
