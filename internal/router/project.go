package router

import (
	"github.com/gin-gonic/gin"
	"word/internal/modules/project/controller"
)

type project struct {
}

func (project) Router(router gin.IRouter, handlerFunc ...gin.HandlerFunc) {
	project := router.Group("project", handlerFunc...)
	project.GET("", controller.Lists) // 列表
	project.POST("", controller.Create)                       // 新增
	project.PUT("up/:id", controller.Update)                  // 编辑
	project.PUT("post/:id/allow", controller.AllowPost)       // 通过文章加入专题
	project.PUT("post/:id/refuse", controller.RefusePost)     // 拒绝文章加入专题
	project.PUT("state/:id/allow", controller.AllowProject)   // 启用专题状态
	project.PUT("state/:id/refuse", controller.RefuseProject) // 禁用专题状态

}
