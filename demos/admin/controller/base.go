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
		base: &admin.Controller{
			Auth:     auth,
			LogSaver: &TraceLogSaver{},
		},
	}
	auth.CheckTokenFunc = func(ctx context.Context, c *app.RequestContext, claims *admin.Claims) (bool, *bizerr.Error) {
		return logic.TokenCheck(claims.Subject, claims.TokenId)
	}
	return ctrl
}

type TraceLogSaver struct {
}

func (t TraceLogSaver) AddLog(ctx *app.RequestContext, result resp.BaseResp) {
	if ctx.IsPost() {
		module, action, _ := admin.ParseRouterPath(ctx.FullPath())
		if len(module) == 0 || len(action) == 0 {
			return
		}
		params := map[string]interface{}{}
		if err := ctx.Bind(&params); err != nil {
			logger.Errorf("params: %v", err)
			return
		}
		if module == "platform" && action == "login" {
			if result.NotOK() {
				return
			}
			delete(params, "password")
			result.Data = nil
		}
		if module == "staff" && action == "reset_password" {
			result.Data = nil
		}

		uid := admin.GetSubject(ctx)
		roles := admin.GetAudience(ctx)
		logic.AddLog(module, action, params, result, roles, uid)
	}
}

type PlatformHandler struct {
	ctrl *admin.Controller
}

func (md PlatformHandler) Login(ctx context.Context, c *app.RequestContext) {
	resp := md.ctrl.Resp(c)
	params := defines.LoginReq{}
	if err := c.Bind(&params); err != nil {
		resp.ReplyErrorParam()
		return
	}
	if len(params.Username) == 0 {
		resp.ReplyErrorParam2("登录账号未输入")
		return
	}
	if len(params.Username) == 0 {
		resp.ReplyErrorParam2("登录密码未输入")
		return
	}

	uid, roles, be := logic.CheckStaff(params.Username, params.Password)
	if be != nil {
		logger.Errorf("error: %v", be.String())
		resp.ReplyErr(be.ToResp())
		return
	}

	token, claims, err := md.ctrl.CreateToken(c, uid, roles)
	if err != nil {
		logger.Errorf("error: %v", err)
		resp.ReplyErrorInternal()
		return
	}

	be = logic.TokenUpdate(claims.Subject, claims.TokenId, claims.ExpiredAt, claims.IssuedAt)
	if be != nil {
		logger.Errorf("error: %v", be.String())
		resp.ReplyErr(be.ToResp())
		return
	}

	resp.ReplyData(defines.TokenData{
		Token: token,
		Et:    claims.ExpiredAt,
	})
}

func (md PlatformHandler) Logout(ctx context.Context, c *app.RequestContext) {
	resp := md.ctrl.Resp(c)
	uid := md.ctrl.Uid(c)
	be := logic.TokenUpdate(uid, "", 0, 0)
	if be != nil {
		logger.Errorf("error: %v", be.String())
		resp.ReplyErr(be.ToResp())
		return
	}

	resp.ReplyOK()
}

func (md PlatformHandler) Profile(ctx context.Context, c *app.RequestContext) {
	resp := md.ctrl.Resp(c)
	uid := md.ctrl.Uid(c)
	userInfo, be := logic.GetStaff(uid)
	if be != nil {
		logger.Errorf("error: %v", be.String())
		resp.ReplyErr(be.ToResp())
		return
	}

	resp.ReplyData(userInfo)
}
