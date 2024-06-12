package platform

import (
	"github.com/xhigher/hzgo/bizerr"
	model "github.com/xhigher/hzgo/demos/admin/model/platform"
	"github.com/xhigher/hzgo/demos/model/db/admin"
	"github.com/xhigher/hzgo/types"
)

func GetRolesMenuList(roles []string) (data []*admin.RoleMenusModel, be *bizerr.Error) {
	rids := types.StringArray(roles)
	data, err := model.GetRolesMenuList(rids)
	if err != nil {
		be = bizerr.New(err)
		return
	}
	return
}

func GetRoleList(status int32) (data []*admin.RoleInfoModel, be *bizerr.Error) {
	data, err := model.GetRoleList(status)
	if err != nil {
		be = bizerr.New(err)
		return
	}
	return
}

func GetMenuList(status int32) (data []*admin.MenuInfoModel, be *bizerr.Error) {
	data, err := model.GetMenuList(status)
	if err != nil {
		be = bizerr.New(err)
		return
	}
	return
}

func GetRolePermissionList(role string, offset, limit int32) (total int64, data []*admin.RolePermissionsModel, be *bizerr.Error) {
	total, data, err := model.GetRolePermissionList(role, offset, limit)
	if err != nil {
		be = bizerr.New(err)
		return
	}
	return
}

func GetRoleMenuList(role string, offset, limit int32) (total int64, data []*admin.RoleMenusModel, be *bizerr.Error) {
	total, data, err := model.GetRoleMenuList(role, offset, limit)
	if err != nil {
		be = bizerr.New(err)
		return
	}
	return
}

func GetAllRolesPermissions() (data map[string]map[string]bool, err error) {
	tempData, err := model.GetAllRolesPermissions()
	if err != nil {
		return
	}
	data = map[string]map[string]bool{}
	for _, item := range tempData {
		if _, ok := data[item.Rid]; !ok {
			data[item.Rid] = map[string]bool{}
		}
		data[item.Rid][item.Path] = true
	}
	return
}
