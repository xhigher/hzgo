package misc

import (
	"fmt"
	"github.com/xhigher/hzgo/consts"
	"github.com/xhigher/hzgo/demos/model/db/admin"
	"github.com/xhigher/hzgo/demos/model/db/misc"
	"github.com/xhigher/hzgo/demos/model/db/user"
	"github.com/xhigher/hzgo/logger"
	"github.com/xhigher/hzgo/utils"
	"gorm.io/gorm"
)

func SaveBannerInfo(id int32, site string, typ int32, name, img, data string) (err error) {
	tableName := misc.BannerInfoModel{}.TableName()
	ts := utils.NowTime()
	sn := ts
	sql := fmt.Sprintf("INSERT INTO %s (`id`,`site`,`type`,`name`,`img`,`data`,`sn`,`status`,`ut`) VALUES (?,?,?,?,?,?) ON DUPLICATE KEY UPDATE `type`=?,`name`=?,`img`=?,`data`=?,`ut`=?", tableName)
	err = admin.DB().Exec(sql, id, site, typ, name, img, data, sn, consts.StatusEditing, ts, typ, name, img, data, ts).Error
	if err != nil {
		logger.Errorf("error: %v", err)
		return
	}
	return
}

func UpdateBannerStatus(id int32, status int32) (err error) {
	updates := map[string]interface{}{
		"status": status,
		"ut":     utils.NowTime(),
	}
	tx := user.DB().Model(misc.BannerInfoModel{}).Where("id=?", id)
	if status == consts.StatusOnline {
		tx.Where("status IN (?)", []int32{consts.StatusEditing, consts.StatusOffline})
	} else if status == consts.StatusOffline {
		tx.Where("status=?", consts.StatusOnline)
	} else {
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

func DeleteBannerInfo(id int32) (err error) {
	err = admin.DB().Where("id=? and status=?", id, consts.StatusEditing).Delete(&misc.BannerInfoModel{}).Error
	if err != nil {
		logger.Errorf("error: %v", err)
	}
	return
}

func GetBannerInfo(id int32) (data *misc.BannerInfoModel, err error) {
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

func GetBannerList(site string, status, offset, limit int32) (total int64, data []*misc.BannerInfoModel, err error) {
	tx := admin.DB().Model(misc.BannerInfoModel{}).Where("status = ?", status)
	if len(site) > 0 {
		tx.Where("site=?", site)
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
