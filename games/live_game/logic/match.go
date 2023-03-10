package logic

import (
	"context"
	"github.com/xhigher/hzgo/server/ws"
	"github.com/xhigher/hzgo/utils"
	"math"
	"math/rand"
	"sync"
	"time"
)

type MatchStatus int
const (
	MatchWait MatchStatus = 0 //等待到时开始
	MatchReady MatchStatus = 1 //到时开始，可加入
	MatchOngoing MatchStatus = 2 //比赛进行中
	MatchPause MatchStatus = 3 //小局结束，赛中暂停
	MatchEnd MatchStatus = 4 //结束

	tickerDuration = 100*time.Millisecond
)

var (
	lock sync.Once
	defaultMatch *Match
)

func StartMatch(){
	lock.Do(func() {
		defaultMatch = newMatch()
	})
}

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
	startTime int64
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
	m.startTime = config.StartTime
	m.players = nil
	m.rooms = nil
}

func (m *Match) HandleTick(){
	switch m.status {
	case MatchWait:
		m.checkReady()
	case MatchReady:
	case MatchOngoing:
	case MatchPause:
	case MatchEnd:
	}
}

func (m *Match) HandleRoomResult(result *RoomResult) {

}

func (m *Match) checkReady(){
	if m.startTime > 0 && m.startTime<=utils.NowTime() {
		m.status = MatchReady
	}
}

func (m *Match) AddPlayer(user *UserInfo, pipe *ws.Pipe) *Player{
	player := &Player{
		id: user.Id,
		name: user.Name,
		avatar: user.Avatar,
		pipe: pipe,
	}

	return player
}

func (m *Match) JoinUser(player *Player){
	if m.status == MatchWait {

	}else if m.status == MatchOngoing || m.status == MatchPause {

	}else if m.status == MatchEnd {

	} else if m.status == MatchReady {
		for _,p := range m.players {
			if p.id == player.id {
				return
			}
		}
		m.players = append(m.players, player)
		if len(m.players) == m.playerCount {
			m.status = MatchPause
			m.curRound ++
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