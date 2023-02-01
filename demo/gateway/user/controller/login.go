package controller

import (
	"context"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/xhigher/hzgo/demo/gateway/api"
	"github.com/xhigher/hzgo/req"
	"github.com/xhigher/hzgo/resp"
)

type LoginReq struct {
	Username string `form:"username" json:"username" query:"username" vd:"(len($) > 0 && len($) < 128); msg:'Illegal format'"`
	Password string `form:"password" json:"password" query:"password" vd:"(len($) > 0 && len($) < 128); msg:'Illegal format'"`
}

type TokenData struct {
	Token string `json:"token"`
	Et    int64  `json:"et"`
}

type LoginData struct {
	Userid string `json:"userid"`
}

func (ctrl Controller) Login(ctx context.Context, c *app.RequestContext) {
	var params LoginReq
	if err := c.BindAndValidate(&params); err != nil {
		resp.ReplyErrorParam(c)
		return
	}

	result := api.User().LoginCheck(req.LoginReq{
		Username: params.Username,
		Password: params.Password,
	})
	if result.NotOK() {
		resp.ReplyErr(c, result)
		return
	}
	data := &LoginData{}
	result.GetData(data)

	token, claims, err := ctrl.Auth.CreateToken(data.Userid)
	if err != nil {
		resp.ReplyErrorInternal(c)
		return
	}

	result = api.User().TokenUpdate(req.TokenUpdateReq{
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
