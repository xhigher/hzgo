package biz

import (
	"context"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/xhigher/hzgo/consts"
	"github.com/xhigher/hzgo/defines"
	"github.com/xhigher/hzgo/demo/admin/rbac"
	"github.com/xhigher/hzgo/server/admin"
)

type BannerModule struct {
	ctrl *admin.Controller
}

func (md BannerModule) Name() string{
	return rbac.BizBanner
}

func (md BannerModule) Roles() []string{
	return []string{admin.RoleMaintainer, admin.RoleOperator}
}

func (md BannerModule) Routers() []admin.Router{
	return []admin.Router{
		{
			Method:  consts.MethodGet,
			Path:    "info",
			Handler: md.Info,
		},
		{
			Method:  consts.MethodGet,
			Path:    "list",
			Handler: md.List,
		},
		{
			Method:  consts.MethodPost,
			Path:    "create",
			Handler: md.Create,
		},
		{
			Method:  consts.MethodPost,
			Path:    "set_online",
			Handler: md.SetOnline,
		},
		{
			Method:  consts.MethodPost,
			Path:    "set_offline",
			Handler: md.SetOffline,
		},
	}
}


func (md BannerModule) Info(ctx context.Context, c *app.RequestContext) {
	resp := md.ctrl.Resp(c)

	resp.ReplyData(nil)
}

func (md BannerModule) List(ctx context.Context, c *app.RequestContext) {
	resp := md.ctrl.Resp(c)
	params := defines.StatusPageReq{}
	if err := c.Bind(&params); err != nil {
		resp.ReplyErrorParam()
		return
	}
	if !defines.CheckPageLimit(params.Limit) {
		resp.ReplyErrorParam2("limit")
		return
	}

	resp.ReplyData(defines.PageData{
		Total: int32(0),
		Offset: params.Offset,
		Limit: params.Limit,
		Data: nil,
	})
}

func (md BannerModule) Create(ctx context.Context, c *app.RequestContext) {
	resp := md.ctrl.Resp(c)

	resp.ReplyOK()
}

func (md BannerModule) SetOnline(ctx context.Context, c *app.RequestContext) {
	resp := md.ctrl.Resp(c)

	resp.ReplyOK()
}

func (md BannerModule) SetOffline(ctx context.Context, c *app.RequestContext) {
	resp := md.ctrl.Resp(c)

	resp.ReplyOK()
}
