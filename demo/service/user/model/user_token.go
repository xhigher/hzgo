package model

import (
	"github.com/xhigher/hzgo/demo/model/db/user"
	"github.com/xhigher/hzgo/logger"
	"github.com/xhigher/hzgo/utils"
	"gorm.io/gorm"
)

func SaveToken(userid, token string, et, it int64) (err error) {
	ts := utils.NowTime()
	updates := map[string]interface{}{
		"token": token,
		"et":    et,
		"ut": ts,
	}
	if it > 0 {
		updates["it"] = it
	}
	res := user.DB().Model(&user.UserTokenModel{}).Where("userid = ?", userid).Updates(updates)
	err = res.Error
	if err != nil {
		logger.Errorf("SaveToken update error: %v", err)
		return
	}
	if res.RowsAffected == 0 {
		data := &user.UserTokenModel{
			Userid: userid,
			Token:  token,
			Et:     et,
			It:     it,
			Ut: ts,
		}
		err = user.DB().Create(data).Error
		if err != nil {
			logger.Errorf("SaveToken create error: %v", err)
			return
		}
	}
	return
}

func CheckToken(userid, token string) (bool, error) {
	data := &user.UserTokenModel{}
	err := user.DB().Where("userid = ?", userid).First(data).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return false, nil
		}
		logger.Errorf("CheckToken error: %v", err)
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
