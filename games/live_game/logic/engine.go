package logic

import (
	"errors"
	"github.com/xhigher/hzgo/server/ws"
	"github.com/xhigher/hzgo/utils"
	"go.uber.org/atomic"
	"sync"
	"time"
)

const (
	userMaxCount = 5000
)

var (
	lock sync.Once
	defaultEngine *Engine
)

type Engine struct {
	ticker   *time.Ticker
	playerList sync.Map
	playerCount atomic.Uint32
	match *Match
}


func StartEngine(){
	lock.Do(func() {
		defaultEngine = &Engine{
			ticker: time.NewTicker(tickerDuration),
			playerList: sync.Map{},
			playerCount: atomic.Uint32{},
			match: newMatch(),
		}
		initRobots(10)
	})
}

func (e *Engine) AddPlayer(user *UserInfo, pipe *ws.Pipe) *Player{
	player := &Player{
		id: user.Id,
		name: user.Name,
		avatar: user.Avatar,
		pipe: pipe,
	}

	e.playerList.Store(player.id, player)

	return player
}

func (e *Engine) GetPlayer(id string) *Player{
	if player, ok := e.playerList.Load(id); ok {
		return player.(*Player)
	}
	return nil
}

func (e *Engine) DeletePlayer(id string){
	e.playerList.Delete(id)
}

func (e *Engine) JoinMatch(id string, data *JoinData){
	player := e.GetPlayer(id)
	if player == nil {
		return
	}
	player.bubbleColor = data.BombColor
	player.role = data.Role
	e.match.JoinPlayer(player)
}

func (e *Engine) UserLogin(pipe *ws.Pipe, data *LoginData) (player *Player, err error) {
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
		}else if typ == errTokenInvalid {
			info = "token验证失败！请重新进入游戏！"
		}
		err = errors.New(info)
		return
	}

	player = &Player{
		id: user.Id,
		name: user.Name,
		avatar: user.Avatar,
		pipe: pipe,
	}
	e.playerList.Store(player.id, player)

	return
}

func (e *Engine) checkToken(id, token string) (user *UserInfo, typ errType, err error){
	user = &UserInfo{
		Id: id,
		Name: utils.RandString(20),
		Avatar: "",
	}
	return
}