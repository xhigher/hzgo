package auth

import (
	"context"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/xhigher/hzgo/router"
)

type HzgoAuthController interface {
	Register(ctx context.Context, c *app.RequestContext)
	Login(ctx context.Context, c *app.RequestContext)
	Logout(ctx context.Context, c *app.RequestContext)
	Refresh(ctx context.Context, c *app.RequestContext)
}

func (mw *HzgoJWTMiddleware) Routers() []router.Router{
	return []router.Router{
		{
			Method: router.MethodPost,
			Version: 1,
			Path: "register",
			Auth: false,
			Handler: mw.Controller.Register,
		},
		{
			Method: router.MethodPost,
			Version: 1,
			Path: "login",
			Auth: false,
			Handler: mw.Controller.Login,
		},
		{
			Method: router.MethodGet,
			Version: 1,
			Path: "logout",
			Auth: true,
			Handler: mw.Controller.Logout,
		},
		{
			Method: router.MethodGet,
			Version: 1,
			Path: "token",
			Auth: true,
			Handler: mw.Controller.Refresh,
		},
	}
}