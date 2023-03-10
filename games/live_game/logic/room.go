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
	map2 maps.Map
	players []*Player
	bubbles []*Bubble
	props []*Prop
	bubbleId int
	TickChan chan int
	context context.Context
	closeFunc context.CancelFunc
	loadProcess int
	tickCount int
}

type RoomMsg struct {
	Id int `json:"id"`
	Type int `json:"type"`
	Status int `json:"status"`
	Map maps.Map `json:"map"`
	Players []PlayerMsg `json:"players"`
	Bubble []BubbleMsg `json:"bubble"`
	Props []PropMsg `json:"props"`
}

type RoomResult struct {
	Id int `json:"id"`
	Winner int `json:"winner"`
}

func NewRoom(m *Match, id int, playerCount int) *Room{
	map2 := maps.GetMap(playerCount)
	r := &Room{
		id: id,
		status: Idle,
		st: utils.NowTime(),
		et:0,
		winner:0,
		aliveNum: 0,
		map2: map2,
		bubbleId:0,
		TickChan:make(chan int, 1),
		loadProcess:0,
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

func (r *Room) GetMsg() RoomMsg{
	return RoomMsg{
		Id:      r.id,
		Type:    r.typ,
		Status:  int(r.status),
		Map:     r.map2,
		Players: r.getPlayersMsg(),
		Bubble:  r.getBubblesMsg(),
		Props:   r.getPropsMsg(),
	}
}

func (r *Room) getPlayersMsg() []PlayerMsg{
	var msg []PlayerMsg
	for _, p := range r.players {
		msg = append(msg, p.GetMsg())
	}
	return msg
}

func (r *Room) getBubblesMsg() []BubbleMsg{
	var msg []BubbleMsg
	for _, b := range r.bubbles {
		msg = append(msg, b.GetMsg())
	}
	return msg
}

func (r *Room) getPropsMsg() []PropMsg{
	var msg []PropMsg
	for _, p := range r.props {
		msg = append(msg, p.GetMsg())
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
		r.loadProcess = r.loadProcess + rand.Intn(20)
		if r.loadProcess > 100 {
			r.loadProcess = 100
		}
		msg := &ws.Message{
			Event: events.LoadProcess,
			Data: encodeMsgData(LoadProcessData{
				Id: r.id,
				Process: r.loadProcess,
			}),
		}
		r.BroadcastMsg(msg)
		if r.loadProcess >= 100 {
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
}