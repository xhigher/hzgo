package logic

import (
	"github.com/xhigher/hzgo/games/live_game/maps"
	"github.com/xhigher/hzgo/server/ws"
	"github.com/xhigher/hzgo/utils"
)
type PlayerStatus int
const (
	PlayerActive PlayerStatus = 0
	PlayerCaught PlayerStatus = 1
	PlayerDied PlayerStatus = 2
)

type Player struct {
	id string
	name string
	avatar string
	room *Room
	lastSite maps.Site
	curSite maps.Site
	status PlayerStatus
	bubbleColor int
	bubblePower int
	bubbleCount int
	pinCount int
	props []Prop
	stepTime int
	result int
	awards []Award
	pipe *ws.Pipe

}

type PlayerMsg struct {
	Id string
	Name string
	Avatar string
	RoomId int
	Site maps.Site
	Status int
	StepTime int
}

func (p *Player) GetMsg() PlayerMsg{
	msg := PlayerMsg{
		Id:       p.id,
		Name:     p.name,
		Avatar:   p.avatar,
		Site:     p.curSite,
		Status:   int(p.status),
		StepTime: p.stepTime,
	}
	if p.room != nil {
		msg.RoomId = p.room.id
	}
	return msg
}

func (p *Player) IsStop() bool {
	return p.curSite.X== p.lastSite.X && p.curSite.Y== p.lastSite.Y
}

func (p *Player) Stop() {
	p.lastSite.X = p.curSite.X
	p.lastSite.Y = p.curSite.Y
}

func (p *Player) SetSite() {
	p.lastSite.X = p.curSite.X
	p.lastSite.Y = p.curSite.Y
}

func (p *Player) CreateBubble() {
	if p.status != PlayerActive || p.bubbleCount<=0 {
		return
	}
	if p.room.CheckSiteBubble(p.lastSite) {
		return
	}

	p.room.bubbleId ++
	p.bubbleCount ++
	bubble := &Bubble{
		Id:     p.room.bubbleId,
		site:   p.lastSite,
		color:  p.bubbleColor,
		power:  p.bubblePower,
		player: p,
		room:   p.room,
		ct:     utils.NowTime(),
	}
	p.room.bubbles = append(p.room.bubbles, bubble)
}

func (p *Player) SendMsg(msg *ws.Message) {
	p.pipe.SendMessage(msg)
}
