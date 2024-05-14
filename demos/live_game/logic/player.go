package logic

import (
	"github.com/xhigher/hzgo/demos/live_game/events"
	"github.com/xhigher/hzgo/demos/live_game/maps"
	"github.com/xhigher/hzgo/demos/live_game/model/store"
	"github.com/xhigher/hzgo/logger"
	"github.com/xhigher/hzgo/server/game"
	"github.com/xhigher/hzgo/utils"
	"math/rand"
	"sort"
	"time"
)

type PlayerStatus int

const (
	PlayerInitial PlayerStatus = 0
	PlayerActive  PlayerStatus = 1
	PlayerTrapped PlayerStatus = 2
	PlayerDied    PlayerStatus = 4
)

type PlayerRole int

const (
	PlayerRobot    PlayerRole = 0
	PlayerHuman    PlayerRole = 1
	PlayerAudience PlayerRole = 2
)

type Player struct {
	id          string
	name        string
	avatar      string
	role        PlayerRole
	room        *Room
	lastSite    maps.Site
	curSite     maps.Site
	status      PlayerStatus
	bubbleColor int
	bubblePower int
	bubbleCount int
	pinCount    int
	skin        int
	props       []*Prop
	stepTime    int
	result      int
	awards      []Award
	pipe        *game.Pipe
	index       int
	dieTime     time.Time
}

func NewPlayer(pipe *game.Pipe, user store.UserInfo, role PlayerRole) *Player {
	return &Player{
		id:          user.Id,
		name:        user.Nickname,
		avatar:      user.Avatar,
		role:        role,
		pipe:        pipe,
		stepTime:    playerStepTime,
		skin:        user.Skin,
		bubbleColor: randomBubbleColor(),
		index:       -1,
	}
}

type PlayerData struct {
	Id          string `json:"id"`
	Name        string `json:"name"`
	Avatar      string `json:"avatar"`
	RoomId      int    `json:"roomId"`
	X           int    `json:"x"`
	Y           int    `json:"y"`
	Status      int    `json:"status"`
	StepTime    int    `json:"moveInterval"`
	Skin        int    `json:"role"`
	BubbleColor int    `json:"bombColor"`
	BubbleCount int    `json:"bombCount"`
	PinCount    int    `json:"thumbtackCount"`
}

func (p *Player) GetData() PlayerData {
	msg := PlayerData{
		Id:          p.id,
		Name:        p.name,
		Avatar:      p.avatar,
		X:           p.curSite.X,
		Y:           p.curSite.Y,
		Status:      int(p.status),
		StepTime:    p.stepTime,
		Skin:        p.skin,
		BubbleColor: p.bubbleColor,
		BubbleCount: p.bubbleCount,
		PinCount:    p.pinCount,
	}
	if p.room != nil {
		msg.RoomId = p.room.id
	}
	return msg
}

func (p *Player) IsStop() bool {
	return p.curSite.X == p.lastSite.X && p.curSite.Y == p.lastSite.Y
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

	logger.Infof("SetSite: id=%v, %v, %v", p.id, utils.JSONString(p.curSite), utils.JSONString(p.lastSite))

	//判断该位置是否有道具，如果有，则吃掉
	if yes, prop := p.room.ExistProp(p.curSite); yes {
		logger.Infof("PickUpProp: id=%v, %v", p.id, utils.JSONString(prop.GetData()))
		p.PickUpProp(prop)
	}

	msg := &game.Message{
		Event: events.Move,
		Data: encodeMsgData(&MoveData{
			I: p.index,
			X: p.curSite.X,
			Y: p.curSite.Y,
			T: p.stepTime,
		}),
	}
	p.room.BroadcastMsg(msg)

	//判断该位置是否有被泡住的玩家，如果有，让玩家死亡
	for _, tp := range p.room.players {
		if tp.IsTrapped() && tp.curSite.Equal(s) {
			tp.Die()
		}
	}
}

func (p *Player) PickUpProp(o *Prop) {
	switch o.typ {
	case PropBubble:
		if p.bubbleCount < 10 {
			p.bubbleCount++
		}
	case PropPower:
		if p.bubblePower < 10 {
			p.bubblePower++
		}
	case PropShoes:
		if p.stepTime > 100 {
			p.stepTime -= 20
		}
	case PropPin:
		p.pinCount++
	}

	logger.Infof("PickUpProp: %v", utils.JSONString(o))
	//讲该道具放入玩家吃掉的道具数组中，死亡时爆出来。
	p.props = append(p.props, o)
	p.room.DisappearProp(o, true, p)
}

