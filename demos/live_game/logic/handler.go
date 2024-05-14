package logic

import (
	"errors"
	"github.com/xhigher/hzgo/games/live_game/events"
	"github.com/xhigher/hzgo/logger"
	"github.com/xhigher/hzgo/server/game"
	"github.com/xhigher/hzgo/utils"
)

type Handler struct {
}

func NewHandler() *Handler {
	return &Handler{}
}

func (h *Handler) HandleMessage(pipe *game.Pipe, msg *game.Message) {
	logger.Infof("HandleMessage event: %v, data: %v", msg.Event, utils.JSONString(msg.Data))
	if msg.Event == events.Login {
		h.handleLogin(pipe, msg)
		return
	} else {
		if len(pipe.GetTag()) == 0 {
			pipe.SendMessage(&game.Message{
				Event: events.LoginError,
				Info:  "请先登录",
			})
			return
		}
		logger.Infof("HandleMessage event: %v, data: %v", msg.Event, utils.JSONString(msg.Data))
		switch msg.Event {
		case events.Join:
			h.handleJoin(pipe.GetTag(), msg)
		case events.LoadProcess:
			h.HandleRoomEvent(pipe.GetTag(), msg.Event, msg)
		}
	}
}

func (h *Handler) handleJoin(id string, msg *game.Message) {
	data := &JoinData{}
	err := decodeMsgData(msg.Data, data)
	if err != nil {
		err = errors.New("数据格式错误")
		return
	}
	defaultEngine.JoinMatch(id, data)

}

func (h *Handler) HandleRoomEvent(id, event string, msg *game.Message) {
	defaultEngine.HandleRoomEvent(id, event, msg)

}

type errType int

const (
	errOK           = 0
	errUserNull     = 1
	errTokenInvalid = 2
)

func (h *Handler) handleLogin(pipe *game.Pipe, msg *game.Message) {
	var err error
	defer func() {
		if err != nil {
			pipe.SendMessage(&game.Message{
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

	pipe.SendMessage(&game.Message{
		Event: events.LoginSuccess,
		Data:  encodeMsgData(player.GetData()),
	})
}
