package model

import (
	"github.com/xhigher/hzgo/consts"
	"github.com/xhigher/hzgo/demos/model/db/admin"
	"github.com/xhigher/hzgo/logger"
	"github.com/xhigher/hzgo/mysql"
	"github.com/xhigher/hzgo/types"
	"github.com/xhigher/hzgo/utils"
	"gorm.io/gorm"
)

func GetRoleList(status int32) (data []*admin.RoleInfoModel, err error) {
	err = admin.DB().Where("status = ?", status).Find(&data).Error
	if err != nil {
		logger.Errorf("error: %v", err)
		return
	}
	return
}

func GetRoleInfo(rid string) (data *admin.RoleInfoModel, err error) {
	err = admin.DB().Where("rid = ?", rid).First(&data).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			data = nil
			err = nil
			return
		}
		logger.Errorf("error: %v", err)
		return
	}
	return
}

func AddRole(rid, name string) (existed bool, err error) {
	model := &admin.RoleInfoModel{
		Rid:    rid,
		Name:   name,
		Status: consts.StatusOnline,
		Ut:     utils.NowTime(),
	}
	err = admin.DB().Create(model).Error
	if err != nil {
		logger.Errorf("error: %v", err)
		if mysql.ErrDuplicateKey(err) {
			existed = true
			return
		}
	}
	return
}

func UpdateRole(rid string, name string) (err error) {
	updates := map[string]interface{}{
		"name": name,
		"ut":   utils.NowTime(),
	}
	err = admin.DB().Model(admin.RoleInfoModel{}).Where("rid = ?", rid).Updates(&updates).Error
	if err != nil {
		logger.Errorf("error: %v", err)
		return
	}
	return
}

func UpdateRoleStatus(rid string, status int32) (err error) {
	updates := map[string]interface{}{
		"status": status,
		"ut":     utils.NowTime(),
	}
	tx := admin.DB().Model(admin.RoleInfoModel{}).Where("rid = ?", rid)
	if status == consts.StatusOnline {
		tx.Where("status=?", consts.StatusOffline)
	} else {
		tx.Where("status=?", consts.StatusOnline)
	}
	err = tx.Updates(&updates).Error
	if err != nil {
		logger.Errorf("error: %v", err)
		return
	}
	return
}

func GetRolePermissionList(rid string, offset, limit int32) (total int64, data []*admin.RolePermissionsModel, err error) {
	tx := admin.DB().Model(admin.RolePermissionsModel{})
	if rid != "" {
		tx.Where("rid=?", rid)
	}
	tx = tx.Session(&gorm.Session{})
	err = tx.Count(&total).Error
	if err != nil {
		logger.Errorf("error: %v", err)
		return
	}
	if total == 0 {
		return
	}

	err = tx.Offset(int(offset)).Limit(int(limit)).Find(&data).Error
	if err != nil {
		logger.Errorf("error: %v", err)
		return
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

func RemoveRolePermission(rid string, path string) (err error) {
	err = admin.DB().Where("rid=? AND path=?", rid, path).Delete(admin.RolePermissionsModel{}).Error
	if err != nil {
		logger.Errorf("error: %v", err)
	}
	return
}

func GetMenuList(status int32) (data []*admin.MenuInfoModel, err error) {
	err = admin.DB().Where("status = ?", status).Find(&data).Error
	if err != nil {
		logger.Errorf("error: %v", err)
		return
	}
	return
}

func GetMenuInfo(rid string) (data *admin.MenuInfoModel, err error) {
	err = admin.DB().Where("rid = ?", rid).First(&data).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			data = nil
			err = nil
			return
		}
		logger.Errorf("error: %v", err)
		return
	}
	return
}

func AddMenu(name, icon, path string, upMid int) (existed bool, err error) {
	model := &admin.MenuInfoModel{
		Name:   name,
		Icon:   icon,
		Path:   path,
		UpMid:  upMid,
		Status: consts.StatusOnline,
		Ut:     utils.NowTime(),
	}
	err = admin.DB().Create(model).Error
	if err != nil {
		logger.Errorf("error: %v", err)
		if mysql.ErrDuplicateKey(err) {
			existed = true
			return
		}
	}
	return
}

func UpdateMenu(mid int, name, icon, path string, upMid int) (err error) {
	updates := map[string]interface{}{
		"name":   name,
		"icon":   icon,
		"path":   path,
		"up_mid": upMid,
		"ut":     utils.NowTime(),
	}
	err = admin.DB().Model(admin.MenuInfoModel{}).Where("mid = ?", mid).Updates(&updates).Error
	if err != nil {
		logger.Errorf("error: %v", err)
		return
	}
	return
}

func UpdateMenuStatus(mid string, status int32) (err error) {
	updates := map[string]interface{}{
		"status": status,
		"ut":     utils.NowTime(),
	}
	tx := admin.DB().Model(admin.MenuInfoModel{}).Where("mid = ?", mid)
	if status == consts.StatusOnline {
		tx.Where("status=?", consts.StatusOffline)
	} else {
		tx.Where("status=?", consts.StatusOnline)
	}
	err = tx.Updates(&updates).Error
	if err != nil {
		logger.Errorf("error: %v", err)
		return
	}
	return
}

func GetRoleMenuList(rid string, offset, limit int32) (total int64, data []*admin.RoleMenusModel, err error) {
	tx := admin.DB().Model(admin.RoleMenusModel{})
	if rid != "" {
		tx.Where("rid=?", rid)
	}
	tx = tx.Session(&gorm.Session{})
	err = tx.Count(&total).Error
	if err != nil {
		logger.Errorf("error: %v", err)
		return
	}
	if total == 0 {
		return
	}

	err = tx.Offset(int(offset)).Limit(int(limit)).Find(&data).Error
	if err != nil {
		logger.Errorf("error: %v", err)
		return
	}
	return
}

func GetRolesMenuList(rids types.StringArray) (data []*admin.RoleMenusModel, err error) {
	tx := admin.DB().Model(admin.RoleMenusModel{})
	if len(rids) > 0 {
		tx.Where("rid IN (?)", rids)
	}
	err = tx.Find(&data).Error
	if err != nil {
		logger.Errorf("error: %v", err)
		return
	}
	return
}

func AddRoleMenu(rid string, mid, upMid int, path string) (existed bool, err error) {
	model := &admin.RoleMenusModel{
		Rid:   rid,
		Mid:   mid,
		UpMid: upMid,
		Path:  path,
		Ut:    utils.NowTime(),
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
