package controller

import (
	"context"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/xhigher/hzgo/defines"
	"github.com/xhigher/hzgo/demo/service/user/model"
	"github.com/xhigher/hzgo/resp"
)

func (ctrl Controller) TokenCheck(ctx context.Context, c *app.RequestContext) {
	params := defines.TokenCheckReq{}
	if err := c.Bind(&params); err != nil {
		resp.ReplyErrorParam(c)
		return
	}
	ok, err := model.CheckToken(params.Userid, params.TokenId)
	if err != nil {
		resp.ReplyErrorInternal(c)
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
	err := model.SaveToken(params.Audience, params.TokenId, params.ExpiredAt, params.IssuedAt)
	if err != nil {
		resp.ReplyErrorInternal(c)
		return
	}

	resp.ReplyOK(c)
}

