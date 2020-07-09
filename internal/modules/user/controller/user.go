package controller

import (
	"github.com/gin-gonic/gin"
	"word/internal/modules/user/contracts"
	"word/internal/modules/user/service"
	service2 "word/internal/service"
	"word/pkg/app"
	"word/pkg/validator"
)

// Login 登陆
func Login(c *gin.Context) {
	var (
		rq contracts.Login
	)
	if err := validator.Bind(c, &rq); !err.IsValid() {
		app.F(c, app.Fail, err.String())
		return
	}
	user := service.User{}
	rs, err := user.Login(rq)
	if err != nil {
		app.F(c, app.Fail, err.Error())
		return
	}
	app.OK(c, rs)
}

// Register 注册用户
func Register(c *gin.Context) {
	var (
		rq contracts.Register
	)
	if err := validator.Bind(c, &rq); !err.IsValid() {
		app.F(c, app.Fail, err.String())
		return
	}
	if rq.ConfirmPassword != rq.Password {
		app.F(c, app.Fail, "password are inconsistent")
		return
	}
	user := service.User{}
	rs, err := user.Register(rq)
	if err != nil {
		app.F(c, app.Fail, err.Error())
		return
	}
	app.OK(c, rs)
}

// Perfect 完善用户信息
func Perfect(c *gin.Context) {
	var q contracts.PerfectInfo
	if v := validator.Bind(c, &q); !v.IsValid() {
		app.F(c, app.Fail, v.String())
		return
	}
	user := service2.GetUser(c)
	if user.ID.IsZero() {
		app.F(c, app.AuthFail, app.AuthFailMessage)
		return
	}
	svc := service.User{}
	err := svc.Perfect(q, user.ID.Hex())
	if err != nil {
		app.F(c, app.Fail, err.Error())
		return
	}
	app.OK(c, true)
}
