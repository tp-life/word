package router

import "github.com/gin-gonic/gin"

type auth struct{}

func (auth) Router(router gin.IRouter, mid ...gin.HandlerFunc) {
	//a := router.Group("auth", mid...)
	//a.GET("")
}
