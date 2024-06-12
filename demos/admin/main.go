package main

import (
	"github.com/xhigher/hzgo/config"
	"github.com/xhigher/hzgo/demos/admin/controller"
	"github.com/xhigher/hzgo/demos/admin/logic/platform"
	"github.com/xhigher/hzgo/server/admin"
)

func main() {

	config.Init()
	conf := config.GetConfig()

	svr := admin.NewServer(conf)
	svr.InitRouters(controller.New(svr.Auth))
	initPermissions()
	svr.Start()
}

func initPermissions() {
	data, err := platform.GetAllRolesPermissions()
	if err != nil {
		return
	}
	admin.InitRolePermissions(data)
}
