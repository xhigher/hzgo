package logic

import (
	"context"
	"math"
	"math/rand"
	"time"
)

type MatchStatus int
const (
	MatchWait MatchStatus = 0 //等待到时开始
	MatchEnroll MatchStatus = 1 //到时开始，报名阶段
	MatchReady MatchStatus = 2 //人数已满，分配房间
	MatchOngoing MatchStatus = 3 //比赛进行中
	MatchPause MatchStatus = 4 //小局结束，赛中暂停
	MatchEnd MatchStatus = 5 //结束

	tickerDuration = 100*time.Millisecond
	readyDuration = 10 *time.Second
)

type MatchConfig struct {
	Id string `json:"id"`
	RoundCount int `json:"round_count"`
	RoomPlayerCount int `json:"room_player_count"`
	StartTime int64 `json:"start_time"`
}

type Match struct {
	id string
	playerCount int
	roomPlayerCount int
	roundCount int
	startTime time.Time
	nextTime time.Time
	curRound int
	status MatchStatus
	players []*Player
	rooms []*Room
	ticker *time.Ticker
	roomContext context.Context
	resultChan chan *RoomResult
}

func newMatch() *Match{
	m := &Match{
		ticker: time.NewTicker(tickerDuration),
		resultChan: make(chan *RoomResult, 100),
	}

	go func() {
		for {
			select {
			case <-m.ticker.C:
				m.HandleTick()
			}
		}
	}()

	go func() {
		for {
			var result *RoomResult
			select {
				case result = <-m.resultChan:
					m.HandleRoomResult(result)
			}
		}
	}()

	return m
}

func (m *Match) Restart(config *MatchConfig){
	if config == nil {
		return
	}
	if m.status != MatchWait {
		return
	}
	m.id = config.Id
	m.roundCount = config.RoundCount
	m.roomPlayerCount = config.RoomPlayerCount
	m.playerCount = int(math.Pow(float64(config.RoomPlayerCount), float64(config.RoundCount)))
	m.curRound = 0
	m.startTime = time.Unix(config.StartTime, 0)
	m.players = nil
	m.rooms = nil
}

func (m *Match) HandleTick(){
	switch m.status {
	case MatchWait:
		m.checkEnroll()
	case MatchReady:
		m.checkReady()
	case MatchOngoing:
	case MatchPause:
	case MatchEnd:
	}
}

func (m *Match) HandleRoomResult(result *RoomResult) {

}

func (m *Match) checkEnroll(){
	if time.Now().After(m.startTime) {
		m.status = MatchEnroll
	}
}

func (m *Match) checkReady(){
	if m.status == MatchReady && time.Now().After(m.nextTime) {
		m.status = MatchEnroll
		for _, r := range m.rooms {
			r.readyGo()
		}
	}
}

func (m *Match) JoinPlayer(player *Player){
	if m.status == MatchWait {

	}else if m.status == MatchReady || m.status == MatchOngoing || m.status == MatchPause {

	}else if m.status == MatchEnd {

	} else if m.status == MatchEnroll {
		for _,p := range m.players {
			if p.id == player.id {
				return
			}
		}
		m.players = append(m.players, player)
		if len(m.players) == m.playerCount {
			m.status = MatchReady
			m.curRound ++
			m.nextTime = time.Now().Add(readyDuration)
			m.startRound()
		}
	}else{
		return
	}
}

func (m *Match) startRound(){
	players := m.players
	rand.Shuffle(len(players), func(i, j int) {
		players[i], players[j] = players[j], players[i]
	})
	roomCount := m.playerCount / m.roomPlayerCount
	m.rooms = make([]*Room, roomCount)
	roomId := m.curRound * 100
	for i:=0; i<roomCount; i++ {
		m.rooms[i] = NewRoom(m, roomId+i, m.roomPlayerCount)
		for j:=i*m.roomPlayerCount; j<(i+1)*m.roomPlayerCount; j++ {
			m.rooms[i].JoinPlayer(m.players[j])
		}
	}
}

func (m *Match) checkRoundState(){

}