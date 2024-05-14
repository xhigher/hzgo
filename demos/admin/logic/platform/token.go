package platform

import (
	"github.com/xhigher/hzgo/bizerr"
	model "github.com/xhigher/hzgo/demo/admin/model/platform"
	"github.com/xhigher/hzgo/logger"
)

func TokenCheck(uid, token string) (ok bool, be *bizerr.Error) {
	ok, err := model.CheckToken(uid, token)
	if err != nil {
		logger.Errorf("token check error: %v, %v", uid, token)
		be = bizerr.New(err)
		return
	}
	return
}

func TokenUpdate(uid, token string, et, it int64) (be *bizerr.Error) {
	err := model.SaveToken(uid, token, et, it)
	if err != nil {
		logger.Errorf("token check error: %v, %v", uid, token)
		be = bizerr.New(err)
		return
	}
	return
}
