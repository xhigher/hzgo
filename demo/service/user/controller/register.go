package controller

import (
	"context"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/xhigher/hzgo/defines"
	"github.com/xhigher/hzgo/demo/service/user/logic"
	"github.com/xhigher/hzgo/logger"
	"github.com/xhigher/hzgo/resp"
)



func (ctrl Controller) Register(ctx context.Context, c *app.RequestContext) {
	params := defines.RegisterReq{}
	if err := c.Bind(&params); err != nil {
		resp.ReplyErrorParam(c)
		return
	}

	userid, be := logic.CreateUser(params.Username, params.Password)
	if be != nil {
		logger.Errorf("error: %v", be.String())
		resp.ReplyErr(c, be.ToResp())
		return
	}

	data := &defines.UseridData{
		Userid: userid,
	}

	resp.ReplyData(c, data)
}
