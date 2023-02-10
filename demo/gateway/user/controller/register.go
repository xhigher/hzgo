package controller

import (
	"context"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/xhigher/hzgo/defines"
	"github.com/xhigher/hzgo/demo/api"
	"github.com/xhigher/hzgo/resp"
)

func (ctrl Controller) Register(ctx context.Context, c *app.RequestContext) {
	params := defines.RegisterReq{}
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
