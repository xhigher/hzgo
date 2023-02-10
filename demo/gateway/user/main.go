package main

import (
	"github.com/xhigher/hzgo/config"
	"github.com/xhigher/hzgo/demo/api"
	"github.com/xhigher/hzgo/demo/gateway/user/controller"
	"github.com/xhigher/hzgo/server/gateway"
)

func main() {

	config.Init()
	conf := config.GetConfig()

	svr := gateway.NewServer(conf)
	svr.InitRouter(controller.NewWithAuth(conf.Name, svr.Auth))
	api.Init()
	svr.Start()
}
