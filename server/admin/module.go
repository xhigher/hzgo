package admin

import (
	"context"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/xhigher/hzgo/consts"
)

type Module interface {
	Name() string
	Routers() []Router
}

type PlatformModuleHandler interface {
	Login(ctx context.Context, c *app.RequestContext)
	Logout(ctx context.Context, c *app.RequestContext)
	Profile(ctx context.Context, c *app.RequestContext)
	Menus(ctx context.Context, c *app.RequestContext)
}

type ModuleManager interface {
	Modules() []Module
	PlatformHandler() PlatformModuleHandler
	BaseController() *Controller
}

type PlatformModule struct {
	ctrl    *Controller
	handler PlatformModuleHandler
}

func (md PlatformModule) Name() string {
	return "platform"
}

func (md PlatformModule) Routers() []Router {
	return []Router{
		{
			Method:  consts.MethodPost,
			Name:    "login",
			NoAuth:  true,
			Handler: md.Login,
		},
		{
			Method:  consts.MethodGet,
			Name:    "logout",
			Handler: md.Logout,
		},
		{
			Method:  consts.MethodGet,
			Name:    "profile",
			Handler: md.Profile,
		},
		{
			Method:  consts.MethodGet,
			Name:    "menus",
			Handler: md.Menus,
		},
	}
}

func (md PlatformModule) Login(ctx context.Context, c *app.RequestContext) {
	md.handler.Login(ctx, c)
}

func (md PlatformModule) Logout(ctx context.Context, c *app.RequestContext) {
	md.handler.Logout(ctx, c)
}

func (md PlatformModule) Profile(ctx context.Context, c *app.RequestContext) {
	md.handler.Profile(ctx, c)
}

func (md PlatformModule) Menus(ctx context.Context, c *app.RequestContext) {
	md.handler.Menus(ctx, c)
}
