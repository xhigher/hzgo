package admin

import (
	"github.com/xhigher/hzgo/logger"
	"github.com/xhigher/hzgo/types"
)

const (
	RoleMaintainer = "maintainer" // 超管
	RoleDeveloper = "developer" // 开发
	RoleOperator = "operator"  // 运营
	RoleTreasurer = "treasurer" // 财务
	RoleCustomer = "customer"  // 客服
)

type CRUDType int

const (
	CRUDCreate CRUDType = 1
	CRUDDelete CRUDType = 2
	CRUDUpdate CRUDType = 3
	CRUDRead CRUDType = 4
	CRUDWrite CRUDType = 5
)

type CRUD struct {
	create bool `json:"create"`
	delete bool `json:"delete"`
	update bool `json:"update"`
	read bool `json:"read"`
	write bool `json:"write"` //contains create,delete,update
}

func (m CRUD) Create() bool {
	return m.create
}
func (m CRUD) Delete() bool {
	return m.delete
}
func (m CRUD) Update() bool {
	return m.update
}
func (m CRUD) Read() bool {
	return m.read
}

func (m CRUD) Write() bool {
	return m.write
}

var (
	CRUDAll = CRUD{
		create: true,
		delete: true,
		update: true,
		read: true,
		write: true,
	}

	CRUDReadonly = CRUD{
		create: false,
		delete: false,
		update: false,
		read: true,
		write: false,
	}

)


var rolePermissions = map[string]map[string]CRUD {
	RoleMaintainer: {},
	RoleDeveloper: {},
	RoleOperator: {},
	RoleTreasurer: {},
	RoleCustomer: {},
}

func InitRolePermissions(role string, permissions map[string]CRUD){
	rolePermissions[role] = permissions
}

func getModulePermission(role, module string) CRUD {
	if ps, ok := rolePermissions[role]; ok {
		logger.Infof("getModulePermission: %v", ps)
		return ps[module]
	}
	return CRUD{}
}

func CheckRoles(roles types.StringArray) bool {
	if len(roles) > 0 {
		for _, r := range roles {
			if _,ok := rolePermissions[r]; !ok {
				return false
			}
		}
	}
	return true
}