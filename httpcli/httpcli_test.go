package httpcli

import (
	"github.com/xhigher/hzgo/config"
	"github.com/xhigher/hzgo/defines"
	usermodel "github.com/xhigher/hzgo/demo/model/db/user"
	"github.com/xhigher/hzgo/logger"
	"github.com/xhigher/hzgo/utils"
	"testing"
	"time"
)

type UserProfileResp struct {
	Code int `json:"code"`
	Msg string                    `json:"msg"`
	Data *usermodel.UserInfoModel `json:"data"`
}

func TestHttpCli_PostJSON(t *testing.T) {
	logger.Init(&config.LoggerConfig{
		Filename: "/dev/stdout",
	})
	cli := New(5*time.Second)

	url := "http://127.0.0.1:9000/svc-user/v1/profile"
	data := defines.UseridData{
		Userid: "68c6lcsmkq",
	}

	resp := &UserProfileResp{}
	cli.PostJSON(url, data, resp)
	logger.Infof("user: %v", utils.JSONString(resp))
	time.Sleep(10*time.Second)
}