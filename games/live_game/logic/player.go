package logic

import (
	"github.com/xhigher/hzgo/games/live_game/events"
	"github.com/xhigher/hzgo/games/live_game/maps"
	"github.com/xhigher/hzgo/games/live_game/model/store"
	"github.com/xhigher/hzgo/server/ws"
	"github.com/xhigher/hzgo/utils"
	"time"
)
type PlayerStatus int
const (
	PlayerInitial PlayerStatus = 0
	PlayerActive PlayerStatus = 1
	PlayerTrapped PlayerStatus = 2
	PlayerDied PlayerStatus = 4
)

type PlayerRole int
const (
	PlayerRobot PlayerRole = 0
	PlayerHuman PlayerRole = 1
	PlayerAudience PlayerRole = 2
)

type Player struct {
	id string
	name string
	avatar string
	role PlayerRole
	room *Room
	lastSite maps.Site
	curSite maps.Site
	status PlayerStatus
	bubbleColor int
	bubblePower int
	bubbleCount int
	pinCount int
	props []*Prop
	stepTime int
	result int
	awards []Award
	pipe *ws.Pipe

	dieTime time.Time
}

func NewPlayer(pipe *ws.Pipe, user store.UserInfo, role PlayerRole) *Player{
	return &Player{
		id: user.Id,
		name: user.Name,
		avatar: user.Avatar,
		role: role,
		pipe: pipe,
	}
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

func (p *Player) SetSite(s maps.Site) {
	if p.status != PlayerActive {
		return
	}

	p.lastSite.X = p.curSite.X
	p.lastSite.Y = p.curSite.Y

	p.curSite.X = s.X
	p.curSite.Y = s.Y

	//判断该位置是否有道具，如果有，则吃掉
	if yes, prop := p.room.ExistProp(p.curSite); yes {
		p.PickUpProp(prop)
	}

	msg := &ws.Message{
		Event: events.MoveStop,
		Data: encodeMsgData(&MoveData{
			I: p.id,
			X:p.curSite.X,
			Y:p.curSite.Y,
			T:p.stepTime,
		}),
	}
	p.room.BroadcastMsg(msg)

	//判断该位置是否有被泡住的玩家，如果有，让玩家死亡
}

func (p *Player) PickUpProp(o *Prop){
	switch o.typ {
	case PropBubble:
		if p.bubbleCount < 10 {
			p.bubbleCount++
		}
		break
	case PropPower:
		if p.bubblePower < 10 {
			p.bubblePower ++
		}
		break
	case PropShoes:
		if p.stepTime > 100 {
			p.stepTime -= 20
		}
		break
	case PropPin:
		p.pinCount ++
		break
	}

	//讲该道具放入玩家吃掉的道具数组中，死亡时爆出来。
	p.props = append(p.props, o)
	p.room.DisappearProp(o, true, p)
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
	if !p.IsRobot() && p.pipe != nil {
		p.pipe.SendMessage(msg)
	}
}

func (p *Player) Die() {
	if p.room == nil {
		return
	}
	p.status = PlayerDied
}

//被泡住
func (p *Player) Trapped(){
	if p.room == nil {
		return
	}
	p.status = PlayerTrapped
	p.moveStop()
	p.dieTime = time.Now().Add(6*time.Second)
}

func (p *Player) IsTrapped() bool {
	return p.status == PlayerTrapped
}

func (p *Player) IsPlaying() bool {
	return p.status == PlayerActive || p.status == PlayerTrapped
}

func (p *Player) moveStop(){
	p.lastSite.X = p.curSite.X
	p.lastSite.Y = p.curSite.Y
	if p.status == PlayerActive {
		if p.room != nil {
			msg := &ws.Message{
				Event: events.MoveStop,
				Data: encodeMsgData(&PlayerId{
					Id: p.id,
				}),
			}
			p.room.BroadcastMsg(msg)
		}
	}
}

func (p *Player) usePin(){
	if p.status ==PlayerTrapped && p.pinCount > 0 {
		p.pinCount --
		p.dieTime = time.Time{}
		p.status = PlayerActive
	}
}

func (p *Player) reset(){
	p.role = 0
	p.status = PlayerInitial
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

func (p *Player) IsRobot() bool{
	return p.role == PlayerRobot
}