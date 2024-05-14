package model

import (
	"github.com/xhigher/hzgo/consts"
	"github.com/xhigher/hzgo/demos/model/db/misc"
	"github.com/xhigher/hzgo/demos/model/db/user"
	"github.com/xhigher/hzgo/logger"
	"github.com/xhigher/hzgo/utils"
	"gorm.io/gorm"
	"sort"
	"strings"
)

func getConfigListSum(ids []string, data map[string]*misc.ConfigInfo) string {
	sort.Strings(ids)
	strArray := make([]string, len(ids))
	for i, id := range ids {
		item := data[id]
		strArray[i] = id + ":" + item.Name + ":" + item.Items
	}
	return utils.MD5(strings.Join(strArray, ";"))
}

func GetConfigList() (sum string, data map[string]*misc.ConfigInfo, err error) {
	var tempData []*misc.ConfigInfoModel
	err = user.DB().Model(misc.ConfigInfoModel{}).Select("id,name,items").Where("status=? AND static=?", consts.StatusOnline, consts.YES).Find(&tempData).Error
	if err != nil {
		logger.Errorf("error: %v", err)
		return
	}
	data = make(map[string]*misc.ConfigInfo)
	if len(tempData) > 0 {
		ids := make([]string, len(data))
		for i, item := range tempData {
			data[item.Id] = &misc.ConfigInfo{
				Name:  item.Name,
				Items: item.Items,
			}
			ids[i] = item.Id
		}
		sum = getConfigListSum(ids, data)
	}
	return
}

func GetConfigInfo(id string) (data *misc.ConfigInfoModel, err error) {
	err = user.DB().First(&data, "id=?", id).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			data = nil
			err = nil
			return
		}
		logger.Errorf("error: %v", err)
	}
	return
}
