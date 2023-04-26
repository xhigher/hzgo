package httpcli

import (
	"github.com/xhigher/hzgo/config"
	"github.com/xhigher/hzgo/logger"
	"testing"
)

type UserProfileResp struct {
	Code int                      `json:"code"`
	Msg  string                   `json:"msg"`
	Data interface{} `json:"data"`
}

func TestHttpCli_PostJSON(t *testing.T) {
	logger.Init(&config.LoggerConfig{
		Filename: "/dev/stdout",
	})

}
