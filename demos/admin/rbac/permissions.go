package rbac

import "github.com/xhigher/hzgo/server/admin"

func InitPermissions() {

	admin.InitRolePermissions(admin.RoleMaintainer, map[string]admin.CRUD{
		Staff:       admin.CRUDAll,
		Dashboard:   admin.CRUDReadonly,
		BizUser:     admin.CRUDAll,
		BizConfig:   admin.CRUDAll,
		BizStat:     admin.CRUDAll,
		BizActivity: admin.CRUDAll,
	})

	admin.InitRolePermissions(admin.RoleDeveloper, map[string]admin.CRUD{
		Dashboard:   admin.CRUDReadonly,
		BizUser:     admin.CRUDAll,
		BizConfig:   admin.CRUDAll,
		BizStat:     admin.CRUDAll,
		BizActivity: admin.CRUDAll,
	})

	admin.InitRolePermissions(admin.RoleOperator, map[string]admin.CRUD{
		Dashboard:   admin.CRUDReadonly,
		BizUser:     admin.CRUDAll,
		BizConfig:   admin.CRUDAll,
		BizStat:     admin.CRUDAll,
		BizActivity: admin.CRUDAll,
	})
}
