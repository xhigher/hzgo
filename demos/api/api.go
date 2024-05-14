package api

import (
	"github.com/xhigher/hzgo/consts"
	"github.com/xhigher/hzgo/svcmgr"
)

func Init() {
	svcmgr.Init([]svcmgr.SvcConf{
		{
			Name:     svcUser.name,
			AddrList: []string{"0.0.0.0:9000"},
			Timeout:  consts.TimeSecond5,
		},
	})
}
