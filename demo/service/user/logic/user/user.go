package user

import (
	"github.com/xhigher/hzgo/consts"
	usermodel "github.com/xhigher/hzgo/demo/service/user/model/user"
)

func CheckUser(username, password string) (userid string, err error) {
	userInfo, err := usermodel.GetUser(username)
	if err != nil {
		return
	}
	if userInfo == nil {
		err = consts.ErrUserNotExisted
		return
	}
	if !usermodel.CheckPassword(userInfo, password) {
		err = consts.ErrUserPasswordWrong
		return
	}
	userid = userInfo.Userid
	return
}
