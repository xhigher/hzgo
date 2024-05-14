package logic

import (
	"github.com/xhigher/hzgo/bizerr"
	"github.com/xhigher/hzgo/demo/service/user/model"
)

func CheckToken(userid, token string) (ok bool, be *bizerr.Error) {
	ok, err := model.CheckToken(userid, token)
	if err != nil {
		be = bizerr.New(err)
		return
	}
	return
}

func UpdateToken(userid, token string, et, it int64) (be *bizerr.Error) {
	err := model.SaveToken(userid, token, et, it)
	if err != nil {
		be = bizerr.New(err)
		return
	}
	return
}
