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

type BannerModule struct {
	ctrl *admin.Controller
}

func (md BannerModule) Name() string {
	return rbac.BizBanner
}

func (md BannerModule) Roles() []string {
	return []string{admin.RoleMaintainer, admin.RoleOperator}
}

func (md BannerModule) Routers() []admin.Router {
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

func (md BannerModule) Info(ctx context.Context, c *app.RequestContext) {
	resp := md.ctrl.Resp(c)
	params := defines.CommonIdReq{}
	if err := c.Bind(&params); err != nil {
		resp.ReplyErrorParam()
		return
	}
	data, be := logic.GetBannerInfo(params.IntId())
	if be != nil {
		logger.Errorf("error: %v", be.String())
		resp.ReplyErr(be.ToResp())
		return
	}
	resp.ReplyData(data)
}

type BannerListReq struct {
	defines.StatusPageReq
	Site string `form:"site" json:"site" query:"site"`
}

func (md BannerModule) List(ctx context.Context, c *app.RequestContext) {
	resp := md.ctrl.Resp(c)
	params := BannerListReq{}
	if err := c.Bind(&params); err != nil {
		resp.ReplyErrorParam()
		return
	}
	if !defines.CheckPageLimit(params.Limit) {
		resp.ReplyErrorParam2("limit")
		return
	}

	total, data, be := logic.GetBannerList(params.Site, params.Status, params.Offset, params.Limit)
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

type BannerSaveReq struct {
	Id   int32  `form:"id" json:"id" query:"id"`
	Site string `form:"site" json:"site" query:"site"`
	Type int32  `form:"type" json:"type" query:"type"`
	Name string `form:"name" json:"name" query:"name"`
	Img  string `form:"img" json:"img" query:"img"`
	Data string `form:"data" json:"data" query:"data"`
}

func (md BannerModule) Save(ctx context.Context, c *app.RequestContext) {
	resp := md.ctrl.Resp(c)
	params := BannerSaveReq{}
	if err := c.Bind(&params); err != nil {
		resp.ReplyErrorParam()
		return
	}
	reload, be := logic.SaveBannerInfo(params.Id, params.Site, params.Type, params.Name, params.Img, params.Data)
	if be != nil {
		logger.Errorf("error: %v", be.String())
		resp.ReplyErr(be.ToResp())
		return
	}
	if reload {

	}
	resp.ReplyOK()
}

func (md BannerModule) Delete(ctx context.Context, c *app.RequestContext) {
	resp := md.ctrl.Resp(c)
	params := defines.CommonIdReq{}
	if err := c.Bind(&params); err != nil {
		resp.ReplyErrorParam()
		return
	}
	be := logic.DeleteBannerInfo(params.IntId())
	if be != nil {
		logger.Errorf("error: %v", be.String())
		resp.ReplyErr(be.ToResp())
		return
	}
	resp.ReplyOK()
}

func (md BannerModule) ChangeStatus(ctx context.Context, c *app.RequestContext) {
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
		be = logic.SetBannerOnline(params.IntId())
	} else if params.Status == consts.StatusOffline {
		be = logic.SetBannerOffline(params.IntId())
	}
	if be != nil {
		logger.Errorf("error: %v", be.String())
		resp.ReplyErr(be.ToResp())
		return
	}

	resp.ReplyOK()
}
