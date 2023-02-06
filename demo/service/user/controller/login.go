package controller

import (
	"context"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/xhigher/hzgo/defines"
	logic "github.com/xhigher/hzgo/demo/service/user/logic/user"
	"github.com/xhigher/hzgo/resp"
)

func (ctrl Controller) LoginCheck(ctx context.Context, c *app.RequestContext) {
	params := defines.LoginReq{}
	if err := c.Bind(&params); err != nil {
		resp.ReplyErrorParam(c)
		return
	}
	userid, err := logic.CheckUser(params.Username, params.Password)
	if err != nil {
		resp.ReplyErrorInternal(c)
		return
	}

	data := &defines.UseridData{
		Userid: userid,
	}

	resp.ReplyData(c, data)
}


