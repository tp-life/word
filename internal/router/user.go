package router

import (
	"github.com/gin-gonic/gin"
	"word/internal/modules/user/controller"
)

type user struct {
}

func (user) Router(router gin.IRouter, handlerFunc ...gin.HandlerFunc) {
	router.POST("register", controller.Register) // 注册用户
	router.POST("login", controller.Login)       // 登陆
	g := router.Group("user", handlerFunc...)
	g.PUT("info", controller.Perfect) // 完善用户信息
}
