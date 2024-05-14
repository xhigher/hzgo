package misc

import (
	"github.com/xhigher/hzgo/bizerr"
	"github.com/xhigher/hzgo/consts"
	model "github.com/xhigher/hzgo/demo/admin/model/misc"
	"github.com/xhigher/hzgo/demo/model/db/misc"
)

func SaveBannerInfo(id int32, site string, typ int32, name, img, data string) (reload bool, be *bizerr.Error) {
	err := model.SaveBannerInfo(id, site, typ, name, img, data)
	if err != nil {
		be = bizerr.New(err)
		return
	}
	banner, err := model.GetBannerInfo(id)
	if err != nil {
		be = bizerr.New(err)
		return
	}
	if banner.Status == consts.StatusOnline {
		reload = true
	}
	return
}

func SetBannerOnline(id int32) (be *bizerr.Error) {
	err := model.UpdateBannerStatus(id, consts.StatusOnline)
	if err != nil {
		be = bizerr.New(err)
		return
	}
	return
}

func SetBannerOffline(id int32) (be *bizerr.Error) {
	err := model.UpdateBannerStatus(id, consts.StatusOffline)
	if err != nil {
		be = bizerr.New(err)
		return
	}
	return
}

func DeleteBannerInfo(id int32) (be *bizerr.Error) {
	err := model.DeleteBannerInfo(id)
	if err != nil {
		be = bizerr.New(err)
		return
	}
	return
}

func GetBannerInfo(id int32) (data *misc.BannerInfoModel, be *bizerr.Error) {
	data, err := model.GetBannerInfo(id)
	if err != nil {
		be = bizerr.New(err)
		return
	}
	return
}

func GetBannerList(site string, status, offset, limit int32) (total int64, data []*misc.BannerInfoModel, be *bizerr.Error) {
	total, data, err := model.GetBannerList(site, status, offset, limit)
	if err != nil {
		be = bizerr.New(err)
		return
	}
	return
}