func (p *Player) CreateBubble() {
	logger.Infof("CreateBubble: id=%v, %v", p.id, p.bubbleCount)
	if p.status != PlayerActive || p.bubbleCount <= 0 {
		return
	}
	if p.room.CheckSiteBubble(p.lastSite) {
		return
	}

	p.room.bubbleId++
	p.bubbleCount++
	bubble := &Bubble{
		Id:       p.room.bubbleId,
		site:     p.lastSite,
		color:    p.bubbleColor,
		power:    p.bubblePower,
		player:   p,
		room:     p.room,
		State:    BubbleAlive,
		bombTime: time.Now().Add(bubbleBombDuration),
	}
	p.room.bubbles = append(p.room.bubbles, bubble)

	msg := &game.Message{
		Event: events.CreateBubble,
		Data:  encodeMsgData(bubble.GetData()),
	}
	p.room.BroadcastMsg(msg)
}

func (p *Player) SendMsg(msg *game.Message) {
	if !p.IsRobot() && p.pipe != nil {
		p.pipe.SendMessage(msg)
	}
}

func (p *Player) Die() {
	if p.room == nil {
		return
	}
	p.ChangeStatus(PlayerDied)
}

func (p *Player) CheckDie(now time.Time) bool {
	if p.IsTrapped() && now.After(p.dieTime) {
		p.Die()
		return true
	}
	return false
}

//被泡住
func (p *Player) Trapped() {
	if p.room == nil {
		return
	}
	p.ChangeStatus(PlayerTrapped)
	p.moveStop()
	p.dieTime = time.Now().Add(6 * time.Second)
}

func (p *Player) IsTrapped() bool {
	return p.status == PlayerTrapped
}

func (p *Player) IsPlaying() bool {
	return p.status == PlayerActive || p.status == PlayerTrapped
}

func (p *Player) moveStop() {
	p.lastSite.X = p.curSite.X
	p.lastSite.Y = p.curSite.Y
	if p.status == PlayerActive {
		if p.room != nil {
			msg := &game.Message{
				Event: events.MoveStop,
				Data: encodeMsgData(&PlayerIndex{
					I: p.index,
				}),
			}
			p.room.BroadcastMsg(msg)
		}
	}
}

func (p *Player) ChangeStatus(status PlayerStatus) {
	if p.status == status {
		return
	}
	p.status = status
	msg := &game.Message{
		Event: events.ChangeUserStatus,
		Data: encodeMsgData(&ChangeUserStatusData{
			Id: p.id,
			//Index:p.index,
			X:      p.curSite.X,
			Y:      p.curSite.Y,
			Status: int(p.status),
		}),
	}
	p.room.BroadcastMsg(msg)
}

func (p *Player) usePin() {
	if p.status == PlayerTrapped && p.pinCount > 0 {
		p.pinCount--
		p.dieTime = time.Time{}
		p.ChangeStatus(PlayerActive)
	}
}

func (p *Player) GetPropsWithDead() (props []*Prop) {
	propCount := len(p.props)
	if propCount == 0 {
		return
	}
	tempProps := p.props
	if propCount >= 3 {
		sort.SliceStable(tempProps, func(i, j int) bool {
			return rand.Intn(10) < 5
		})
		tempProps = tempProps[0:3]
	}

	var nearSites []maps.Site
	for i := -3; i <= 3; i++ {
		for j := -3; j <= 3; j++ {
			if i != 0 && j != 0 {
				site := maps.Site{
					X: p.curSite.X + j,
					Y: p.curSite.Y + i,
				}
				if p.room.mapData.IsValidSite(site) {
					if !p.room.mapData.ExistBox(site) && !p.room.mapData.ExistObstacle(site) {
						nearSites = append(nearSites, site)
					}
				}
			}
		}
	}
	if len(nearSites) == 0 {
		return
	}

	for i := 0; i < len(tempProps); i++ {
		if i < len(nearSites) {
			tempProps[i].site = nearSites[i]
			tempProps[i].disappearTime = time.Now().Add(propAppearDuration)
			p.room.props = append(p.room.props, tempProps[i])
			props = append(props, tempProps[i])
		}
	}

	return
}

func (p *Player) reset() {
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

func (p *Player) IsRobot() bool {
	return p.role == PlayerRobot
}
