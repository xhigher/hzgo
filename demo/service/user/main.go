package main

import (
	"github.com/xhigher/hzgo/config"
	"github.com/xhigher/hzgo/demo/service/user/controller"
	"github.com/xhigher/hzgo/server/service"
)

func main() {

	config.Init()
	conf := config.GetConfig()

	svr := service.NewServer(conf)
	svr.InitRouter(controller.New(conf.Name))
	svr.Start()
}
