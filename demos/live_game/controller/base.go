package controller

import (
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/xhigher/hzgo/consts"
	"github.com/xhigher/hzgo/resp"
	"github.com/xhigher/hzgo/server/game"
)

type Controller struct {
}

func New() Controller {
	return Controller{}
}

func (c Controller) Routers() []game.Router {
	return []game.Router{
		{
			Method:  consts.MethodGet,
			Path:    "wechat_login",
			Handler: c.WechatLogin,
		},
		{
			Method:  consts.MethodGet,
			Path:    "user_info",
			Handler: c.UserInfo,
		},
	}
}

func (c Controller) Resp(ctx *app.RequestContext) resp.Responder {
	return resp.Responder{Ctx: ctx, LogOut: true}
}

func (c Controller) checkToken(id string, token string) bool {

	return true
}
