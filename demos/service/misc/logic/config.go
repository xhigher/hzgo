package logic

import (
	"github.com/xhigher/hzgo/bizerr"
	"github.com/xhigher/hzgo/demo/model/db/misc"
	"github.com/xhigher/hzgo/demo/service/misc/cache"
	"github.com/xhigher/hzgo/demo/service/misc/model"
)

func GetConfigList() (sum string, data map[string]*misc.ConfigInfo, be *bizerr.Error) {
	var err error
	sum, data, err = cache.GetConfigList()
	if err != nil {
		be = bizerr.New(err)
		return
	}
	if data != nil {
		return
	} else {
		sum, data, err = model.GetConfigList()
		if err != nil {
			be = bizerr.New(err)
			return
		}
		cache.SetConfigList(sum, data)
	}
	return
}

func GetConfigInfo(id string) (data *misc.ConfigInfoModel, be *bizerr.Error) {
	var err error
	data, err = cache.GetConfigInfo(id)
	if err != nil {
		be = bizerr.New(err)
		return
	}
	if data != nil {
		return
	} else {
		data, err = model.GetConfigInfo(id)
		if err != nil {
			be = bizerr.New(err)
			return
		}
		cache.SetConfigInfo(id, data)
	}
	return
}

func ReloadConfig(id string) (be *bizerr.Error) {
	sum, data, err := model.GetConfigList()
	if err != nil {
		be = bizerr.New(err)
		return
	}
	err = cache.SetConfigList(sum, data)
	if err != nil {
		be = bizerr.New(err)
		return
	}
	if len(id) == 0 {
		return
	}
	info, err := model.GetConfigInfo(id)
	if err != nil {
		be = bizerr.New(err)
		return
	}
	err = cache.SetConfigInfo(id, info)
	if err != nil {
		be = bizerr.New(err)
		return
	}
	return
}
