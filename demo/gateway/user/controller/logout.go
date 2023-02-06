package controller

import (
	"context"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/xhigher/hzgo/defines"
	"github.com/xhigher/hzgo/demo/gateway/api"
	"github.com/xhigher/hzgo/resp"
)

func (ctrl Controller) Logout(ctx context.Context, c *app.RequestContext) {
	userid := ctrl.Userid(c)

	result := api.User().TokenUpdate(defines.TokenUpdateReq{
		Audience: userid,
		TokenId: "",
		ExpiredAt: 0,
		IssuedAt: 0,
	})
	if result.NotOK() {
		resp.ReplyErr(c, result)
		return
	}

	resp.ReplyOK(c)
}
