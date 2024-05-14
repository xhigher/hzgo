package controller

import (
	"context"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/xhigher/hzgo/consts"
	"github.com/xhigher/hzgo/defines"
	"github.com/xhigher/hzgo/demo/api"
	"github.com/xhigher/hzgo/server/gateway"
	"github.com/xhigher/hzgo/server/gateway/middlewares"
)

type Controller struct {
	*gateway.Controller
}

func (ctrl Controller) Name() string {
	return ctrl.Controller.Name
}

func (ctrl Controller) Routers() []gateway.Router {
	return []gateway.Router{
		{
			Method:  consts.MethodPost,
			Version: 1,
			Path:    "register",
			Sign:    true,
			Auth:    false,
			Handler: ctrl.Register,
		},
		{
			Method:  consts.MethodPost,
			Version: 1,
			Path:    "login",
			Sign:    true,
			Auth:    false,
			Handler: ctrl.Login,
		},
		{
			Method:  consts.MethodGet,
			Version: 1,
			Path:    "logout",
			Sign:    true,
			Auth:    true,
			Handler: ctrl.Logout,
		},
		{
			Method:  consts.MethodGet,
			Version: 1,
			Path:    "renewal",
			Sign:    true,
			Auth:    true,
			Handler: ctrl.Renewal,
		},
		{
			Method:  consts.MethodGet,
			Version: 1,
			Path:    "profile",
			Sign:    true,
			Auth:    true,
			Handler: ctrl.Profile,
		},
	}
}

func New(name string) Controller {
	return Controller{
		&gateway.Controller{
			Name: name,
		},
	}
}

func NewWithAuth(name string, auth *middlewares.JWTAuth) Controller {
	auth.CheckTokenFunc = func(ctx context.Context, c *app.RequestContext, claims *middlewares.AuthClaims) bool {
		result := api.User().TokenCheck(defines.TokenCheckReq{
			Userid:  claims.Subject,
			TokenId: claims.TokenId,
		})
		if result.NotOK() {
			return false
		}
		return true
	}
	return Controller{
		&gateway.Controller{
			Name: name,
			Auth: auth,
		},
	}
}

func (ctrl Controller) Test(ctx context.Context, c *app.RequestContext) {
	resp := ctrl.Resp(c)

	resp.ReplyOK()
}
