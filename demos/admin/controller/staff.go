package controller

import (
	"context"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/xhigher/hzgo/consts"
	"github.com/xhigher/hzgo/defines"
	logic "github.com/xhigher/hzgo/demo/admin/logic/platform"
	"github.com/xhigher/hzgo/demo/admin/rbac"
	"github.com/xhigher/hzgo/logger"
	"github.com/xhigher/hzgo/server/admin"
	"github.com/xhigher/hzgo/types"
)

type StaffModule struct {
	ctrl *admin.Controller
}

func (md StaffModule) Name() string {
	return rbac.Staff
}

func (md StaffModule) Roles() []string {
	return []string{admin.RoleMaintainer}
}

func (md StaffModule) Routers() []admin.Router {
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
			Path:    "reset_password",
			Handler: md.ResetPassword,
		},
		{
			Method:  consts.MethodPost,
			Path:    "update_roles",
			Handler: md.UpdateRoles,
		},
		{
			Method:  consts.MethodPost,
			Path:    "change_status",
			Handler: md.ChangeStatus,
		},
		{
			Method:  consts.MethodGet,
			Path:    "trace_logs",
			Handler: md.TraceLogs,
		},
	}
}

type UidReq struct {
	Uid string `json:"uid"`
}

func (md StaffModule) Info(ctx context.Context, c *app.RequestContext) {
	resp := md.ctrl.Resp(c)
	params := UidReq{}
	if err := c.Bind(&params); err != nil {
		resp.ReplyErrorParam()
		return
	}

	staffInfo, be := logic.GetStaff(params.Uid)
	if be != nil {
		logger.Errorf("error: %v", be.String())
		resp.ReplyErr(be.ToResp())
		return
	}

	resp.ReplyData(staffInfo)
}

func (md StaffModule) List(ctx context.Context, c *app.RequestContext) {
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

	total, staffList, be := logic.GetStaffList(params.Status, params.Offset, params.Limit)
	if be != nil {
		logger.Errorf("error: %v", be.String())
		resp.ReplyErr(be.ToResp())
		return
	}

	resp.ReplyData(defines.PageData{
		Total:  int32(total),
		Offset: params.Offset,
		Limit:  params.Limit,
		Data:   staffList,
	})
}

type CreateReq struct {
	Username string `json:"username"`
	Nickname string `json:"nickname"`
	Phone    string `json:"phone"`
	Email    string `json:"email"`
}

func (md StaffModule) Create(ctx context.Context, c *app.RequestContext) {
	resp := md.ctrl.Resp(c)
	params := CreateReq{}
	if err := c.Bind(&params); err != nil {
		resp.ReplyErrorParam()
		return
	}
	if len(params.Username) < 4 || len(params.Username) > 30 {
		resp.ReplyErrorParam2("username")
		return
	}
	if len(params.Phone) != 11 {
		resp.ReplyErrorParam2("phone")
		return
	}

	be := logic.CreateStaff(params.Username, params.Nickname, params.Phone, params.Email)
	if be != nil {
		logger.Errorf("error: %v", be.String())
		resp.ReplyErr(be.ToResp())
		return
	}

	resp.ReplyOK()
}

type ResetPasswordData struct {
	Uid      string `json:"uid"`
	Password string `json:"password"`
}

func (md StaffModule) ResetPassword(ctx context.Context, c *app.RequestContext) {
	resp := md.ctrl.Resp(c)
	params := UidReq{}
	if err := c.Bind(&params); err != nil {
		resp.ReplyErrorParam()
		return
	}

	password, be := logic.ResetStaffPassword(params.Uid)
	if be != nil {
		logger.Errorf("error: %v", be.String())
		resp.ReplyErr(be.ToResp())
		return
	}

	resp.ReplyData(ResetPasswordData{
		Uid:      params.Uid,
		Password: password,
	})
}

type UpdateRolesReq struct {
	Uid   string            `json:"uid"`
	Roles types.StringArray `json:"roles"`
}

func (md StaffModule) UpdateRoles(ctx context.Context, c *app.RequestContext) {
	resp := md.ctrl.Resp(c)
	params := UpdateRolesReq{}
	if err := c.Bind(&params); err != nil {
		resp.ReplyErrorParam()
		return
	}
	if !admin.CheckRoles(params.Roles) {
		resp.ReplyErrorParam()
		return
	}

	be := logic.UpdateStaffRoles(params.Uid, params.Roles)
	if be != nil {
		logger.Errorf("error: %v", be.String())
		resp.ReplyErr(be.ToResp())
		return
	}

	resp.ReplyOK()
}

type ChangeStatusReq struct {
	Uid    string `json:"uid"`
	Status int32  `json:"roles"`
}

func (md StaffModule) ChangeStatus(ctx context.Context, c *app.RequestContext) {
	resp := md.ctrl.Resp(c)
	params := ChangeStatusReq{}
	if err := c.Bind(&params); err != nil {
		resp.ReplyErrorParam()
		return
	}

	myUid := md.ctrl.Uid(c)
	if myUid == params.Uid {
		resp.ReplyNOK()
		return
	}

	be := logic.UpdateStaffStatus(params.Uid, params.Status)
	if be != nil {
		logger.Errorf("error: %v", be.String())
		resp.ReplyErr(be.ToResp())
		return
	}

	resp.ReplyOK()
}

type TraceLogPageReq struct {
	Uid    string `form:"uid" json:"uid" query:"uid"`
	Module string `form:"module" json:"module" query:"module"`
	Offset int32  `form:"offset" json:"offset" query:"offset"`
	Limit  int32  `form:"limit" json:"limit" query:"limit"`
}

func (md StaffModule) TraceLogs(ctx context.Context, c *app.RequestContext) {
	resp := md.ctrl.Resp(c)
	params := TraceLogPageReq{}
	if err := c.Bind(&params); err != nil {
		resp.ReplyErrorParam()
		return
	}
	if !defines.CheckPageLimit(params.Limit) {
		resp.ReplyErrorParam2("limit")
		return
	}

	total, logs, be := logic.GetTraceLogs(params.Uid, params.Module, params.Offset, params.Limit)
	if be != nil {
		logger.Errorf("error: %v", be.String())
		resp.ReplyErr(be.ToResp())
		return
	}

	resp.ReplyData(defines.PageData{
		Total:  int32(total),
		Offset: params.Offset,
		Limit:  params.Limit,
		Data:   logs,
	})
}
