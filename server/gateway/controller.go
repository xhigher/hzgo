package gateway

import (
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/xhigher/hzgo/defines"
	"github.com/xhigher/hzgo/resp"
	"github.com/xhigher/hzgo/server/gateway/middlewares"
)

type Controller struct {
	Name string
	Auth *middlewares.JWTAuth
	LogSaver resp.TraceLogSaver
}

func (ctrl Controller) Resp(c *app.RequestContext) resp.Responder {
	return resp.Responder{Ctx:c, LogSaver: ctrl.LogSaver, LogOut: true}
}

func (ctrl Controller) Userid(c *app.RequestContext) string {
	return middlewares.GetAudience(c)
}

func (ctrl Controller) Token(c *app.RequestContext) string {
	return middlewares.GetToken(c)
}

func (ctrl Controller) BaseParams(c *app.RequestContext) *defines.BaseParams {
	return middlewares.GetBaseParams(c)
}