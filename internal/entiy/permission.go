package entiy

import (
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"word/pkg/database/mongo"
)

// Permission 权限列表
type Permission struct {
	Base   `bson:"#expand"`
	Name   string `json:"name" bson:"name" binding:"max=24" form:"name"` // 权限名称
	Value  string `bson:"value"`                                         // 权限值
	Path   string `json:"path" bson:"path" form:"path"`                  // 资源定位路径
	Method string `json:"method" bson:"method" form:"method"`            // 请求方式
}

// Permissions 一个权限列表
type Permissions []Permission

// HasPermission 查询权限列表中是否包含指定的权限
func (permissions Permissions) HasPermission(id primitive.ObjectID, path, method string) bool {
	for _, permission := range permissions {
		// 超级管理员
		if permission.Method == "*" && permission.Path == "*" {
			return true
		}
		if permission.ID.Hex() == id.Hex() {
			return true
		}
	}
	return false
}

// TableName 表名
func (Permission) TableName() string { return "permissions" }

// GetPermissionsByIDs 根据权限ID获取权限
func (m Permission) GetPermissionsByIDs(ids ...primitive.ObjectID) Permissions {
	var (
		permissions Permissions
	)
	if len(ids) == 0 {
		return permissions
	}
	_ = mongo.Collection(m).Where(bson.M{"_id": bson.M{"$in": ids}}).StFindMany(&permissions)
	return permissions
}

// GetPermissionsByRequest 根据请求参数获取权限列表
//  uri 资源路径
//  method 请求方式
func (m Permission) GetPermissionsByRequest(path string, method string) Permissions {
	var permissions Permissions
	_ = mongo.Collection(m).Where(bson.M{"path": path, "method": method}).StFindMany(&permissions)
	return permissions
}
