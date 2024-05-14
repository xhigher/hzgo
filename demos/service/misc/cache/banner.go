package cache

import (
	"context"
	"encoding/json"
	"github.com/go-redis/redis/v8"
	"github.com/xhigher/hzgo/consts/cachekey"
	"github.com/xhigher/hzgo/demo/model/db/misc"
	"github.com/xhigher/hzgo/logger"
)

func GetAllBannerList() (data map[string][]*misc.BannerItem, err error) {
	ret, err := client().HGetAll(context.Background(), cachekey.BannerList).Result()
	if err != nil {
		if err == redis.Nil {
			err = nil
			return
		}
		logger.Errorf("error %v", err)
		return
	}

	data = make(map[string][]*misc.BannerItem)
	for key, val := range ret {
		var items []*misc.BannerItem
		err = json.Unmarshal([]byte(val), &items)
		if err != nil {
			logger.Errorf("error %v", err)
			return
		}
		data[key] = items
	}
	return
}

func GetSiteBannerList(site string) (data []*misc.BannerItem, err error) {
	ret, err := client().HGet(context.Background(), cachekey.BannerList, site).Result()
	if err != nil {
		if err == redis.Nil {
			err = nil
			return
		}
		logger.Errorf("error %v", err)
		return
	}

	err = json.Unmarshal([]byte(ret), &data)
	if err != nil {
		logger.Errorf("error %v", err)
		return
	}

	return
}

func SetBannerList(data map[string][]*misc.BannerItem) (err error) {
	values := make([]interface{}, 0)
	var bytes []byte
	for key, items := range data {
		bytes, err = json.Marshal(items)
		if err != nil {
			logger.Errorf("error %v", err)
			return
		}
		values = append(values, key, string(bytes))
	}

	err = client().HSet(context.Background(), cachekey.BannerList, values...).Err()
	if err != nil {
		logger.Errorf("error %v", err)
		return
	}
	return
}
