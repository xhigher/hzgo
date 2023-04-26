package utils

import (
	"github.com/cloudwego/hertz/pkg/common/json"
	"github.com/xhigher/hzgo/logger"
)

func JSONString(v interface{}) string {
	bs, err := json.Marshal(v)
	if err != nil {
		logger.Errorf("JSONString error: %v", err)
		return ""
	}
	return string(bs)
}

func JSONObject(data string, obj interface{}) error {
	return json.Unmarshal([]byte(data), obj)
}
