// Package rbac 权限控制
package middlewares

import (
	"github.com/gin-gonic/gin"
	"word/internal/entiy"
)

// HasPermission 判断指定用户是否有当前访问资源的权限
func HasPermission(userID string, ctx *gin.Context) bool {
	var (
		binding    = entiy.Binding{}
		role       = entiy.RoleGroup{}
		permission = entiy.Permission{}
	)
	bindings := binding.GetRoleIDs(userID)

	if len(bindings) == 0 {
		return false
	}
	permissionIDs := role.GetRoles(bindings...).GetPermissionIDs()
	if len(permissionIDs) == 0 {
		return false
	}
	permissions := permission.GetPermissionsByIDs(permissionIDs...)

	for _, id := range permissionIDs {
		if permissions.HasPermission(id, ctx.Request.URL.Path, ctx.Request.Method) {
			return true
		}
	}

	return false
}
