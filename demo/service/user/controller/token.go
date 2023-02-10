package controller

import (
	"context"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/xhigher/hzgo/defines"
	"github.com/xhigher/hzgo/demo/service/user/logic"
	"github.com/xhigher/hzgo/logger"
	"github.com/xhigher/hzgo/resp"
)

func (ctrl Controller) TokenCheck(ctx context.Context, c *app.RequestContext) {
	params := defines.TokenCheckReq{}
	if err := c.Bind(&params); err != nil {
		resp.ReplyErrorParam(c)
		return
	}
	ok, be := logic.CheckToken(params.Userid, params.TokenId)
	if be != nil {
		logger.Errorf("error: %v", be.String())
		resp.ReplyErr(c, be.ToResp())
		return
	}

	if !ok {
		resp.ReplyNOK(c)
		return
	}
	resp.ReplyOK(c)
}

func (ctrl Controller) TokenUpdate(ctx context.Context, c *app.RequestContext) {
	params := defines.TokenUpdateReq{}
	if err := c.Bind(&params); err != nil {
		resp.ReplyErrorParam(c)
		return
	}

	be := logic.UpdateToken(params.Audience, params.TokenId, params.ExpiredAt, params.IssuedAt)
	if be != nil {
		logger.Errorf("error: %v", be.String())
		resp.ReplyErr(c, be.ToResp())
		return
	}

	resp.ReplyOK(c)
}

