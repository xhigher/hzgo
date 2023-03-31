package logic

import (
	"context"
	"github.com/xhigher/hzgo/games/live_game/events"
	"github.com/xhigher/hzgo/games/live_game/maps"
	"github.com/xhigher/hzgo/logger"
	"github.com/xhigher/hzgo/server/ws"
	"github.com/xhigher/hzgo/utils"
	"math/rand"
	"sort"
	"time"
)

type RoomStatus int
const (
	Idle RoomStatus = 0
	Reading RoomStatus = 1
	Ongoing RoomStatus = 2
	Ended RoomStatus = 3
)

type Room struct {
	id int
	typ int
	status RoomStatus
	st int64
	et int64
	aliveNum int
	mapData maps.MapData
	players []*Player
	bubbles []*Bubble
	robots []*Robot
	props []*Prop
	bubbleId int
	propId int
	TickChan chan int
	context context.Context
	closeFunc context.CancelFunc
	audiences map[string]*Player
	propWeights []PropWeight
	endTime time.Time
	winner *Player
}

type RoomData struct {
	Id int `json:"id"`
	Type int `json:"type"`
	Status int `json:"status"`
	Map maps.MapData `json:"map"`
	Players []PlayerData `json:"players"`
	Bubble []BubbleData `json:"bubble"`
	Props []PropData `json:"props"`
}

type RoomResult struct {
	Id int `json:"id"`
	Winner int `json:"winner"`
}

func NewRoom(id int, playerCount int) *Room{
	mapData := maps.GetMap(playerCount)
	r := &Room{
		id: id,
		status: Idle,
		st: utils.NowTime(),
		et:0,
		winner:nil,
		aliveNum: 0,
		mapData: mapData,
		bubbleId:0,
		propWeights:[]PropWeight{
			{
				Type: 0,
				Weight: 20,
			},
			{
				Type: 1,
				Weight: 20,
			},
			{
				Type: 2,
				Weight: 20,
			},
			{
				Type: 3,
				Weight: 1,
			},
		},
	}
	return r
}

func (r *Room) GetData() RoomData{
	return RoomData{
		Id:      r.id,
		Type:    r.typ,
		Status:  int(r.status),
		Map:     r.mapData,
		Players: r.getPlayersData(),
		Bubble:  r.getBubblesData(),
		Props:   r.getPropsData(),
	}
}

func (r *Room) getPlayersData() []PlayerData{
	var msg []PlayerData
	for _, p := range r.players {
		msg = append(msg, p.GetData())
	}
	return msg
}

func (r *Room) getBubblesData() []BubbleData{
	var msg []BubbleData
	for _, b := range r.bubbles {
		msg = append(msg, b.GetData())
	}
	return msg
}

func (r *Room) getPropsData() []PropData{
	var msg []PropData
	for _, p := range r.props {
		msg = append(msg, p.GetData())
	}
	return msg
}


func (r *Room) handleTimerEvent(now time.Time){
	if r.status != Ongoing {
		return
	}
	r.handleProps(now)
	r.handleBubbleBomb(now)
	if r.checkPlayerDie(now){
		r.checkRoundOver()
	}

	//if now.After(r.endTime) {
	//	r.status = Ended
	//	r.handleResult()
	//}

}

func (r *Room) handleProps(now time.Time){
	for _, b := range r.props {
		if now.After(b.disappearTime) {

		}
	}
}

func (r *Room) handleBubbleBomb(now time.Time){
	for _, b := range r.bubbles {
		r.checkBubbleBomb(b, now)
	}
}

func (r *Room) checkPlayerDie(now time.Time) bool {
	yes := false
	for _, p := range r.players {
		if p.CheckDie(now){
			r.aliveNum --
			yes = true
		}
	}
	return yes
}

func (r *Room) checkRoundOver() {
	if r.aliveNum <= 1 {
		r.status = Ended
		results := make([]RoundOverResult, len(r.players))
		for i, p := range r.players {
			win := 0
			if p.IsPlaying(){
				win = 1
			}
			results[i] = RoundOverResult{
				Player: p.id,
				Index: i,
				Win: win,
			}
		}
		msg := &ws.Message{
			Event: events.GameOver,
			Data: encodeMsgData(&RoundOverData{
				Result: results,
				Win: 5,
				Lose: -5,
			}),
		}
		r.BroadcastMsg(msg)
	}
}

func (r *Room) ReadyGo(endTime time.Time){
	r.status = Ongoing
	r.endTime = endTime
	msg := &ws.Message{
		Event: events.GameReady,
	}
	r.BroadcastMsg(msg)

	for _,p := range r.players {
		p.status = PlayerActive
	}

	r.RunRobotPlayers()
}

func (r *Room) handleResult(){

}

func (r *Room) Finish(){
	r.closeFunc()
}

func (r *Room) JoinPlayer(player *Player){
	player.pinCount = 1
	player.bubbleCount = 2
	player.bubblePower = 1
	player.room = r
	r.players = append(r.players, player)
}

func (r *Room) JoinRobot(robot *Robot){
	logger.Infof("JoinRobot id: %v", robot.id)
	robot.Player.room = r
	r.robots = append(r.robots, robot)
	r.players = append(r.players, robot.Player)
}

