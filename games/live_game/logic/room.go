package logic

import (
	"context"
	"github.com/xhigher/hzgo/games/live_game/events"
	"github.com/xhigher/hzgo/games/live_game/maps"
	"github.com/xhigher/hzgo/server/ws"
	"github.com/xhigher/hzgo/utils"
	"math/rand"
)

type RoomStatus int
const (
	Idle RoomStatus = 0
	Loading RoomStatus = 1
	Ongoing RoomStatus = 2
	End RoomStatus = 3
)

type Room struct {
	id int
	typ int
	status RoomStatus
	st int64
	et int64
	winner int
	aliveNum int
	mapData maps.MapData
	players []*Player
	robots []*Robot
	bubbles []*Bubble
	props []*Prop
	bubbleId int
	TickChan chan int
	context context.Context
	closeFunc context.CancelFunc
	tickCount int
	audiences map[string]*Player
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

func NewRoom(m *Match, id int, playerCount int) *Room{
	mapData := maps.GetMap(playerCount)
	r := &Room{
		id: id,
		status: Idle,
		st: utils.NowTime(),
		et:0,
		winner:0,
		aliveNum: 0,
		mapData: mapData,
		bubbleId:0,
		TickChan:make(chan int, 1),
		tickCount:0,
	}

	r.context, r.closeFunc = context.WithCancel(m.roomContext)

	go func() {
		for {
			select {
			case <-r.TickChan:
				r.HandleTick()
			case <-r.context.Done():

				return
			}
		}
	}()

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

func (r *Room) HandleTick(){
	switch r.status {
	case Loading:
		r.handleLoading()
	case Ongoing:

	}
}

func (r *Room) handleLoading(){
	r.tickCount ++
	if r.tickCount % 5 == 0 {//500ms
		loading := 0
		for _, rb := range r.robots {
			rb.loadProcess = rb.loadProcess + rand.Intn(20)
			if rb.loadProcess > 100 {
				rb.loadProcess = 100
			}
			msg := &ws.Message{
				Event: events.LoadProcess,
				Data: encodeMsgData(LoadProcessData{
					Id: r.id,
					Process: rb.loadProcess,
				}),
			}
			r.BroadcastMsg(msg)
			if rb.loadProcess >= 100 {
				loading ++
			}
		}
		if loading == len(r.robots) {
			r.readyGo()
		}
	}
}

func (r *Room) readyGo(){
	r.status = Ongoing
	msg := &ws.Message{
		Event: events.GameReady,
	}
	r.BroadcastMsg(msg)
}

func (r *Room) Finish(){
	r.closeFunc()
}

func (r *Room) JoinPlayer(player *Player){

}

func (r *Room) ExistPlayer(player *Player) bool {
	for _,p := range r.players {
		if p.id == player.id {
			return true
		}
	}
	return false
}

func (r *Room) CheckSiteBubble(site maps.Site) bool {
	for _,b := range r.bubbles {
		if b.site.Equal(site) {
			return true
		}
	}
	return false
}

func (r Room) BroadcastMsg(msg *ws.Message){
	if r.status == End {
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

func (r Room) RoundStart(){
	data := &RoundStartData{
		RoomId: r.id,
		Type:    r.typ,
		Status:  int(r.status),
		Map:     r.mapData,
		Users: r.getPlayersData(),
		Bubbles:  nil,
		Props:   r.getPropsData(),
	}
	msg := &ws.Message{
		Event: events.GameStart,
		Data: encodeMsgData(data),
	}

	r.BroadcastMsg(msg)
}

func  (r Room) Exit(id string){
	if _, ok := r.audiences[id]; ok {
		delete(r.audiences, id)
	}
	for _,p := range r.players {
		if p.id == id {
			p.status = PlayerExit
		}
	}
}