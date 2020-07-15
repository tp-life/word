package controller

import (
	"github.com/gin-gonic/gin"
	"word/internal/modules/post/contracts"
	"word/internal/modules/post/service"
	"word/pkg/app"
	"word/pkg/validator"
)

// Lists 列表
func Lists(c *gin.Context) {
	var q contracts.PostSearch
	if err := validator.Bind(c, &q); !err.IsValid() {
		app.F(c, app.Fail, err.String())
		return
	}
	post := service.NewPost(c, "")
	lists, count := post.Lists(q, app.Offset(c), app.Limit(c))
	app.OriOK(c, map[string]interface{}{
		"count": count,
		"data":  lists,
	})
}

// Detail 详情
func Detail(c *gin.Context) {
	id := c.Param("post_id")
	if id == "" {
		app.F(c, app.Fail, app.ParamError)
		return
	}
	post := service.NewPost(c, id)
	info, err := post.Detail()
	if err != nil {
		app.F(c, app.Fail, err.Error())
		return
	}
	app.OriOK(c, info)
}

// Created 新增
func Created(c *gin.Context) {
	var q contracts.CreatePost
	if err := validator.Bind(c, &q); !err.IsValid() {
		app.F(c, app.Fail, err.String())
		return
	}
	post := service.NewPost(c, "")
	id, err := post.CreatePost(q)
	if err != nil {
		app.F(c, app.Fail, err.Error())
		return
	}
	app.OK(c, id)
}

// Update 编辑
func Update(c *gin.Context) {
	var q contracts.CreatePost
	if err := validator.Bind(c, &q); !err.IsValid() {
		app.F(c, app.Fail, err.String())
		return
	}
	id := c.Param("post_id")
	if id == "" {
		app.F(c, app.Fail, app.ParamError)
		return
	}
	post := service.NewPost(c, id)
	err := post.UpdatePost(q)
	if err != nil {
		app.F(c, app.Fail, err.Error())
		return
	}
	app.OK(c, true)
}

// Delete 删除
func Delete(c *gin.Context) {
	id := c.Param("post_id")
	if id == "" {
		app.F(c, app.Fail, app.ParamError)
		return
	}
	post := service.NewPost(c, id)
	err := post.Delete()
	if err != nil {
		app.F(c, app.Fail, err.Error())
		return
	}
	app.OK(c, true)
}

// ChangeStatus 状态变更
func ChangeStatus(c *gin.Context) {
	var q contracts.PostStatus
	if err := validator.Bind(c, &q); !err.IsValid() {
		app.F(c, app.Fail, err.String())
		return
	}
	id := c.Param("post_id")
	if id == "" {
		app.F(c, app.Fail, app.ParamError)
		return
	}
	post := service.NewPost(c, id)
	err := post.UpdateStatus(q)
	if err != nil {
		app.F(c, app.Fail, err.Error())
		return
	}
	app.OK(c, true)
}

// JoinProject 添加专题
func JoinProject(c *gin.Context) {
	var q contracts.JoinProject
	if err := validator.Bind(c, &q); !err.IsValid() {
		app.F(c, app.Fail, err.String())
		return
	}
	id := c.Param("post_id")
	if id == "" {
		app.F(c, app.Fail, app.ParamError)
		return
	}
	post := service.NewPost(c, id)
	err := post.JoinProject(id, q.Project)
	if err != nil {
		app.F(c, app.Fail, err.Error())
		return
	}
	app.OK(c, true)
}
