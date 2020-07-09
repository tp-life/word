package service

import (
	"errors"
	"github.com/gin-gonic/gin"
	"word/internal/entiy"
	"word/pkg/app"
	"word/pkg/middlewares"
)

var (
	ErrLogin = errors.New("login fail")
)

// 生成token
func GenJwt(user entiy.User) (string, error) {
	token, err := middlewares.NewToken(user)
	if err != nil {
		app.Logger().Info("open login error: ", err)
		return "", ErrLogin
	}
	return token, nil
}

// GetUser 获取用户信息
func GetUser(c *gin.Context) entiy.User {
	staff := entiy.User{}
	us, b := c.Get("user")
	if !b {
		return staff
	}
	user, ok := us.(entiy.User)
	if !ok {
		return staff
	}
	return user
}
