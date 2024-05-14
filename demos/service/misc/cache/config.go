package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/go-redis/redis/v8"
	"github.com/xhigher/hzgo/consts/cachekey"
	"github.com/xhigher/hzgo/consts/cachetime"
	"github.com/xhigher/hzgo/demo/model/db/misc"
	"github.com/xhigher/hzgo/logger"
)

func GetConfigList() (sum string, data map[string]*misc.ConfigInfo, err error) {
	sum, err = client().Get(context.Background(), cachekey.ConfigListSum).Result()
	if err != nil {
		if err == redis.Nil {
			err = nil
			return
		}
		logger.Errorf("error %v", err)
		return
	}

	ret, err := client().Get(context.Background(), cachekey.ConfigList).Result()
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

func SetConfigList(sum string, data map[string]*misc.ConfigInfo) (err error) {
	err = client().Set(context.Background(), cachekey.ConfigListSum, sum, cachetime.Permanent).Err()
	if err != nil {
		logger.Errorf("error %v", err)
		return
	}

	bytes, err := json.Marshal(data)
	if err != nil {
		logger.Errorf("error %v", err)
		return
	}
	err = client().Set(context.Background(), cachekey.ConfigList, string(bytes), cachetime.Permanent).Err()
	if err != nil {
		logger.Errorf("error %v", err)
		return
	}
	return
}

func GetConfigInfo(id string) (data *misc.ConfigInfoModel, err error) {
	ret, err := client().Get(context.Background(), fmt.Sprintf(cachekey.ConfigInfo, id)).Result()
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

func SetConfigInfo(id string, data *misc.ConfigInfoModel) (err error) {
	bytes, err := json.Marshal(data)
	if err != nil {
		logger.Errorf("error %v", err)
		return
	}
	err = client().Set(context.Background(), fmt.Sprintf(cachekey.ConfigInfo, id), string(bytes), cachetime.Permanent).Err()
	if err != nil {
		logger.Errorf("error %v", err)
		return
	}
	return
}
