package router

import (
	"github.com/gin-gonic/gin"
	"word/internal/modules/post/controller"
)

type post struct {
}

func (post) Router(router gin.IRouter, handleFun ...gin.HandlerFunc) {
	p := router.Group("post", handleFun...)
	p.POST("", controller.Created)                    // 添加文章
	p.PUT(":post_id/st", controller.ChangeStatus)     // 更改状态
	p.PUT(":post_id/up", controller.Update)           // 修改内容
	p.DELETE(":post_id", controller.Delete)           // 删除
	p.GET("", controller.Lists)                       // 列表
	p.GET(":post_id", controller.Detail)              // 详情
	p.GET(":post_id/project", controller.JoinProject) // 加入专题

	article := router.Group("article")
	article.GET("", controller.ArticleLists)          // 文章列表
	article.GET(":post_id", controller.ArticleDetail) // 文章详情

}
