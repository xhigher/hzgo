package admin

import (
	"context"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/xhigher/hzgo/consts"
)

type Module interface {
	Name() string
	Routers() []Router
	Roles() []string
}

type PlatformModuleHandler interface {
	Login(ctx context.Context, c *app.RequestContext)
	Logout(ctx context.Context, c *app.RequestContext)
	Profile(ctx context.Context, c *app.RequestContext)
}

type ModuleManager interface {
	Modules() []Module
	PlatformHandler() PlatformModuleHandler
	BaseController() *Controller
}

type PlatformModule struct {
	ctrl *Controller
	handler PlatformModuleHandler
}

func (md PlatformModule) Name() string{
	return "platform"
}

func (md PlatformModule) Routers() []Router{
	return []Router{
		{
			Method:  consts.MethodPost,
			Path:    "login",
			NoAuth: true,
			Handler: md.Login,
		},
		{
			Method:  consts.MethodGet,
			Path:    "logout",
			Handler: md.Logout,
		},
		{
			Method:  consts.MethodGet,
			Path:    "profile",
			Handler: md.Profile,
		},
	}
}

func (md PlatformModule) Login(ctx context.Context, c *app.RequestContext){
	md.handler.Login(ctx, c)
}

func (md PlatformModule) Logout(ctx context.Context, c *app.RequestContext) {
	md.handler.Logout(ctx, c)
}

func (md PlatformModule) Profile(ctx context.Context, c *app.RequestContext) {
	md.handler.Profile(ctx, c)
}
