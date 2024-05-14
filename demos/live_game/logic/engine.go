package logic

import (
	"errors"
	"github.com/xhigher/hzgo/games/live_game/events"
	"github.com/xhigher/hzgo/games/live_game/maps"
	"github.com/xhigher/hzgo/games/live_game/model/store"
	"github.com/xhigher/hzgo/logger"
	"github.com/xhigher/hzgo/server/game"
	"github.com/xhigher/hzgo/utils"
	"go.uber.org/atomic"
	"sync"
	"time"
)

const (
	userMaxCount   = 5000
	tickerDuration = 50 * time.Millisecond
	playerStepTime = 300
)

var (
	lock          sync.Once
	defaultEngine *Engine
)

type Engine struct {
	ticker      *time.Ticker
	playerList  sync.Map
	playerCount atomic.Uint32
	match       *Match
}

func StartEngine(config MatchConfig) {
	lock.Do(func() {
		maps.Init()

		defaultEngine = &Engine{
			ticker:      time.NewTicker(tickerDuration),
			playerList:  sync.Map{},
			playerCount: atomic.Uint32{},
			match:       newMatch(config),
		}
		defaultEngine.startTicker()

		InitRobots(10)
	})
}

func (e *Engine) startTicker() {
	go func() {
		for {
			select {
			case <-e.ticker.C:
				e.match.HandleTick()
			}
		}
	}()
}

func (e *Engine) AddPlayer(player *Player) {
	logger.Infof("AddPlayer: %v", player.GetData())
	e.playerList.Store(player.id, player)
}

func (e *Engine) GetPlayer(id string) *Player {
	if player, ok := e.playerList.Load(id); ok {
		logger.Infof("GetPlayer: %v", utils.JSONString(player))
		return player.(*Player)
	}
	return nil
}

func (e *Engine) DeletePlayer(id string) {
	logger.Infof("DeletePlayer: %v", id)
	e.playerList.Delete(id)
}

func (e *Engine) JoinMatch(id string, data *JoinData) {
	logger.Infof("JoinMatch: %v", id)
	player := e.GetPlayer(id)
	if player == nil {
		return
	}
	player.skin = data.Skin
	player.bubbleColor = data.BubbleColor
	err := e.match.JoinPlayer(player)
	if err != nil {
		player.pipe.SendMessage(&game.Message{
			Event: events.JoinError,
			Info:  err.Error(),
		})
		return
	}

	e.match.JoinPlayerRobots()
}

func (e *Engine) HandleRoomEvent(id, event string, msg *game.Message) {
	player := e.GetPlayer(id)
	if player == nil {
		return
	}
	if player.room == nil {
		return
	}

	switch msg.Event {
	case events.LoadProcess:
		e.handleLoadProcess(player, msg)
	}

}

func (e *Engine) handleLoadProcess(player *Player, msg *game.Message) {
	data := &LoadProcessData{}
	err := decodeMsgData(msg.Data, data)
	if err != nil {
		err = errors.New("数据格式错误")
		return
	}
	player.room.BroadcastMsg(msg)
}

func (e *Engine) UserLogin(pipe *game.Pipe, data *LoginData) (player *Player, err error) {
	if e.playerCount.Load() >= userMaxCount {
		err = errors.New("连接人数已满")
		return
	}
	if e.playerCount.Inc() > userMaxCount {
		defaultEngine.playerCount.Dec()
		err = errors.New("连接人数已满")
		return
	}

	user, typ, err := e.checkToken(data.Id, data.Token)
	if err != nil {
		err = errors.New("服务器错误")
		return
	}
	if typ != errOK {
		info := "登录失败"
		if typ == errUserNull {
			info = "用户不存在"
		} else if typ == errTokenInvalid {
			info = "token验证失败！请重新进入游戏！"
		}
		err = errors.New(info)
		return
	}

	player = NewPlayer(pipe, user, PlayerHuman)

	e.AddPlayer(player)

	return
}

func (e *Engine) checkToken(id, token string) (user store.UserInfo, typ errType, err error) {
	ok, user := GetUser(id)
	if !ok {
		typ = errUserNull
		return
	}
	return
}
