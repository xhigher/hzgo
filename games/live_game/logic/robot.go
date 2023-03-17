package logic

import (
	"github.com/xhigher/hzgo/utils"
	"sync"
)

var (
	robots = sync.Pool{
		New: func() interface{} {
			return newRobot(0)
		},
	}
)

type Robot struct {
	player *Player
	loadProcess int
}

func initRobots(count int){
	for i:=0; i<count; i++ {
		robots.Put(newRobot(i))
	}
}

func newRobot(i int) *Robot{
	return &Robot{
		player: &Player{
			id: utils.IntToBase36(utils.NowTime()-725846400+int64(i)),
			name: utils.RandString(20),
			avatar: "",
		},
	}
}

func getRobot() *Robot{
	return robots.Get().(*Robot)
}

func releaseRobot(r *Robot){
	r.player.reset()
	robots.Put(r)
}