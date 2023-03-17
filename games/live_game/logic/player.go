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
	PlayerExit PlayerStatus = 2
	PlayerDied PlayerStatus = 3

)

type Player struct {
	id string
	name string
	avatar string
	role int
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

type PlayerData struct {
	Id string `json:"id"`
	Name string `json:"name"`
	Avatar string `json:"avatar"`
	Role int `json:"role"`
	RoomId int `json:"room_id"`
	Site maps.Site `json:"site"`
	Status int `json:"status"`
	StepTime int `json:"step_time"`
}

func (p *Player) GetData() PlayerData{
	msg := PlayerData{
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
	if p.pipe != nil {
		p.pipe.SendMessage(msg)
	}
}


func (p *Player) reset(){
	p.role = 0
	p.status = PlayerDied
	p.room = nil
	p.bubbleColor = 0
	p.bubblePower = 1
	p.bubbleCount = 0
	p.lastSite = maps.Site{}
	p.curSite = maps.Site{}
	p.pinCount = 0
	p.props = nil
	p.stepTime = 0
	p.result = 0
	p.awards = nil
	p.pipe = nil
}