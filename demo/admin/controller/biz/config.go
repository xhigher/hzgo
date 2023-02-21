package biz

import (
	"context"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/xhigher/hzgo/consts"
	"github.com/xhigher/hzgo/defines"
	"github.com/xhigher/hzgo/demo/admin/rbac"
	"github.com/xhigher/hzgo/server/admin"
)

type ConfigModule struct {
	ctrl *admin.Controller
}

func (md ConfigModule) Name() string{
	return rbac.BizConfig
}

func (md ConfigModule) Roles() []string{
	return []string{admin.RoleMaintainer, admin.RoleOperator}
}

func (md ConfigModule) Routers() []admin.Router{
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
			Path:    "change_status",
			Handler: md.ChangeStatus,
		},
	}
}

type TraceLogPageReq struct {
	Uid string `form:"uid" json:"uid" query:"uid"`
	Offset int32 `form:"offset" json:"offset" query:"offset"`
	Limit int32  `form:"limit" json:"limit" query:"limit"`
}
func (md ConfigModule) Info(ctx context.Context, c *app.RequestContext) {
	resp := md.ctrl.Resp(c)

	resp.ReplyData(nil)
}

func (md ConfigModule) List(ctx context.Context, c *app.RequestContext) {
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

func (md ConfigModule) Create(ctx context.Context, c *app.RequestContext) {
	resp := md.ctrl.Resp(c)

	resp.ReplyOK()
}

func (md ConfigModule) ChangeStatus(ctx context.Context, c *app.RequestContext) {
	resp := md.ctrl.Resp(c)
	params := defines.ChangeStatusReq{}
	if err := c.Bind(&params); err != nil {
		resp.ReplyErrorParam()
		return
	}
	if !defines.CheckChangeStatus(params.Status) {
		resp.ReplyErrorParam2("status")
		return
	}


	resp.ReplyOK()
}

