package model

import (
	"github.com/xhigher/hzgo/demo/model/db/user"
	"github.com/xhigher/hzgo/utils"
	"gorm.io/gorm"
)

func GetUserInfo(userid string) (data *user.UserInfoModel, err error) {
	err = user.DB().Where("userid = ?", userid).First(&data).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			data = nil
			err = nil
			return
		}
	}
	return
}

func GetUserList(status, offset, limit int32) (total int64, data []*user.UserInfoModel, err error) {
	err = user.DB().Model(user.UserInfoModel{}).Where("status = ?", status).Count(&total).Error
	if err != nil {
		return
	}
	if total == 0 {
		return
	}

	err = user.DB().Where("status = ?", status).Offset(int(offset)).Limit(int(limit)).First(&data).Error
	if err != nil {
		return
	}
	return
}


func UpdateUserStatus(userid string, status int32) (err error){
	updates := map[string]interface{}{
		"status": status,
		"ut": utils.NowTime(),
	}
	err = user.DB().Model(user.UserInfoModel{}).Where("userid = ?", userid).Updates(&updates).Error
	if err != nil {
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
		return
	}
	return
}

