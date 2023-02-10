package controller

import (
	"context"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/xhigher/hzgo/consts"
	"github.com/xhigher/hzgo/defines"
	logic "github.com/xhigher/hzgo/demo/admin/logic/platform"
	"github.com/xhigher/hzgo/demo/admin/rbac"
	"github.com/xhigher/hzgo/logger"
	"github.com/xhigher/hzgo/resp"
	"github.com/xhigher/hzgo/server/admin"
	"github.com/xhigher/hzgo/types"
)

type StaffModule struct {
	ctrl *admin.Controller
}

func (md StaffModule) Name() string{
	return rbac.Staff
}

func (md StaffModule) Roles() []string{
	return []string{admin.RoleMaintainer}
}

func (md StaffModule) Routers() []admin.Router{
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
	}
}

type UidReq struct {
	Uid string `json:"uid"`
}

func (md StaffModule) Info(ctx context.Context, c *app.RequestContext) {
	params := UidReq{}
	if err := c.Bind(&params); err != nil {
		resp.ReplyErrorParam(c)
		return
	}

	staffInfo, be := logic.GetStaff(params.Uid)
	if be != nil {
		logger.Errorf("error: %v", be.String())
		resp.ReplyErr(c, be.ToResp())
		return
	}

	resp.ReplyData(c, staffInfo)
}

func (md StaffModule) List(ctx context.Context, c *app.RequestContext) {
	params := defines.StatusPageReq{}
	if err := c.Bind(&params); err != nil {
		resp.ReplyErrorParam(c)
		return
	}
	if !defines.CheckPageLimit(params.Limit) {
		resp.ReplyErrorParam2(c, "limit")
		return
	}

	total, staffList, be := logic.GetStaffList(params.Status, params.Offset, params.Limit)
	if be != nil {
		logger.Errorf("error: %v", be.String())
		resp.ReplyErr(c, be.ToResp())
		return
	}

	resp.ReplyData(c, defines.PageData{
		Total: int32(total),
		Offset: params.Offset,
		Limit: params.Limit,
		Data: staffList,
	})
}

type CreateReq struct {
	Username string `json:"username"`
	Nickname string `json:"nickname"`
	Phone string `json:"phone"`
	Email string `json:"email"`
}

func (md StaffModule) Create(ctx context.Context, c *app.RequestContext){
	params := CreateReq{}
	if err := c.Bind(&params); err != nil {
		resp.ReplyErrorParam(c)
		return
	}
	if len(params.Username) <4 || len(params.Username) > 30 {
		resp.ReplyErrorParam2(c, "username")
		return
	}
	if len(params.Phone) != 11 {
		resp.ReplyErrorParam2(c, "phone")
		return
	}

	be := logic.CreateStaff(params.Username, params.Nickname, params.Phone, params.Email)
	if be != nil {
		logger.Errorf("error: %v", be.String())
		resp.ReplyErr(c, be.ToResp())
		return
	}

	resp.ReplyOK(c)
}

type ResetPasswordData struct {
	Uid string `json:"uid"`
	Password string `json:"password"`
}

func (md StaffModule) ResetPassword(ctx context.Context, c *app.RequestContext){
	params := UidReq{}
	if err := c.Bind(&params); err != nil {
		resp.ReplyErrorParam(c)
		return
	}

	password, be := logic.ResetStaffPassword(params.Uid)
	if be != nil {
		logger.Errorf("error: %v", be.String())
		resp.ReplyErr(c, be.ToResp())
		return
	}

	resp.ReplyData(c, ResetPasswordData{
		Uid: params.Uid,
		Password: password,
	})
}

type UpdateRolesReq struct {
	Uid string                `json:"uid"`
	Roles types.StringArray `json:"roles"`
}

func (md StaffModule) UpdateRoles(ctx context.Context, c *app.RequestContext) {
	params := UpdateRolesReq{}
	if err := c.Bind(&params); err != nil {
		resp.ReplyErrorParam(c)
		return
	}
	if !admin.CheckRoles(params.Roles){
		resp.ReplyErrorParam(c)
		return
	}

	be := logic.UpdateStaffRoles(params.Uid, params.Roles)
	if be != nil {
		logger.Errorf("error: %v", be.String())
		resp.ReplyErr(c, be.ToResp())
		return
	}

	resp.ReplyOK(c)
}

type ChangeStatusReq struct {
	Uid string `json:"uid"`
	Status int32 `json:"roles"`
}

func (md StaffModule) ChangeStatus(ctx context.Context, c *app.RequestContext) {
	params := ChangeStatusReq{}
	if err := c.Bind(&params); err != nil {
		resp.ReplyErrorParam(c)
		return
	}
	myUid := md.ctrl.Uid(c)
	if myUid == params.Uid {
		resp.ReplyNOK(c)
		return
	}

	be := logic.UpdateStaffStatus(params.Uid, params.Status)
	if be != nil {
		logger.Errorf("error: %v", be.String())
		resp.ReplyErr(c, be.ToResp())
		return
	}

	resp.ReplyOK(c)
}
