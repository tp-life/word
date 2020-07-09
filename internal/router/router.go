package router

import (
	"github.com/gin-gonic/gin"
	"word/internal/entiy"
	"word/pkg/middlewares"
)

type IRouter interface {
	Router(router gin.IRouter, handlerFunc ...gin.HandlerFunc)
}

func Service(router gin.IRouter) {
	route := []IRouter{
		auth{},
		user{},
		post{},
		project{},
	}
	for _, v := range route {
		v.Router(router, middlewares.VerifyAuth("user", "JWT", entiy.User{}))
	}
}
