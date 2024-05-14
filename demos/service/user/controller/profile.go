package controller

import (
	"context"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/xhigher/hzgo/demo/service/user/logic"
	"github.com/xhigher/hzgo/logger"
)

func (ctrl Controller) Profile(ctx context.Context, c *app.RequestContext) {
	resp := ctrl.Resp(c)
	userid, ok := ctrl.Userid(c)
	if !ok {
		resp.ReplyErrorParam2("userid")
		return
	}

	userInfo, be := logic.GetUser(userid)
	if be != nil {
		logger.Errorf("error: %v", be.String())
		resp.ReplyErr(be.ToResp())
		return
	}

	resp.ReplyData(userInfo)
}
