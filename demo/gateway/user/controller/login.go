package controller

import (
	"context"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/xhigher/hzgo/defines"
	"github.com/xhigher/hzgo/demo/api"
	"github.com/xhigher/hzgo/resp"
)

func (ctrl Controller) Login(ctx context.Context, c *app.RequestContext) {
	params := defines.LoginReq{}
	if err := c.Bind(&params); err != nil {
		resp.ReplyErrorParam(c)
		return
	}
	baseParams := ctrl.BaseParams(c)

	result := api.User().LoginCheck(defines.LoginReq{
		Username: params.Username,
		Password: params.Password,
	})
	if result.NotOK() {
		resp.ReplyErr(c, result)
		return
	}
	data := &defines.UseridData{}
	result.GetData(data)

	token, claims, err := ctrl.Auth.CreateToken(data.Userid, baseParams.Ap)
	if err != nil {
		resp.ReplyErrorInternal(c)
		return
	}

	result = api.User().TokenUpdate(defines.TokenUpdateReq{
		Audience: claims.Audience,
		TokenId: claims.TokenId,
		ExpiredAt: claims.ExpiredAt,
		IssuedAt: claims.IssuedAt,
	})
	if result.NotOK() {
		resp.ReplyErr(c, result)
		return
	}

	resp.ReplyData(c, defines.TokenData{
		Token: token,
		Et:    claims.ExpiredAt,
	})
}
