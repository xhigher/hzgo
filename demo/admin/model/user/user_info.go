package model

import (
	"fmt"
	"github.com/xhigher/hzgo/consts"
	"github.com/xhigher/hzgo/demo/model/db/user"
	"github.com/xhigher/hzgo/logger"
	"github.com/xhigher/hzgo/utils"
	"gorm.io/gorm"
)

func GetUserInfo(userid string) (data *user.UserInfoModel, err error) {
	err = user.DB().First(&data, "userid = ?", userid).Error
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

func GetUserList(status, offset, limit int32) (total int64, data []*user.UserInfoModel, err error) {
	tx := user.DB().Model(user.UserInfoModel{}).Where("status = ?", status).Session(&gorm.Session{})
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


func UpdateUserStatus(userid string, status int32) (err error){
	updates := map[string]interface{}{
		"status": status,
		"ut": utils.NowTime(),
	}
	tx := user.DB().Model(user.UserInfoModel{}).Where("userid = ?", userid)
	if status == consts.UserStatusActive {
		tx.Where("status=?", consts.UserStatusBlocked)
	}else if status == consts.UserStatusBlocked {
		tx.Where("status=?", consts.UserStatusActive)
	}else if status == consts.UserStatusCancelled {
		tx.Where("status=?", consts.UserStatusActive)
	}else{
		err = fmt.Errorf("status[%d] error", status)
		return
	}
	err = tx.Updates(&updates).Error
	if err != nil {
		logger.Errorf("error: %v", err)
		return
	}
	return
}

func CleanUserTokens(userid string) (err error){
	updates := map[string]interface{}{
		"token": "",
		"et":0,
		"it":0,
		"ut": utils.NowTime(),
	}
	err = user.DB().Model(user.UserTokenModel{}).Where("userid = ?", userid).Updates(&updates).Error
	if err != nil {
		logger.Errorf("error: %v", err)
		return
	}
	return
}

