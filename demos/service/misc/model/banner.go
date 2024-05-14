package model

import (
	"github.com/xhigher/hzgo/consts"
	"github.com/xhigher/hzgo/demo/model/db/misc"
	"github.com/xhigher/hzgo/demo/model/db/user"
	"github.com/xhigher/hzgo/logger"
)

func GetAllBannerList() (data map[string][]*misc.BannerItem, err error) {
	var tempData []*misc.BannerInfoModel
	err = user.DB().Model(misc.BannerInfoModel{}).Where("status=?", consts.StatusOnline).Order("sn DESC").Find(&tempData).Error
	if err != nil {
		logger.Errorf("error: %v", err)
		return
	}
	data = make(map[string][]*misc.BannerItem)
	if len(tempData) > 0 {
		for _, item := range tempData {
			data[item.Site] = append(data[item.Site], &misc.BannerItem{
				Id:   item.Id,
				Type: item.Type,
				Name: item.Name,
				Img:  item.Img,
				Data: item.Data,
			})
		}
	}
	return
}
