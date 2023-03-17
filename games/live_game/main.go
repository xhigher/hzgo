package main

import (
	"github.com/xhigher/hzgo/config"
	"github.com/xhigher/hzgo/games/live_game/logic"
	"github.com/xhigher/hzgo/server/ws"
)

func main() {

	config.Init()
	conf := config.GetConfig()

	svr := ws.NewServer(conf, logic.NewHandler())
	logic.StartEngine()
	svr.Start()
}
