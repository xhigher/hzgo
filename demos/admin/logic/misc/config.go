package misc

import (
	"github.com/xhigher/hzgo/bizerr"
	"github.com/xhigher/hzgo/consts"
	model "github.com/xhigher/hzgo/demo/admin/model/misc"
	"github.com/xhigher/hzgo/demo/model/db/misc"
)

func SaveConfigInfo(id, name, items string, static bool, filters string) (reload bool, be *bizerr.Error) {
	err := model.SaveConfigInfo(id, name, items, static, filters)
	if err != nil {
		be = bizerr.New(err)
		return
	}
	data, err := model.GetConfigInfo(id)
	if err != nil {
		be = bizerr.New(err)
		return
	}
	if data.Status == consts.StatusOnline {
		reload = true
	}
	return
}

func SetConfigOnline(id string) (be *bizerr.Error) {
	err := model.UpdateConfigStatus(id, consts.StatusOnline)
	if err != nil {
		be = bizerr.New(err)
		return
	}
	return
}

func SetConfigOffline(id string) (be *bizerr.Error) {
	err := model.UpdateConfigStatus(id, consts.StatusOffline)
	if err != nil {
		be = bizerr.New(err)
		return
	}
	return
}

func DeleteConfigInfo(id string) (be *bizerr.Error) {
	err := model.DeleteConfigInfo(id)
	if err != nil {
		be = bizerr.New(err)
		return
	}
	return
}

func GetConfigInfo(id string) (data *misc.ConfigInfoModel, be *bizerr.Error) {
	data, err := model.GetConfigInfo(id)
	if err != nil {
		be = bizerr.New(err)
		return
	}
	return
}

func GetConfigList(status, offset, limit int32) (total int64, data []*misc.ConfigInfoModel, be *bizerr.Error) {
	total, data, err := model.GetConfigList(status, offset, limit)
	if err != nil {
		be = bizerr.New(err)
		return
	}
	return
}
