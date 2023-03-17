package logic

import (
	"errors"
	"github.com/xhigher/hzgo/games/live_game/events"
	"github.com/xhigher/hzgo/server/ws"
)

type Handler struct {

}

func NewHandler() *Handler{
	return &Handler{}
}

func (h *Handler) HandleMessage(pipe *ws.Pipe, msg *ws.Message) {
	if msg.Event == events.Login {
		h.handleLogin(pipe, msg)
		return
	}else{
		if len(pipe.GetTag()) == 0 {
			pipe.SendMessage(&ws.Message{
				Event: events.LoginError,
				Info:  "请先登录",
			})
			return
		}
		switch msg.Event {
		case events.Join:
			h.handleJoin(pipe.GetTag(), msg)
		}
	}
}



func (h *Handler) handleJoin(id string, msg *ws.Message){
	data := &JoinData{}
	err := decodeMsgData(msg.Data, data)
	if err != nil {
		err = errors.New("数据格式错误")
		return
	}
	defaultEngine.JoinMatch(id, data)

}

type errType int

const (
	errOK = 0
	errUserNull = 1
	errTokenInvalid = 2
)

type UserInfo struct {
	Id string `json:"id"`
	Name string `json:"name"`
	Avatar string `json:"avatar"`
}

func (h *Handler) handleLogin(pipe *ws.Pipe, msg *ws.Message) {
	var err error
	defer func() {
		if err != nil {
			pipe.SendMessage(&ws.Message{
				Event: events.LoginError,
				Info:  err.Error(),
			})
			pipe.Close()
		}
	}()
	if len(pipe.GetTag()) > 0 {
		err = errors.New("您已登录")
		return
	}

	data := &LoginData{}
	err = decodeMsgData(msg.Data, data)
	if err != nil {
		err = errors.New("数据格式错误")
		return
	}

	player, err := defaultEngine.UserLogin(pipe, data)
	if err != nil {
		return
	}

	pipe.SetTag(player.id)

	pipe.SendMessage(&ws.Message{
		Event: events.LoginSuccess,
		Data:encodeMsgData(player.GetData()),
	})
}