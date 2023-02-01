package controller

import (
	"context"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/xhigher/hzgo/demo/gateway/api"
	"github.com/xhigher/hzgo/resp"
)

func (ctrl Controller) Profile(ctx context.Context, c *app.RequestContext) {
	userid := ctrl.Userid(c)

	result := api.User().Profile(userid)
	if result.NotOK() {
		resp.ReplyErr(c, result)
		return
	}

	resp.ReplyData(c, result.Data)
}
