package logic

import (
	"github.com/xhigher/hzgo/games/live_game/events"
	"github.com/xhigher/hzgo/server/ws"
	"go.uber.org/atomic"
	"sync"
)

const (
	userMaxCount = 5000
)

type Handler struct {
	userList sync.Map
	userCount atomic.Uint32
}

func NewHandler() *Handler{
	return &Handler{
		userList: sync.Map{},
	}
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
			h.handleJoin(pipe.GetTag())

		}
	}

}


func (h *Handler) handleJoin(id string){

}
