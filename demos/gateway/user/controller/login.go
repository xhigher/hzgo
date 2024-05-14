package controller

import (
	"context"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/xhigher/hzgo/defines"
	"github.com/xhigher/hzgo/demo/api"
)

func (ctrl Controller) Login(ctx context.Context, c *app.RequestContext) {
	resp := ctrl.Resp(c)
	params := defines.LoginReq{}
	if err := c.Bind(&params); err != nil {
		resp.ReplyErrorParam()
		return
	}
	baseParams := ctrl.BaseParams(c)

	result := api.User().LoginCheck(defines.LoginReq{
		Username: params.Username,
		Password: params.Password,
	})
	if result.NotOK() {
		resp.ReplyErr(result)
		return
	}
	data := &defines.UseridData{}
	result.GetData(data)

	token, claims, err := ctrl.Auth.CreateToken(data.Userid, baseParams.Ap)
	if err != nil {
		resp.ReplyErrorInternal()
		return
	}

	result = api.User().TokenUpdate(defines.TokenUpdateReq{
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
