package controller

import (
	"context"
	"github.com/cloudwego/hertz/pkg/app"
	usermodel "github.com/xhigher/hzgo/demo/service/user/model/user"
	"github.com/xhigher/hzgo/resp"
)

func (ctrl Controller) Renewal(ctx context.Context, c *app.RequestContext) {
	token, claims, err := ctrl.Auth.RenewalToken(c)
	if err != nil {
		resp.ReplyErrorInternal(c)
		return
	}

	claims.IssuedAt = 0
	err = usermodel.SaveToken(claims.Audience, claims.TokenId, claims.ExpiredAt, claims.IssuedAt)
	if err != nil {
		resp.ReplyErrorInternal(c)
		return
	}

	resp.ReplyData(c, TokenData{
		Token: token,
		Et:    claims.ExpiredAt,
	})
}
