package user

import (
	"github.com/xhigher/hzgo/bizerr"
	"github.com/xhigher/hzgo/consts"
	model "github.com/xhigher/hzgo/demo/admin/model/user"
	"github.com/xhigher/hzgo/demo/model/db/user"
)

func GetUser(userid string) (userInfo *user.UserInfoModel, be *bizerr.Error) {
	userInfo, err := model.GetUserInfo(userid)
	if err != nil {
		be = bizerr.New(err)
		return
	}
	if userInfo == nil {
		be = bizerr.UserNull()
		return
	}
	userInfo.Password = ""
	return
}

func GetUserList(status, offset, limit int32) (total int64, userList []*user.UserInfoModel, be *bizerr.Error) {
	total, userList, err := model.GetUserList(status, offset, limit)
	if err != nil {
		be = bizerr.New(err)
		return
	}
	return
}

func StartUser(userid string) (be *bizerr.Error) {
	userInfo, err := model.GetUserInfo(userid)
	if err != nil {
		be = bizerr.New(err)
		return
	}
	if userInfo == nil {
		be = bizerr.UserNull()
		return
	}
	if userInfo.Status == consts.UserStatusActive {
		return
	}
	if userInfo.Status == consts.UserStatusCancelled {
		be = bizerr.UserCanceled()
		return
	}

	err = model.UpdateUserStatus(userid, consts.UserStatusActive)
	if err != nil {
		be = bizerr.New(err)
		return
	}

	CleanStaffToken(userid)

	return
}

func StopUser(userid string) (be *bizerr.Error) {
	userInfo, err := model.GetUserInfo(userid)
	if err != nil {
		be = bizerr.New(err)
		return
	}
	if userInfo == nil {
		be = bizerr.UserNull()
		return
	}
	if userInfo.Status == consts.UserStatusBlocked {
		return
	}
	if userInfo.Status == consts.UserStatusCancelled {
		be = bizerr.UserCanceled()
		return
	}

	err = model.UpdateUserStatus(userid, consts.UserStatusBlocked)
	if err != nil {
		be = bizerr.New(err)
		return
	}

	CleanStaffToken(userid)

	return
}

func CleanStaffToken(userid string) (be *bizerr.Error) {
	err := model.CleanUserTokens(userid)
	if err != nil {
		be = bizerr.New(err)
		return
	}

	return
}
