package entiy

import (
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
	"word/pkg/database/mongo"
	"word/pkg/middlewares"
)

type User struct {
	Base     `bson:"#expand"`
	Account  string   `bson:"account"`
	Name     string   `bson:"name"`
	Avatar   string   `bson:"avatar"`
	Password string   `bson:"password"`
	Phone    string   `bson:"phone"`
	Desc     string   `bson:"desc"` // 个人介绍
	Tags     []string `bson:"tags"` // 标签
	LoginAt  int64    `bson:"login_at"`
	LogOutAt int64    `bson:"logout_at"` // 退出登录时间
	Status   uint8    `bson:"status"`    //
}

func (User) TableName() string {
	return "user"
}

// GetTopic 获取用户ID
func (m User) GetTopic() interface{} {
	return m.ID.Hex()
}

// CreateIndex 创建索引
func (m *User) CreateIndex() {
	index := []map[string]interface{}{
		{
			"state": 1,
		},
		{
			"phone": 1,
		},
	}
	mongo.Collection(m).CreateIndexes(index)
}

// CreateUniqueIndexes 创建唯一索引
func (m *User) CreateUniqueIndexes() {
	index := []map[string]interface{}{
		{
			"email": 1,
		},
	}
	mongo.Collection(m).CreateUniqueIndexes(index)
}

// FindByTopic 根据用户ID查找用户
func (m User) FindByTopic(topic interface{}) middlewares.AuthInterface {
	var id primitive.ObjectID
	var err error
	var user User
	if jwtID, ok := topic.(string); ok {
		id, err = primitive.ObjectIDFromHex(jwtID)
		if err != nil {
			return user
		}
	} else {
		if jwtID, ok := topic.(primitive.ObjectID); ok {
			id = jwtID
		}
	}

	_ = mongo.Collection(user).Where(bson.M{"_id": id}).StFindOne(&user)
	return user
}

// GetCheckData 获取 TODO 前期为了方便开发，后期加验证
func (m User) GetCheckData() string {
	return m.ID.Hex()
}

// Check 检测
func (m User) Check(ctx *gin.Context, checkData string) bool {
	if m.LogOutAt > 0 && m.LoginAt <= m.LogOutAt {
		return false
	}
	return true
}

// ExpiredAt 过期时间
func (m User) ExpiredAt() int64 {
	return time.Now().Add(86400 * time.Second).Unix()
}

// Login 记录登陆时间
func (m *User) Login() {
	t := time.Now().Unix()
	m.LoginAt = t
	mongo.Collection(m).Where(bson.M{"_id": m.ID}).UpdateOne(m)
}

// Logout 退出登录
func (m *User) Logout() {
	t := time.Now().Unix()
	m.LogOutAt = t
	mongo.Collection(m).Where(bson.M{"_id": m.ID}).UpdateOne(m)
}

func (m *User) FindByAccount(account string) {
	_ = mongo.Collection(m).Where(bson.M{"account": account}).StFindOne(m)
}
