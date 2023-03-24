package logic

import (
	"context"
	"github.com/xhigher/hzgo/games/live_game/events"
	"github.com/xhigher/hzgo/logger"
	"github.com/xhigher/hzgo/server/ws"
	"github.com/xhigher/hzgo/utils"
	"math"
	"math/rand"
	"time"
)

type MatchStatus int
const (
	MatchWaiting MatchStatus = 0 //等待到时开始
	MatchEnrolling MatchStatus = 1 //到时开始，报名阶段
	MatchReadying MatchStatus = 2 //分配房间准备中
	MatchOngoing MatchStatus = 3 //比赛进行中
	MatchEnded MatchStatus = 4 //结束
	readyDuration = 5 *time.Second
	roundTimeoutDuration = 5*60 *time.Second
)

type MatchConfig struct {
	Id string `json:"id"`
	RoundCount int `json:"round_count"`
	RoomPlayerCount int `json:"room_player_count"`
	StartTime string `json:"start_time"`
	IntervalHours int64 `json:"interval_hours"`
}

type Match struct {
	config MatchConfig
	playerCount int
	startTime time.Time
	nextTime time.Time
	curRound int
	status MatchStatus
	players []*Player
	robots map[string]*Robot
	rooms []*Room
	roomContext context.Context
	resultChan chan *RoomResult
}

func newMatch(config MatchConfig) *Match{
	m := &Match{
		config: config,
		resultChan: make(chan *RoomResult, 100),
		playerCount: int(math.Pow(float64(config.RoomPlayerCount), float64(config.RoundCount))),
		startTime: utils.ParseYmdhms(config.StartTime),
		status: MatchWaiting,
	}

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

func (m *Match) Restart(){
	if m.status != MatchWaiting {
		return
	}
	if m.config.IntervalHours == 0 {
		return
	}

	m.status= MatchWaiting
	m.curRound = 0
	m.startTime = m.startTime.Add(time.Duration(m.config.IntervalHours) * time.Hour)
	m.players = nil
	m.rooms = nil
}

func (m *Match) HandleTick(){
	now := time.Now()
	logger.Infof("HandleTick: %v, %v", m.status, now.Format(utils.TimeFormatYMDHMS))
	switch m.status {
	case MatchWaiting:
		m.checkStarted(now)
	case MatchReadying:
		m.checkReady(now)
	case MatchOngoing:
		m.handleTimerEvent(now)
	case MatchEnded:
	}
}

func (m *Match) HandleRoomResult(result *RoomResult) {

}

func (m *Match) checkStarted(now time.Time){
	logger.Infof("checkStarted: %v, %v", now.Format(utils.TimeFormatYMDHMS),m.startTime.Format(utils.TimeFormatYMDHMS))
	if now.After(m.startTime) {
		m.status = MatchEnrolling
	}
}

func (m *Match) checkReady(now time.Time){
	if m.status == MatchReadying && now.After(m.nextTime) {
		m.status = MatchOngoing
		endTime := now.Add(roundTimeoutDuration)
		for _, r := range m.rooms {
			r.ReadyGo(endTime)
		}
	}
}

func (m *Match) handleTimerEvent(now time.Time){
	if m.status == MatchOngoing {
		endedCount := 0
		var winners []*Player
		for _, r := range m.rooms {
			if r.IsEnded() {
				endedCount ++
				winners = append(winners, r.winner)
			}
		}
		if endedCount == len(m.rooms) {
			if endedCount == 1 {
				m.status = MatchEnded
			}else{
				m.status = MatchReadying
				m.players = winners
				m.nextTime = time.Now().Add(readyDuration)
				m.startRound()
			}
		}

		for _, r := range m.rooms {
			r.handleTimerEvent(now)
		}
	}
}

func (m *Match) JoinRobot(robot *Robot) (err error){
	err = m.JoinPlayer(robot.Player)
	if err != nil {
		return
	}
	m.robots[robot.id] = robot
	return
}

func (m *Match) JoinPlayer(player *Player) (err error){
	if m.status == MatchWaiting {
		return errMatchWait
	}else if m.status == MatchReadying || m.status == MatchOngoing {
		return errMatchOngoing
	}else if m.status == MatchEnded {
		return errMatchEnd
	} else if m.status == MatchEnrolling {
		for _,p := range m.players {
			if p.id == player.id {
				return errPlayerJoined
			}
		}
		m.players = append(m.players, player)
		m.BroadcastMsg(&ws.Message{
			Event: events.JoinSuccess,
			Data: encodeMsgData(player.GetData()),
		})

		if len(m.players) == m.playerCount {
			m.status = MatchReadying
			m.nextTime = time.Now().Add(readyDuration)
			m.startRound()
		}
	}
	return nil
}

func (m *Match) startRound(){
	m.curRound ++
	players := m.players
	rand.Shuffle(len(players), func(i, j int) {
		players[i], players[j] = players[j], players[i]
	})
	roomCount := m.playerCount / m.config.RoomPlayerCount
	m.rooms = make([]*Room, roomCount)
	roomId := m.curRound * 100
	for i:=0; i<roomCount; i++ {
		room := NewRoom(roomId+i, m.config.RoomPlayerCount)
		for j:=i*m.config.RoomPlayerCount; j<(i+1)*m.config.RoomPlayerCount; j++ {
			player := m.players[j]
			if player.IsRobot() {
				room.JoinRobot(m.robots[player.id])
			}else{
				room.JoinPlayer(player)
			}
		}
		room.RoundStart()
		m.rooms[i] = room
	}
}

func (m *Match) JoinPlayerRobots(){
	for {
		time.Sleep(time.Duration(utils.RandInt64(500, 1500))*time.Millisecond)
		if err:= m.JoinRobot(GeRobot()); err != nil {
			return
		}
	}
}

func (m *Match) BroadcastMsg(msg *ws.Message){
	for _,p := range m.players {
		p.SendMsg(msg)
	}
}