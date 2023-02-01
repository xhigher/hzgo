package controller

import (
	"context"
	"github.com/cloudwego/hertz/pkg/app"
	usermodel "github.com/xhigher/hzgo/demo/service/user/model/user"
	"github.com/xhigher/hzgo/req"
	"github.com/xhigher/hzgo/resp"
)

func (ctrl Controller) TokenCheck(ctx context.Context, c *app.RequestContext) {
	params := req.TokenCheckReq{}
	if err := c.Bind(&params); err != nil {
		resp.ReplyErrorParam(c)
		return
	}
	ok, err := usermodel.CheckToken(params.Userid, params.TokenId)
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
	params := req.TokenUpdateReq{}
	if err := c.Bind(&params); err != nil {
		resp.ReplyErrorParam(c)
		return
	}
	err := usermodel.SaveToken(params.Audience, params.TokenId, params.ExpiredAt, params.IssuedAt)
	if err != nil {
		resp.ReplyErrorInternal(c)
		return
	}

	resp.ReplyOK(c)
}

