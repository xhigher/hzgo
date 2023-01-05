package utils

import (
	"github.com/cloudwego/hertz/pkg/common/json"
	"github.com/xhigher/hzgo/logger"
)

func JsonString(v interface{}) string {
	bs, err := json.Marshal(v)
	if err != nil {
		logger.Errorf("JsonString error: %v", err)
		return ""
	}
	return string(bs)
}
