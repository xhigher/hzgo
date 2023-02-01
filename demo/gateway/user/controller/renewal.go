package controller

import (
	"context"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/xhigher/hzgo/demo/gateway/api"
	"github.com/xhigher/hzgo/req"
	"github.com/xhigher/hzgo/resp"
)

func (ctrl Controller) Renewal(ctx context.Context, c *app.RequestContext) {
	token, claims, err := ctrl.Auth.RenewalToken(c)
	if err != nil {
		resp.ReplyErrorInternal(c)
		return
	}

	claims.IssuedAt = 0
	result := api.User().TokenUpdate(req.TokenUpdateReq{
		Audience: claims.Audience,
		TokenId: claims.TokenId,
		ExpiredAt: claims.ExpiredAt,
		IssuedAt: claims.IssuedAt,
	})
	if result.NotOK() {
		resp.ReplyErr(c, result)
		return
	}

	resp.ReplyData(c, TokenData{
		Token: token,
		Et:    claims.ExpiredAt,
	})
}
