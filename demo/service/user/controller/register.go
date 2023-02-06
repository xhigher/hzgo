package controller

import (
	"context"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/xhigher/hzgo/defines"
	model "github.com/xhigher/hzgo/demo/service/user/model/user"
	"github.com/xhigher/hzgo/resp"
)



func (ctrl Controller) Register(ctx context.Context, c *app.RequestContext) {
	params := defines.RegisterReq{}
	if err := c.Bind(&params); err != nil {
		resp.ReplyErrorParam(c)
		return
	}
	userInfo, err := model.GetUser(params.Username)
	if err != nil {
		resp.ReplyErrorInternal(c)
		return
	}
	if userInfo != nil {
		resp.ReplyErr(c, resp.ErrorUserExisted)
		return
	}

	modelLogic := model.CreateUserLogic{
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

	data := &defines.UseridData{
		Userid: userInfo.Userid,
	}

	resp.ReplyData(c, data)
}
