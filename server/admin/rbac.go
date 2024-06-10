package admin

import (
	"github.com/xhigher/hzgo/types"
)

const (
	RoleMaintainer = "maintainer" // 超管
	RoleDeveloper  = "developer"  // 开发

)

var rolePermissions = map[string]map[string]bool{}

func InitRolePermissions(data map[string]map[string]bool) {
	rolePermissions = data
}

func CheckPermission(rid string, path string) bool {
	if ps, ok := rolePermissions[rid]; ok {
		return ps[path]
	}
	return false
}

func CheckRoles(roles types.StringArray) bool {
	if len(roles) > 0 {
		for _, r := range roles {
			if _, ok := rolePermissions[r]; !ok {
				return false
			}
		}
	}
	return true
}
