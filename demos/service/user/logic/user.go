package logic

import (
	"github.com/xhigher/hzgo/bizerr"
	"github.com/xhigher/hzgo/consts"
	"github.com/xhigher/hzgo/demo/model/db/user"
	"github.com/xhigher/hzgo/demo/service/user/model"
)

func CheckUser(username, password string) (userid string, be *bizerr.Error) {
	userInfo, err := model.GetUser(username)
	if err != nil {
		be = bizerr.New(err)
		return
	}
	if userInfo == nil {
		be = bizerr.UserNull()
		return
	}
	if userInfo.Status == consts.UserStatusBlocked {
		be = bizerr.UserBlocked()
		return
	}
	if userInfo.Status == consts.UserStatusCancelled {
		be = bizerr.UserCanceled()
		return
	}

	if !model.CheckPassword(userInfo, password) {
		be = bizerr.PasswordWrong("")
		return
	}
	userid = userInfo.Userid
	return
}

func GetUser(userid string) (userInfo *user.UserInfoModel, be *bizerr.Error) {
	userInfo, err := model.GetUserById(userid)
	if err != nil {
		be = bizerr.New(err)
		return
	}
	if userInfo == nil {
		be = bizerr.UserNull()
		return
	}
	if userInfo.Status == consts.UserStatusBlocked {
		be = bizerr.UserBlocked()
		userInfo = nil
		return
	}
	if userInfo.Status == consts.UserStatusCancelled {
		be = bizerr.UserCanceled()
		userInfo = nil
		return
	}
	userInfo.Password = ""
	return
}

func CreateUser(username, password string) (userid string, be *bizerr.Error) {
	userInfo, err := model.GetUser(username)
	if err != nil {
		be = bizerr.New(err)
		return
	}
	if userInfo != nil {
		be = bizerr.UserExists()
		return
	}

	modelLogic := model.CreateUserTask{
		Username: username,
		Password: password,
	}
	userInfo, existed, err := modelLogic.Do()
	if err != nil {
		be = bizerr.New(err)
		return
	}
	if existed {
		be = bizerr.UserExists()
		return
	}
	userid = userInfo.Userid
	return
}
