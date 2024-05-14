package main

import (
	"github.com/xhigher/hzgo/config"
	"github.com/xhigher/hzgo/demo/admin/controller"
	"github.com/xhigher/hzgo/demo/admin/rbac"
	"github.com/xhigher/hzgo/demo/api"
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
