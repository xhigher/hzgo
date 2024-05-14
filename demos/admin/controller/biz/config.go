package biz

import (
	"context"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/xhigher/hzgo/bizerr"
	"github.com/xhigher/hzgo/consts"
	"github.com/xhigher/hzgo/defines"
	logic "github.com/xhigher/hzgo/demo/admin/logic/misc"
	"github.com/xhigher/hzgo/demo/admin/rbac"
	"github.com/xhigher/hzgo/logger"
	"github.com/xhigher/hzgo/server/admin"
)

type ConfigModule struct {
	ctrl *admin.Controller
}

func (md ConfigModule) Name() string {
	return rbac.BizConfig
}

func (md ConfigModule) Roles() []string {
	return []string{admin.RoleMaintainer, admin.RoleOperator}
}

func (md ConfigModule) Routers() []admin.Router {
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
			Path:    "save",
			Handler: md.Save,
		},
		{
			Method:  consts.MethodPost,
			Path:    "delete",
			Handler: md.Delete,
		},
		{
			Method:  consts.MethodPost,
			Path:    "change_status",
			Handler: md.ChangeStatus,
		},
	}
}

func (md ConfigModule) Info(ctx context.Context, c *app.RequestContext) {
	resp := md.ctrl.Resp(c)
	params := defines.CommonIdReq{}
	if err := c.Bind(&params); err != nil {
		resp.ReplyErrorParam()
		return
	}
	data, be := logic.GetConfigInfo(params.Id)
	if be != nil {
		logger.Errorf("error: %v", be.String())
		resp.ReplyErr(be.ToResp())
		return
	}
	resp.ReplyData(data)
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

	total, data, be := logic.GetConfigList(params.Status, params.Offset, params.Limit)
	if be != nil {
		logger.Errorf("error: %v", be.String())
		resp.ReplyErr(be.ToResp())
		return
	}

	resp.ReplyData(defines.PageData{
		Total:  int32(total),
		Offset: params.Offset,
		Limit:  params.Limit,
		Data:   data,
	})
}

type ConfigSaveReq struct {
	Id      string `form:"id" json:"id" query:"id"`
	Name    string `form:"name" json:"name" query:"name"`
	Items   string `form:"items" json:"items" query:"items"`
	Static  bool   `form:"static" json:"static" query:"static"`
	Filters string `form:"filters" json:"filters" query:"filters"`
}

func (md ConfigModule) Save(ctx context.Context, c *app.RequestContext) {
	resp := md.ctrl.Resp(c)
	params := ConfigSaveReq{}
	if err := c.Bind(&params); err != nil {
		resp.ReplyErrorParam()
		return
	}
	reload, be := logic.SaveConfigInfo(params.Id, params.Name, params.Items, params.Static, params.Filters)
	if be != nil {
		logger.Errorf("error: %v", be.String())
		resp.ReplyErr(be.ToResp())
		return
	}
	if reload {

	}
	resp.ReplyOK()
}

func (md ConfigModule) Delete(ctx context.Context, c *app.RequestContext) {
	resp := md.ctrl.Resp(c)
	params := defines.CommonIdReq{}
	if err := c.Bind(&params); err != nil {
		resp.ReplyErrorParam()
		return
	}
	be := logic.DeleteConfigInfo(params.Id)
	if be != nil {
		logger.Errorf("error: %v", be.String())
		resp.ReplyErr(be.ToResp())
		return
	}
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

	var be *bizerr.Error
	if params.Status == consts.StatusOnline {
		be = logic.SetConfigOnline(params.Id)
	} else if params.Status == consts.StatusOffline {
		be = logic.SetConfigOffline(params.Id)
	}
	if be != nil {
		logger.Errorf("error: %v", be.String())
		resp.ReplyErr(be.ToResp())
		return
	}

	resp.ReplyOK()
}
