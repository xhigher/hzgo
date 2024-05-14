package controller

import (
	"context"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/xhigher/hzgo/defines"
	"github.com/xhigher/hzgo/demo/api"
)

func (ctrl Controller) Register(ctx context.Context, c *app.RequestContext) {
	resp := ctrl.Resp(c)
	params := defines.RegisterReq{}
	if err := c.Bind(&params); err != nil {
		resp.ReplyErrorParam()
		return
	}

	result := api.User().Register(params)
	if result.NotOK() {
		resp.ReplyErr(result)
		return
	}

	resp.ReplyOK()
}
