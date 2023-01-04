package controller

import (
	"context"
	"github.com/cloudwego/hertz/pkg/app"
	model_user "github.com/xhigher/hzgo/demo/model/user"
	"github.com/xhigher/hzgo/resp"
)

type RegisterReq struct {
	Username string `form:"username" json:"username" query:"username" vd:"(len($) > 0 && len($) < 128); msg:'Illegal format'"`
	Password string `form:"password" json:"password" query:"password" vd:"(len($) > 0 && len($) < 128); msg:'Illegal format'"`
}

func (ctrl Controller) Register(ctx context.Context, c *app.RequestContext) {
	req := RegisterReq{}
	if err := c.BindAndValidate(&req); err != nil {
		resp.ReplyErrorParam(c)
		return
	}
	userInfo, err := model_user.GetUser(req.Username)
	if err != nil {
		resp.ReplyErrorInternal(c)
		return
	}
	if userInfo != nil {
		resp.ReplyErr(c, resp.ErrorUserExisted)
		return
	}

	modelLogic := model_user.CreateUserLogic{
		Username: req.Username,
		Password: req.Password,
	}
	userInfo, existed, err := modelLogic.Do()
	if err != nil {
		resp.ReplyNOK(c)
		return
	}
	if existed {
		resp.ReplyErr(c, resp.ErrorUserExisted)
		return
	}

	resp.ReplyOK(c)
}

