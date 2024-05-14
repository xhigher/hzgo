package logic

import (
	"github.com/xhigher/hzgo/bizerr"
	"github.com/xhigher/hzgo/demo/model/db/misc"
	"github.com/xhigher/hzgo/demo/service/misc/cache"
	"github.com/xhigher/hzgo/demo/service/misc/model"
)

func GetAllBannerList() (data map[string][]*misc.BannerItem, be *bizerr.Error) {
	var err error
	data, err = cache.GetAllBannerList()
	if err != nil {
		be = bizerr.New(err)
		return
	}
	if data != nil {
		return
	} else {
		data, err = model.GetAllBannerList()
		if err != nil {
			be = bizerr.New(err)
			return
		}
		cache.SetBannerList(data)
	}
	return
}

func GetSiteBannerList(site string) (data []*misc.BannerItem, be *bizerr.Error) {
	var err error
	data, err = cache.GetSiteBannerList(site)
	if err != nil {
		be = bizerr.New(err)
		return
	}
	return
}

func ReloadBanner() (be *bizerr.Error) {
	data, err := model.GetAllBannerList()
	if err != nil {
		be = bizerr.New(err)
		return
	}
	err = cache.SetBannerList(data)
	if err != nil {
		be = bizerr.New(err)
		return
	}
	return
}
