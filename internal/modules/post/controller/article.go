package controller

import (
	"github.com/gin-gonic/gin"
	"word/internal/modules/post/contracts"
	"word/internal/modules/post/service"
	"word/pkg/app"
	"word/pkg/validator"
)

// ArticleLists 文章列表
func ArticleLists(c *gin.Context) {
	var q contracts.PostSearch
	if err := validator.Bind(c, &q); !err.IsValid() {
		app.F(c, app.Fail, err.String())
		return
	}

	post := service.NewPost(nil, "")
	lists, count := post.Lists(q, app.Offset(c), app.Limit(c))
	app.OriOK(c, map[string]interface{}{
		"count": count,
		"data":  lists,
	})
}

// Detail 详情
func ArticleDetail(c *gin.Context) {
	id := c.Param("post_id")
	if id == "" {
		app.F(c, app.Fail, app.ParamError)
		return
	}
	post := service.NewPost(nil, id)
	info, err := post.Detail()
	if err != nil {
		app.F(c, app.Fail, err.Error())
		return
	}
	app.OriOK(c, info)
}
