package entiy

import (
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"word/pkg/database/mongo"
)

// Binding 用户id与权限id的关系
type Binding struct {
	Base       `bson:"#expand"`
	UserID     string   `bson:"user_id" `
	Permission []string ` bson:"permission"`
}

// TableName 表名
func (Binding) TableName() string { return "binding" }

// GetRoleIDs 获取指定用户的所有角色ID
func (m Binding) GetRoleIDs(userID string) []primitive.ObjectID {
	var bindings []Binding
	_ = mongo.Collection(m).Where(bson.M{"user_id": userID}).StFindMany(&bindings)

	var roles []primitive.ObjectID
	for _, bind := range bindings {
		for _, v := range bind.Permission {
			id, _ := primitive.ObjectIDFromHex(v)
			roles = append(roles, id)
		}
	}
	return roles
}
