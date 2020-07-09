package service

import (
	"errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"word/internal/entiy"
	"word/internal/modules/auth/service"
	"word/internal/modules/user/contracts"
	service2 "word/internal/service"
	"word/pkg/app"
	"word/pkg/database/mongo"
	"word/pkg/password"
)

type User struct {
	ID   string
	base entiy.User
}

func (u User) Info() {
	objID, err := primitive.ObjectIDFromHex(u.ID)
	if err != nil {
		return
	}
	_ = mongo.Collection(u.base).Where(bson.M{"_id": objID}).StFindOne(&u.base)
}

func (u *User) SetID(id string) {
	u.ID = id
}

// Login 登陆
func (User) Login(data contracts.Login) (contracts.LoginSuccess, error) {
	var (
		model entiy.User
	)
	model.FindByAccount(data.Account)
	if model.ID.IsZero() {
		return contracts.LoginSuccess{}, errors.New("inexistence")
	}
	if !password.Verify(data.Password, model.Password) {
		return contracts.LoginSuccess{}, errors.New("password error")
	}
	jwt, err := service2.GenJwt(model)
	if err != nil {
		return contracts.LoginSuccess{}, err
	}
	model.Login()
	rs := contracts.LoginSuccess{
		PerfectInfo: contracts.PerfectInfo{
			Phone:  model.Phone,
			Name:   model.Name,
			Avatar: model.Avatar,
		},
		Account: model.Account,
		Auth:    nil,
		Token:   jwt,
	}

	auths := service.GetUserRoles(model.ID.Hex())
	rs.Auth = auths
	return rs, nil
}

// Register 注册
func (u User) Register(data contracts.Register) (contracts.LoginSuccess, error) {
	var (
		model entiy.User
	)
	model.FindByAccount(data.Account)
	if !model.ID.IsZero() {
		return contracts.LoginSuccess{}, errors.New("username exists")
	}
	user := entiy.User{
		Account:  data.Account,
		Password: password.Hash(data.Password),
		Name:     data.Name,
		Avatar:   data.Avatar,
		Phone:    data.Phone,
		Tags:     []string{app.NORMALUSER},
	}
	_, err := mongo.Collection(user).InsertOneWithError(user)
	rs, err := u.login(user)
	return rs, err
}

func (User) login(m entiy.User) (contracts.LoginSuccess, error){
	jwt, err := service2.GenJwt(m)
	if err != nil {
		return contracts.LoginSuccess{}, err
	}
	rs := contracts.LoginSuccess{
		PerfectInfo: contracts.PerfectInfo{
			Phone:  m.Phone,
			Name:   m.Name,
			Avatar: m.Avatar,
		},
		Account: m.Account,
		Auth:    nil,
		Token:   jwt,
	}
	auths := service.GetUserRoles(m.ID.Hex())
	rs.Auth = auths
	return rs, nil
}

// Perfect 完善资料
func (u User) Perfect(data contracts.PerfectInfo, userID string) error {
	u.SetID(userID)
	u.Info()
	if u.base.ID.IsZero() {
		return errors.New("not found user")
	}
	up := u.base
	up.Phone = data.Phone
	up.Name = data.Name
	up.Avatar = data.Avatar
	_, err := mongo.Collection(u.base).Where(bson.M{"_id": u.base.ID}).UpdateOneWithError(up)
	return err
}
