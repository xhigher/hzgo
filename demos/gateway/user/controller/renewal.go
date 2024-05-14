package controller

import (
	"context"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/xhigher/hzgo/defines"
	"github.com/xhigher/hzgo/demo/api"
)

func (ctrl Controller) Renewal(ctx context.Context, c *app.RequestContext) {
	resp := ctrl.Resp(c)
	token, claims, err := ctrl.Auth.RenewalToken(c)
	if err != nil {
		resp.ReplyErrorInternal()
		return
	}

	claims.IssuedAt = 0
	result := api.User().TokenUpdate(defines.TokenUpdateReq{
		Audience:  claims.Audience,
		TokenId:   claims.TokenId,
		ExpiredAt: claims.ExpiredAt,
		IssuedAt:  claims.IssuedAt,
	})
	if result.NotOK() {
		resp.ReplyErr(result)
		return
	}

	resp.ReplyData(defines.TokenData{
		Token: token,
		Et:    claims.ExpiredAt,
	})
}
