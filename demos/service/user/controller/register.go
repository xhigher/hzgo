package controller

import (
	"context"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/xhigher/hzgo/defines"
	"github.com/xhigher/hzgo/demo/service/user/logic"
	"github.com/xhigher/hzgo/logger"
)

func (ctrl Controller) Register(ctx context.Context, c *app.RequestContext) {
	resp := ctrl.Resp(c)
	params := defines.RegisterReq{}
	if err := c.Bind(&params); err != nil {
		resp.ReplyErrorParam()
		return
	}

	userid, be := logic.CreateUser(params.Username, params.Password)
	if be != nil {
		if be.UserExists() {
			resp.ReplyErrorUserExists()
			return
		}
		logger.Errorf("error: %v", be.String())
		resp.ReplyErr(be.ToResp())
		return
	}

	data := &defines.UseridData{
		Userid: userid,
	}

	resp.ReplyData(data)
}