func (r *Room) CheckSiteBubble(site maps.Site) bool {
	for _,b := range r.bubbles {
		if b.site.Equal(site) {
			return true
		}
	}
	return false
}

func (r *Room) BroadcastMsg(msg *ws.Message){
	if r.status != Reading && r.status != Ongoing {
		return
	}
	for _,p := range r.players {
		p.SendMsg(msg)
	}
	for _,p := range r.audiences {
		p.SendMsg(msg)
	}
}

func (r Room) SendPlayerMsg(p *Player, msg *ws.Message){
	p.SendMsg(msg)
	for _,a := range p.room.audiences {
		a.SendMsg(msg)
	}
}

func (r *Room) RoundStart(){
	r.status = Reading
	bornSites := r.mapData.BornSites
	sort.SliceStable(bornSites, func(i, j int) bool {
		return utils.RandInt(0,100)<50
	})

	//给玩家赋值初始坐标
	for i, p := range r.players {
		p.curSite = bornSites[i]
		p.lastSite = bornSites[i]
	}

	data := &RoundStartData{
		RoomId: r.id,
		Type:    r.typ,
		Status:  int(r.status),
		Map:     r.mapData,
		Player: r.getPlayersData(),
		Bubbles:  nil,
		Props:   r.getPropsData(),
	}
	msg := &ws.Message{
		Event: events.GameStart,
		Data: encodeMsgData(data),
	}

	r.BroadcastMsg(msg)
}

func  (r *Room) IsEnded() bool{
	return r.status == Ended
}

func  (r *Room) Exit(id string){
	if _, ok := r.audiences[id]; ok {
		delete(r.audiences, id)
	}
	for _,p := range r.players {
		if p.id == id {
			p.status = PlayerDied
		}
	}
}

func (r *Room) DeleteBubble(b *Bubble) bool{
	for i, tb := range r.bubbles {
		if tb.Id == b.Id {
			if i > 0 {
				r.bubbles = append(r.bubbles[:i], r.bubbles[i+1:]...)
			}else{
				r.bubbles = r.bubbles[1:]
			}
			return true
		}
	}
	return false
}

func (r *Room) ExistBubble(site maps.Site) (bool, *Bubble){
	for _, tb := range r.bubbles {
		if tb.site.Equal(site) {
			return true,tb
		}
	}
	return false, nil
}

func (r *Room) ExistProp(site maps.Site) (bool, *Prop){
	for _, tb := range r.props {
		if tb.site.Equal(site) {
			return true,tb
		}
	}
	return false, nil
}

func (r *Room) ExistPlayer(site maps.Site) (bool, *Player){
	for _, tp := range r.players {
		if tp.curSite.Equal(site) {
			return true,tp
		}
	}
	return false, nil
}

func (r *Room) HasHumanPlayer() bool{
	for _, tp := range r.players {
		if !tp.IsRobot() && tp.IsPlaying() {
			return true
		}
	}
	return false
}

func (r *Room) checkBubbleBomb(b *Bubble, now time.Time){
	result := b.CheckBomb(now)
	if result == nil {
		return
	}

	//在地图中销毁需要销毁的箱子
	r.destroyBoxes(result)

	msg := &ws.Message{
		Event: events.BubbleBomb,
		Data:encodeMsgData(result),
	}
	r.BroadcastMsg(msg)

	r.checkBombedUsers(result)

}

func (r *Room) destroyBoxes(result *BombResult){
	for _,bx := range result.Boxes {
		r.mapData.Boxes[bx.Y][bx.X] = 0
		rn := rand.Int31n(100)
		//30%几率产生道具
		if rn < 30 {
			pw := WeightRandom(r.propWeights)
			prop := &Prop{
				id: r.propId,
				site: bx,
				typ: pw.Type,
				roomId: r.id,
			}
			r.props = append(r.props, prop)
			r.propId ++
		}
	}
}

//判断爆炸点数组是否有该用户，如果有，则泡住
func (r *Room) checkBombedUsers(result *BombResult){
	for _,p := range r.players {
		for _,a := range result.Areas {
			if a.Equal(p.curSite) {
				if p.status == 0 {
					p.Trapped()
				}else if p.status == 10 {
					p.Die()
				}
			}
		}
	}
}

func (r *Room) DisappearProp(p *Prop, isPickUp bool, player *Player){
	pi := -1
	for i, tb := range r.props {
		if tb.site.Equal(p.site) {
			pi = i
			break
		}
	}
	if pi >= 0 {
		if pi > 0 {
			r.props = append(r.props[:pi], r.props[pi+1:]...)
		}else{
			r.props = r.props[1:]
		}
		msg := &ws.Message{
			Event:events.DisappearProp,
			Data: encodeMsgData(&PropDisappearData{
				Id: p.id,
				PickUp: isPickUp,
				Player: player.id,
			}),
		}
		r.BroadcastMsg(msg)
	}
}


func (r *Room) RunRobotPlayers(){
	for _, pr := range r.robots {
		go pr.Run()
	}
}
