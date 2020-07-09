package api

import (
	"fmt"
	user "word/internal/router"

	// "word/pkg/database/managers"
	"word/pkg/middlewares"
	"word/pkg/server"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

var (
	swagHandler gin.HandlerFunc
)

func init() {
	swagHandler = ginSwagger.WrapHandler(swaggerFiles.Handler)
}

func swagService(router gin.IRouter) {
	if gin.Mode() != gin.ReleaseMode {
		router.GET(fmt.Sprintf("/api/doc/*any"), middlewares.CORS, swagHandler)
	}
}

func Service(handler *server.Handler) {
	// var admin = managers.New()
	// admin.Register(entity.Pay{}, managers.Mysql)
	handler.Register(
		swagService,
		user.Service,
	)

}
