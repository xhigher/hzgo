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

	config := logic.MatchConfig{
		Id: "100",
		RoomPlayerCount: 4,
		RoundCount: 1,
		StartTime: "2023-03-23 08:50:50",
	}
	logic.StartEngine(config)
	svr.Start()
}
