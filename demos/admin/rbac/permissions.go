package rbac

import (
	model "github.com/xhigher/hzgo/demos/admin/model/platform"
	"github.com/xhigher/hzgo/server/admin"
)

func InitPermissions() {

	data, err := model.ReloadRolePermissions()
	if err != nil {
		return
	}
	admin.InitRolePermissions(data)

}
