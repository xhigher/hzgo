package user

import (
	"github.com/xhigher/hzgo/utils"
	"gorm.io/gorm"
)

func SaveToken(userid, token string, et, it int64) (err error) {
	updates := map[string]interface{}{
		"token": token,
		"et":     et,
		"it": it,
	}
	res := DB().Where("userid = ?", userid).Updates(updates)
	err = res.Error
	if err != nil {
		return
	}
	if res.RowsAffected == 0 {
		data := &UserToken{
			Userid: userid,
			Token: token,
			Et: et,
			It: it,
		}
		err = DB().Create(data).Error
		if err != nil {
			return
		}
	}
	return
}

func CheckToken(userid, token string) (bool, error) {
	data := &UserToken{}
	err := DB().Where("userid = ?", userid).First(data).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return false, nil
		}
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
