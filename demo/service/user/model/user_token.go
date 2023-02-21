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
	model := &user.UserTokenModel{}
	res := user.DB().Model(model).Where("userid = ?", userid).Updates(updates)
	err = res.Error
	if err != nil {
		logger.Errorf("error: %v", err)
		return
	}
	if res.RowsAffected == 0 {
		model = &user.UserTokenModel{
			Userid: userid,
			Token:  token,
			Et:     et,
			It:     it,
			Ut: ts,
		}
		err = user.DB().Create(model).Error
		if err != nil {
			logger.Errorf("error: %v", err)
			return
		}
	}
	return
}

func CheckToken(userid, token string) (ok bool, err error) {
	data := &user.UserTokenModel{}
	err = user.DB().First(data,"userid = ?", userid).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			err = nil
			return
		}
		logger.Errorf("error: %v", err)
		return
	}
	if data.Token != token {
		return
	}
	nt := utils.NowTime()
	if data.Et < nt {
		return
	}
	ok = true
	return
}
