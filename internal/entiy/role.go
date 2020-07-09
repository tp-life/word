package entiy

import (
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"word/pkg/database/mongo"
)

// RoleGroup 权限分组 用于展示作用
type RoleGroup struct {
	Base        `bson:"#expand"`
	Name        string   `json:"name" bson:"name" form:"name" binding:"max=12"`       // 权限分组名称
	Permissions []string `json:"permissions" bson:"permissions" form:"permissions[]"` // 权限id列表
}

// Roles 角色列表
type Roles []RoleGroup

// GetPermissionIDs 获取角色列表中包含的所有权限id
func (roles Roles) GetPermissionIDs() []primitive.ObjectID {
	var ids []primitive.ObjectID
	for _, role := range roles {
		for _, permission := range role.Permissions {
			id, _ := primitive.ObjectIDFromHex(permission)
			ids = append(ids, id)
		}
	}
	return ids
}

// TableName 表名
func (RoleGroup) TableName() string { return "roles" }

// GetRoles 获取权限
func (m RoleGroup) GetRoles(id ...primitive.ObjectID) Roles {
	var (
		roles Roles
		where = bson.M{}
	)
	if len(id) > 0 {
		where = bson.M{"_id": bson.M{"$in": id}}
	}
	mongo.Collection(m).Where(where).FindMany(&roles)
	return roles
}
