package main
import (
	"github.com/xhigher/hzgo/config"
	"github.com/xhigher/hzgo/demo/base"
	"github.com/xhigher/hzgo/demo/controller"
	"github.com/xhigher/hzgo/server"
)
func main()  {

	config.Init()

	svr := server.New(config.GetConfig())
	svr.InitAuth(base.NewAuthMiddleware())

	svr.InitRouter(controller.New())

	svr.Start()
}
