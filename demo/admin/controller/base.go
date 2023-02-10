package controller

import (
	"context"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/xhigher/hzgo/bizerr"
	"github.com/xhigher/hzgo/defines"
	"github.com/xhigher/hzgo/demo/admin/controller/biz"
	logic "github.com/xhigher/hzgo/demo/admin/logic/platform"
	"github.com/xhigher/hzgo/logger"
	"github.com/xhigher/hzgo/resp"
	"github.com/xhigher/hzgo/server/admin"
)

type Controller struct {
	base *admin.Controller
}

func (ctrl Controller) Modules() []admin.Module {
	bizModules := biz.Modules(ctrl.base)
	baseModules := []admin.Module{
		StaffModule{
			ctrl: ctrl.base,
		},
	}
	return append(baseModules, bizModules...)
}

func (ctrl Controller) BaseController() *admin.Controller {
	return ctrl.base
}

func (ctrl Controller) PlatformHandler() admin.PlatformModuleHandler {
	return PlatformHandler{
		ctrl: ctrl.base,
	}
}

func New(auth *admin.Auth) Controller {
	ctrl := Controller{
		base:&admin.Controller{
			Auth: auth,
		},
	}
	auth.CheckTokenFunc = func(ctx context.Context, c *app.RequestContext, claims *admin.Claims) (bool, *bizerr.Error) {
		return logic.TokenCheck(claims.Subject, claims.TokenId)
	}
	return ctrl
}

type PlatformHandler struct {
	ctrl *admin.Controller
}

func (md PlatformHandler) Login(ctx context.Context, c *app.RequestContext){
	params := defines.LoginReq{}
	if err := c.Bind(&params); err != nil {
		resp.ReplyErrorParam(c)
		return
	}
	if len(params.Username) == 0 {
		resp.ReplyErrorParam2(c, "登录账号未输入")
		return
	}
	if len(params.Username) == 0 {
		resp.ReplyErrorParam2(c, "登录密码未输入")
		return
	}

	uid, roles, be := logic.CheckStaff(params.Username,params.Password)
	if be != nil {
		logger.Errorf("error: %v", be.String())
		resp.ReplyErr(c, be.ToResp())
		return
	}

	token, claims, err := md.ctrl.CreateToken(uid, roles)
	if err != nil {
		logger.Errorf("error: %v", err)
		resp.ReplyErrorInternal(c)
		return
	}

	be = logic.TokenUpdate(claims.Subject, claims.TokenId, claims.ExpiredAt, claims.IssuedAt)
	if be != nil {
		logger.Errorf("error: %v", be.String())
		resp.ReplyErr(c, be.ToResp())
		return
	}

	resp.ReplyData(c, defines.TokenData{
		Token: token,
		Et:    claims.ExpiredAt,
	})
}

func (md PlatformHandler) Logout(ctx context.Context, c *app.RequestContext) {
	uid := md.ctrl.Uid(c)

	be := logic.TokenUpdate(uid, "", 0, 0)
	if be != nil {
		logger.Errorf("error: %v", be.String())
		resp.ReplyErr(c, be.ToResp())
		return
	}

	resp.ReplyOK(c)
}

func (md PlatformHandler) Profile(ctx context.Context, c *app.RequestContext) {
	uid := md.ctrl.Uid(c)

	userInfo, be := logic.GetStaff(uid)
	if be != nil {
		logger.Errorf("error: %v", be.String())
		resp.ReplyErr(c, be.ToResp())
		return
	}

	resp.ReplyData(c, userInfo)
}


