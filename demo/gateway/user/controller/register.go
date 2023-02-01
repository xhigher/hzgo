package controller

import (
	"context"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/xhigher/hzgo/demo/gateway/api"
	"github.com/xhigher/hzgo/req"
	"github.com/xhigher/hzgo/resp"
)

func (ctrl Controller) Register(ctx context.Context, c *app.RequestContext) {
	params := req.RegisterReq{}
	if err := c.Bind(&params); err != nil {
		resp.ReplyErrorParam(c)
		return
	}

	result := api.User().Register(params)
	if result.NotOK() {
		resp.ReplyErr(c, result)
		return
	}

	resp.ReplyOK(c)
}
