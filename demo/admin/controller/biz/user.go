package biz

import (
	"context"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/xhigher/hzgo/consts"
	"github.com/xhigher/hzgo/defines"
	logic "github.com/xhigher/hzgo/demo/admin/logic/user"
	"github.com/xhigher/hzgo/demo/admin/rbac"
	"github.com/xhigher/hzgo/logger"
	"github.com/xhigher/hzgo/resp"
	"github.com/xhigher/hzgo/server/admin"
)

type UserModule struct {
	ctrl *admin.Controller
}

func (md UserModule) Name() string{
	return rbac.BizUser
}

func (md UserModule) Routers() []admin.Router{
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
	params := defines.UseridReq{}
	if err := c.Bind(&params); err != nil {
		resp.ReplyErrorParam(c)
		return
	}

	userInfo, be := logic.GetUser(params.Userid)
	if be != nil {
		logger.Errorf("error: %v", be.String())
		resp.ReplyErr(c, be.ToResp())
		return
	}

	resp.ReplyData(c, userInfo)
}

func (md UserModule) List(ctx context.Context, c *app.RequestContext) {
	params := defines.StatusPageReq{}
	if err := c.Bind(&params); err != nil {
		resp.ReplyErrorParam(c)
		return
	}
	if !defines.CheckPageLimit(params.Limit) {
		resp.ReplyErrorParam2(c, "limit")
		return
	}

	total, userList, err := logic.GetUserList(params.Status, params.Offset, params.Limit)
	if err != nil {
		resp.ReplyErrorInternal(c)
		return
	}

	resp.ReplyData(c, defines.PageData{
		Total: int32(total),
		Offset: params.Offset,
		Limit: params.Limit,
		Data: userList,
	})
}

func (md UserModule) Profile(ctx context.Context, c *app.RequestContext) {
	params := defines.UseridReq{}
	if err := c.Bind(&params); err != nil {
		resp.ReplyErrorParam(c)
		return
	}

	userInfo, err := logic.GetUser(params.Userid)
	if err != nil {
		resp.ReplyErrorInternal(c)
		return
	}

	resp.ReplyData(c, userInfo)
}


func (md UserModule) Start(ctx context.Context, c *app.RequestContext) {
	params := defines.UseridReq{}
	if err := c.Bind(&params); err != nil {
		resp.ReplyErrorParam(c)
		return
	}

	err := logic.StartUser(params.Userid)
	if err != nil {
		resp.ReplyErrorInternal(c)
		return
	}

	resp.ReplyOK(c)
}

func (md UserModule) Stop(ctx context.Context, c *app.RequestContext) {
	params := defines.UseridReq{}
	if err := c.Bind(&params); err != nil {
		resp.ReplyErrorParam(c)
		return
	}

	err := logic.StopUser(params.Userid)
	if err != nil {
		resp.ReplyErrorInternal(c)
		return
	}

	resp.ReplyOK(c)
}
