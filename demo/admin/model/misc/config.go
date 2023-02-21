package misc

import (
	"fmt"
	"github.com/xhigher/hzgo/consts"
	"github.com/xhigher/hzgo/demo/model/db/admin"
	"github.com/xhigher/hzgo/demo/model/db/misc"
	"github.com/xhigher/hzgo/demo/model/db/user"
	"github.com/xhigher/hzgo/logger"
	"github.com/xhigher/hzgo/utils"
	"gorm.io/gorm"
)
func SaveConfigInfo(id, name, items string, static bool, filters string) (err error) {
	tableName := misc.ConfigInfoModel{}.TableName()
	ts := utils.NowTime()
	sql := fmt.Sprintf("INSERT INTO %s (`id`,`name`,`items`,`static`,`filters`,`status`,`ut`) VALUES (?,?,?,?,?,?) ON DUPLICATE KEY UPDATE `name`=?,`items`=?,`static`=?,`filters`=?,`ut`=?", tableName)
	err = admin.DB().Exec(sql, id, name, items, static, filters, consts.StatusEditing, ts, name, items, static, filters, ts).Error
	if err != nil {
		logger.Errorf("error: %v", err)
		return
	}
	return
}

func UpdateConfigStatus(id string, status int32) (err error){
	updates := map[string]interface{}{
		"status": status,
		"ut": utils.NowTime(),
	}
	tx := user.DB().Model(misc.ConfigInfoModel{}).Where("id=?", id)
	if status == consts.StatusOnline {
		tx.Where("status IN (?)", []int32{consts.StatusEditing, consts.StatusOffline})
	}else if status == consts.StatusOffline {
		tx.Where("status=?", consts.StatusOnline)
	}else{
		err = fmt.Errorf("status[%d] error", status)
		return
	}
	tx.Updates(&updates)
	if err != nil {
		logger.Errorf("error: %v", err)
		return
	}
	return
}

func DeleteConfigInfo(id string) (err error) {
	err = admin.DB().Where("id=? and status=?", id, consts.StatusEditing).Delete(&misc.ConfigInfoModel{}).Error
	if err != nil {
		logger.Errorf("error: %v", err)
	}
	return
}

func GetConfigInfo(id string) (data *misc.ConfigInfoModel, err error) {
	err = admin.DB().First(&data, "id = ?", id).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			data = nil
			err = nil
			return
		}
		logger.Errorf("error: %v", err)
	}
	return
}


func GetConfigList(status, offset, limit int32) (total int64, data []*misc.ConfigInfoModel, err error) {
	tx := admin.DB().Model(misc.ConfigInfoModel{}).Where("status = ?", status).Session(&gorm.Session{})
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
