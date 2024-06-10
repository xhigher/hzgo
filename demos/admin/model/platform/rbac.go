package model

import (
	"github.com/xhigher/hzgo/consts"
	"github.com/xhigher/hzgo/demos/model/db/admin"
	"github.com/xhigher/hzgo/logger"
	"github.com/xhigher/hzgo/mysql"
	"github.com/xhigher/hzgo/utils"
)

func GetRoleList() (data []*admin.RoleInfoModel, err error) {
	err = admin.DB().Model(admin.RoleInfoModel{}).Where("status = ?", consts.StatusOnline).Find(&data).Error
	if err != nil {
		logger.Errorf("error: %v", err)
		return
	}
	return
}

func AddRole(name string) (existed bool, err error) {
	model := &admin.RoleInfoModel{
		Name:   name,
		Status: consts.UserStatusActive,
		Ut:     utils.NowTime(),
	}
	err = admin.DB().Create(model).Error
	if err != nil {
		if mysql.ErrDuplicateKey(err) {
			existed = true
		} else {
			logger.Errorf("error: %v", err)
		}
	}
	return
}

func AddRolePermission(rid string, path string) (existed bool, err error) {
	model := &admin.RolePermissionsModel{
		Rid:  rid,
		Path: path,
		Ut:   utils.NowTime(),
	}
	err = admin.DB().Create(model).Error
	if err != nil {
		if mysql.ErrDuplicateKey(err) {
			existed = true
		} else {
			logger.Errorf("error: %v", err)
		}
	}
	return
}

func AddRoleMenu(rid string, mid int, path string) (existed bool, err error) {
	model := &admin.RoleMenusModel{
		Rid:  rid,
		Mid:  mid,
		Path: path,
		Ut:   utils.NowTime(),
	}
	err = admin.DB().Create(model).Error
	if err != nil {
		if mysql.ErrDuplicateKey(err) {
			existed = true
		} else {
			logger.Errorf("error: %v", err)
		}
	}
	return
}

func RemoveRoleMenu(rid string, mid int) (err error) {
	err = admin.DB().Where("rid=? AND mid=?", rid, mid).Delete(admin.RoleMenusModel{}).Error
	if err != nil {
		logger.Errorf("error: %v", err)
	}
	return
}

func RemoveMenu(rid string, mid int) (err error) {
	err = admin.DB().Where("rid=? AND mid=?", rid, mid).Delete(admin.RoleMenusModel{}).Error
	if err != nil {
		logger.Errorf("error: %v", err)
	}
	return
}

func ReloadRolePermissions() (data map[string]map[string]bool, err error) {
	var tempData []*admin.RolePermissionsModel
	err = admin.DB().Model(admin.RolePermissionsModel{}).Find(&tempData).Error
	if err != nil {
		logger.Errorf("error: %v", err)
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
