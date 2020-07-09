package service

import (
	"word/internal/entiy"
)

// GetUserRoles 获取用户具体权限值
func GetUserRoles(userID string) []string {
	var (
		gr      entiy.Binding
		pm      entiy.Permission
		permits = make([]string, 0)
	)
	roleID := gr.GetRoleIDs(userID)
	rs := pm.GetPermissionsByIDs(roleID...)
	for _, v := range rs {
		permits = append(permits, v.Value)
	}
	return permits
}
