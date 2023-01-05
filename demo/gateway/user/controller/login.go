package controller

import (
	"context"
	"github.com/cloudwego/hertz/pkg/app"
	userlogic "github.com/xhigher/hzgo/demo/service/user/logic/user"
	usermodel "github.com/xhigher/hzgo/demo/service/user/model/user"
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

func (ctrl Controller) Login(ctx context.Context, c *app.RequestContext) {
	var req LoginReq
	if err := c.BindAndValidate(&req); err != nil {
		resp.ReplyErrorParam(c)
		return
	}
	userid, err := userlogic.CheckUser(req.Username, req.Password)
	if err != nil {
		resp.ReplyErrorInternal(c)
		return
	}

	token, claims, err := ctrl.Auth.CreateToken(userid)
	if err != nil {
		resp.ReplyErrorInternal(c)
		return
	}

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
