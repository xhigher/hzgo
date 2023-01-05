package controller

import (
	"context"
	"github.com/cloudwego/hertz/pkg/app"
	usermodel "github.com/xhigher/hzgo/demo/service/user/model/user"
	"github.com/xhigher/hzgo/req"
	"github.com/xhigher/hzgo/resp"
)



func (ctrl Controller) Register(ctx context.Context, c *app.RequestContext) {
	params := req.RegisterReq{}
	if err := c.Bind(&params); err != nil {
		resp.ReplyErrorParam(c)
		return
	}
	userInfo, err := usermodel.GetUser(params.Username)
	if err != nil {
		resp.ReplyErrorInternal(c)
		return
	}
	if userInfo != nil {
		resp.ReplyErr(c, resp.ErrorUserExisted)
		return
	}

	modelLogic := usermodel.CreateUserLogic{
		Username: params.Username,
		Password: params.Password,
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
