package main

import (
	"github.com/xhigher/hzgo/config"
	"github.com/xhigher/hzgo/server/notice"
)

func main() {

	config.Init()
	conf := config.GetConfig()

	svr := notice.NewNodeServer(conf)

	svr.Start()
}
