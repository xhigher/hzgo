package main

import (
	"github.com/xhigher/hzgo/config"
	"github.com/xhigher/hzgo/demos/admin/controller"
	"github.com/xhigher/hzgo/demos/admin/rbac"
	"github.com/xhigher/hzgo/demos/api"
	"github.com/xhigher/hzgo/server/admin"
)

func main() {

	config.Init()
	conf := config.GetConfig()

	rbac.InitPermissions()
	svr := admin.NewServer(conf)
	svr.InitRouters(controller.New(svr.Auth))
	api.Init()
	svr.Start()
}
