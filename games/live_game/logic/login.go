package logic

import (
	"errors"
	"github.com/xhigher/hzgo/games/live_game/events"
	"github.com/xhigher/hzgo/server/ws"
	"github.com/xhigher/hzgo/utils"
)

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
	if h.userCount.Load() >= userMaxCount {
		err = errors.New("连接人数已满")
		return
	}
	if h.userCount.Inc() > userMaxCount {
		h.userCount.Dec()
		err = errors.New("连接人数已满")
		return
	}

	data := LoginData{}
	err = decodeMsgData(msg.Data, &data)
	if err != nil {
		err = errors.New("数据格式错误")
		return
	}
	user, typ, err := h.checkToken(data.Id, data.Token)
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
	pipe.SetTag(user.Id)

	player := defaultMatch.AddPlayer(user, pipe)

	pipe.SendMessage(&ws.Message{
		Event: events.LoginSuccess,
		Data:encodeMsgData(player.GetMsg()),
	})
}

func (h *Handler) checkToken(id, token string) (user *UserInfo, typ errType, err error){
	user = &UserInfo{
		Id: id,
		Name: utils.RandString(20),
		Avatar: "",
	}
	return
}

