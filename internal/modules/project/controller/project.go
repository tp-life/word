package controller

import (
	"github.com/gin-gonic/gin"
	"word/internal/entiy"
	"word/internal/modules/project/contracts"
	"word/internal/modules/project/service"
	"word/pkg/app"
	"word/pkg/validator"
)

// Lists 专题列表
func Lists(c *gin.Context) {
	project := service.Project{}
	rs := project.Projects()
	app.OK(c, rs)
}

// Create 创建专题
func Create(c *gin.Context) {
	var q contracts.CreateProject
	if v := validator.Bind(c, &q); !v.IsValid() {
		app.F(c, app.Fail, v.String())
		return
	}
	project := service.Project{}
	err := project.Create(q)
	if err != nil {
		app.F(c, app.Fail, err.Error())
		return
	}
	app.OK(c, true)
}

// Update 修改专题
func Update(c *gin.Context) {
	var q contracts.CreateProject
	if v := validator.Bind(c, &q); !v.IsValid() {
		app.F(c, app.Fail, v.String())
		return
	}
	projectID := c.Param("id")
	if projectID == "" {
		app.F(c, app.Fail, "not found this project")
		return
	}
	project := service.Project{}
	err := project.Update(q, projectID)
	if err != nil {
		app.F(c, app.Fail, err.Error())
		return
	}
	app.OK(c, true)
}

// AllowProject 通过专题
func AllowProject(c *gin.Context) {
	project := service.Project{}
	projectID := c.Param("id")
	if projectID == "" {
		app.F(c, app.Fail, "not found this project")
		return
	}
	err := project.ChangeState(projectID, entiy.ENABLE)
	if err != nil {
		app.F(c, app.Fail, err.Error())
		return
	}
	app.OK(c, true)
}

// RefuseProject 禁用专题
func RefuseProject(c *gin.Context) {
	project := service.Project{}
	projectID := c.Param("id")
	if projectID == "" {
		app.F(c, app.Fail, "not found this project")
		return
	}
	err := project.ChangeState(projectID, entiy.Disable)
	if err != nil {
		app.F(c, app.Fail, err.Error())
		return
	}
	app.OK(c, true)
}

// AllowPost 通过文章加入专题
func AllowPost(c *gin.Context) {
	project := service.ProjectPostM{}
	projectID := c.Param("id")
	if projectID == "" {
		app.F(c, app.Fail, "not found this project")
		return
	}
	err := project.AuditPost(projectID, entiy.ENABLE)
	if err != nil {
		app.F(c, app.Fail, err.Error())
		return
	}
	app.OK(c, true)
}

// RefusePost 拒绝文章加入专题
func RefusePost(c *gin.Context) {
	project := service.ProjectPostM{}
	projectID := c.Param("id")
	if projectID == "" {
		app.F(c, app.Fail, "not found this project")
		return
	}
	err := project.AuditPost(projectID, entiy.Disable)
	if err != nil {
		app.F(c, app.Fail, err.Error())
		return
	}
	app.OK(c, true)
}
