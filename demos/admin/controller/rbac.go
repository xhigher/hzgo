package controller

import (
	"context"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/xhigher/hzgo/consts"
	"github.com/xhigher/hzgo/defines"
	logic "github.com/xhigher/hzgo/demos/admin/logic/platform"
	"github.com/xhigher/hzgo/logger"
	"github.com/xhigher/hzgo/server/admin"
)

type RBACModule struct {
	ctrl *admin.Controller
}

func (md RBACModule) Name() string {
	return "rbac"
}

func (md RBACModule) Routers() []admin.Router {
	return []admin.Router{
		{
			Method:  consts.MethodGet,
			Name:    "role_list",
			Handler: md.RoleList,
		},
		{
			Method:  consts.MethodPost,
			Name:    "role_save",
			Handler: md.RoleSave,
		},
		{
			Method:  consts.MethodPost,
			Name:    "role_delete",
			Handler: md.RoleDelete,
		},
		{
			Method:  consts.MethodPost,
			Name:    "role_change_status",
			Handler: md.RoleChangeStatus,
		},

		{
			Method:  consts.MethodGet,
			Name:    "menu_list",
			Handler: md.MenuList,
		},
		{
			Method:  consts.MethodPost,
			Name:    "menu_save",
			Handler: md.MenuSave,
		},
		{
			Method:  consts.MethodPost,
			Name:    "menu_delete",
			Handler: md.MenuDelete,
		},
		{
			Method:  consts.MethodPost,
			Name:    "menu_change_status",
			Handler: md.MenuChangeStatus,
		},
		{
			Method:  consts.MethodPost,
			Name:    "role_permission_add",
			Handler: md.RolePermissionAdd,
		},
		{
			Method:  consts.MethodPost,
			Name:    "role_permission_delete",
			Handler: md.RolePermissionDelete,
		},
		{
			Method:  consts.MethodGet,
			Name:    "role_permission_list",
			Handler: md.RolePermissionList,
		},
		{
			Method:  consts.MethodPost,
			Name:    "role_menu_add",
			Handler: md.RoleMenuAdd,
		},
		{
			Method:  consts.MethodPost,
			Name:    "role_menu_delete",
			Handler: md.RoleMenuDelete,
		},
		{
			Method:  consts.MethodGet,
			Name:    "role_menu_list",
			Handler: md.RoleMenuList,
		},
	}
}

func (md RBACModule) RoleList(ctx context.Context, c *app.RequestContext) {
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

	data, be := logic.GetRoleList(params.Status)
	if be != nil {
		logger.Errorf("error: %v", be.String())
		resp.ReplyErr(be.ToResp())
		return
	}

	resp.ReplyData(data)
}

func (md RBACModule) RoleSave(ctx context.Context, c *app.RequestContext) {
	resp := md.ctrl.Resp(c)

	resp.ReplyOK()
}

func (md RBACModule) RoleDelete(ctx context.Context, c *app.RequestContext) {
	resp := md.ctrl.Resp(c)

	resp.ReplyOK()
}

func (md RBACModule) RoleChangeStatus(ctx context.Context, c *app.RequestContext) {
	resp := md.ctrl.Resp(c)

	resp.ReplyOK()
}

func (md RBACModule) MenuList(ctx context.Context, c *app.RequestContext) {
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

	data, be := logic.GetMenuList(params.Status)
	if be != nil {
		logger.Errorf("error: %v", be.String())
		resp.ReplyErr(be.ToResp())
		return
	}

	resp.ReplyData(data)
}

func (md RBACModule) MenuSave(ctx context.Context, c *app.RequestContext) {
	resp := md.ctrl.Resp(c)

	resp.ReplyOK()
}

func (md RBACModule) MenuDelete(ctx context.Context, c *app.RequestContext) {
	resp := md.ctrl.Resp(c)

	resp.ReplyOK()
}

func (md RBACModule) MenuChangeStatus(ctx context.Context, c *app.RequestContext) {
	resp := md.ctrl.Resp(c)

	resp.ReplyOK()
}

func (md RBACModule) RolePermissionAdd(ctx context.Context, c *app.RequestContext) {
	resp := md.ctrl.Resp(c)

	resp.ReplyOK()
}

func (md RBACModule) RolePermissionDelete(ctx context.Context, c *app.RequestContext) {
	resp := md.ctrl.Resp(c)

	resp.ReplyOK()
}

type RolePageReq struct {
	Rid    string `form:"uid" json:"uid" query:"uid"`
	Offset int32  `form:"offset" json:"offset" query:"offset"`
	Limit  int32  `form:"limit" json:"limit" query:"limit"`
}

func (md RBACModule) RolePermissionList(ctx context.Context, c *app.RequestContext) {
	resp := md.ctrl.Resp(c)
	params := RolePageReq{}
	if err := c.Bind(&params); err != nil {
		resp.ReplyErrorParam()
		return
	}
	if !defines.CheckPageLimit(params.Limit) {
		resp.ReplyErrorParam2("limit")
		return
	}

	total, data, be := logic.GetRolePermissionList(params.Rid, params.Limit, params.Offset)
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

func (md RBACModule) RoleMenuAdd(ctx context.Context, c *app.RequestContext) {
	resp := md.ctrl.Resp(c)

	resp.ReplyOK()
}

func (md RBACModule) RoleMenuDelete(ctx context.Context, c *app.RequestContext) {
	resp := md.ctrl.Resp(c)

	resp.ReplyOK()
}

func (md RBACModule) RoleMenuList(ctx context.Context, c *app.RequestContext) {
	resp := md.ctrl.Resp(c)
	params := RolePageReq{}
	if err := c.Bind(&params); err != nil {
		resp.ReplyErrorParam()
		return
	}
	if !defines.CheckPageLimit(params.Limit) {
		resp.ReplyErrorParam2("limit")
		return
	}

	total, data, be := logic.GetRoleMenuList(params.Rid, params.Limit, params.Offset)
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
