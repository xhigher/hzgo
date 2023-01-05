package gateway

import (
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/xhigher/hzgo/server/gateway/middlewares"
)

type Controller struct {
	Name string
	Auth *middlewares.JWTAuth
}

func (ctrl Controller) Userid(c *app.RequestContext) string {
	return ctrl.Auth.GetAudience(c)
}

func (ctrl Controller) Token(c *app.RequestContext) string {
	return ctrl.Auth.GetToken(c)
}
