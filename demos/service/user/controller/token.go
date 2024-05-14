package controller

import (
	"context"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/xhigher/hzgo/defines"
	"github.com/xhigher/hzgo/demo/service/user/logic"
	"github.com/xhigher/hzgo/logger"
)

func (ctrl Controller) TokenCheck(ctx context.Context, c *app.RequestContext) {
	resp := ctrl.Resp(c)
	params := defines.TokenCheckReq{}
	if err := c.Bind(&params); err != nil {
		resp.ReplyErrorParam()
		return
	}
	ok, be := logic.CheckToken(params.Userid, params.TokenId)
	if be != nil {
		logger.Errorf("error: %v", be.String())
		resp.ReplyErr(be.ToResp())
		return
	}

	if !ok {
		resp.ReplyNOK()
		return
	}
	resp.ReplyOK()
}

func (ctrl Controller) TokenUpdate(ctx context.Context, c *app.RequestContext) {
	resp := ctrl.Resp(c)
	params := defines.TokenUpdateReq{}
	if err := c.Bind(&params); err != nil {
		resp.ReplyErrorParam()
		return
	}

	be := logic.UpdateToken(params.Audience, params.TokenId, params.ExpiredAt, params.IssuedAt)
	if be != nil {
		logger.Errorf("error: %v", be.String())
		resp.ReplyErr(be.ToResp())
		return
	}

	resp.ReplyOK()
}
