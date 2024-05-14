package biz

import (
	"context"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/xhigher/hzgo/consts"
	"github.com/xhigher/hzgo/defines"
	logic "github.com/xhigher/hzgo/demo/admin/logic/user"
	"github.com/xhigher/hzgo/demo/admin/rbac"
	"github.com/xhigher/hzgo/logger"
	"github.com/xhigher/hzgo/server/admin"
)

type UserModule struct {
	ctrl *admin.Controller
}

func (md UserModule) Name() string {
	return rbac.BizUser
}

func (md UserModule) Roles() []string {
	return []string{admin.RoleMaintainer, admin.RoleOperator}
}

func (md UserModule) Routers() []admin.Router {
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
			Method:  consts.MethodGet,
			Path:    "profile",
			Handler: md.Profile,
		},
		{
			Method:  consts.MethodPost,
			Path:    "start",
			Handler: md.Start,
		},
		{
			Method:  consts.MethodPost,
			Path:    "Stop",
			Handler: md.Stop,
		},
	}
}

func (md UserModule) Info(ctx context.Context, c *app.RequestContext) {
	resp := md.ctrl.Resp(c)
	params := defines.UseridReq{}
	if err := c.Bind(&params); err != nil {
		resp.ReplyErrorParam()
		return
	}

	userInfo, be := logic.GetUser(params.Userid)
	if be != nil {
		logger.Errorf("error: %v", be.String())
		resp.ReplyErr(be.ToResp())
		return
	}

	resp.ReplyData(userInfo)
}

func (md UserModule) List(ctx context.Context, c *app.RequestContext) {
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

	total, userList, err := logic.GetUserList(params.Status, params.Offset, params.Limit)
	if err != nil {
		resp.ReplyErrorInternal()
		return
	}

	resp.ReplyData(defines.PageData{
		Total:  int32(total),
		Offset: params.Offset,
		Limit:  params.Limit,
		Data:   userList,
	})
}

func (md UserModule) Profile(ctx context.Context, c *app.RequestContext) {
	resp := md.ctrl.Resp(c)
	params := defines.UseridReq{}
	if err := c.Bind(&params); err != nil {
		resp.ReplyErrorParam()
		return
	}

	userInfo, err := logic.GetUser(params.Userid)
	if err != nil {
		resp.ReplyErrorInternal()
		return
	}

	resp.ReplyData(userInfo)
}

func (md UserModule) Start(ctx context.Context, c *app.RequestContext) {
	resp := md.ctrl.Resp(c)
	params := defines.UseridReq{}
	if err := c.Bind(&params); err != nil {
		resp.ReplyErrorParam()
		return
	}

	err := logic.StartUser(params.Userid)
	if err != nil {
		resp.ReplyErrorInternal()
		return
	}

	resp.ReplyOK()
}

func (md UserModule) Stop(ctx context.Context, c *app.RequestContext) {
	resp := md.ctrl.Resp(c)
	params := defines.UseridReq{}
	if err := c.Bind(&params); err != nil {
		resp.ReplyErrorParam()
		return
	}

	err := logic.StopUser(params.Userid)
	if err != nil {
		resp.ReplyErrorInternal()
		return
	}

	resp.ReplyOK()
}
