package model

import (
	"github.com/xhigher/hzgo/demo/model/db/admin"
	"github.com/xhigher/hzgo/logger"
	"github.com/xhigher/hzgo/utils"
	"gorm.io/gorm"
)

func SaveToken(uid, token string, et, it int64) (err error) {
	ts := utils.NowTime()
	updates := map[string]interface{}{
		"token": token,
		"et":    et,
		"ut":ts,
	}
	if it > 0 {
		updates["it"] = it
	}
	res := admin.DB().Model(&admin.StaffTokenModel{}).Where("uid = ?", uid).Updates(updates)
	err = res.Error
	if err != nil {
		logger.Errorf("error: %v", err)
		return
	}
	if res.RowsAffected == 0 {
		data := &admin.StaffTokenModel{
			Uid: uid,
			Token:  token,
			Et:     et,
			It:     it,
			Ut: ts,

		}
		err = admin.DB().Create(data).Error
		if err != nil {
			logger.Errorf("error: %v", err)
			return
		}
	}
	return
}

func CheckToken(uid, token string) (bool, error) {
	data := &admin.StaffTokenModel{}
	err := admin.DB().First(data, "uid = ?", uid).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return false, nil
		}
		logger.Errorf("error: %v", err)
		return false, err
	}
	if data.Token != token {
		return false, nil
	}
	nt := utils.NowTime()
	if data.Et < nt {
		return false, nil
	}
	return true, nil
}
