package controller

import (
	"context"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/xhigher/hzgo/demo/service/user/logic"
	"github.com/xhigher/hzgo/logger"
	"github.com/xhigher/hzgo/resp"
)

func (ctrl Controller) Profile(ctx context.Context, c *app.RequestContext) {
	userid, ok := ctrl.Userid(c)
	if !ok {
		return
	}

	userInfo, be := logic.GetUser(userid)
	if be != nil {
		logger.Errorf("error: %v", be.String())
		resp.ReplyErr(c, be.ToResp())
		return
	}

	resp.ReplyData(c, userInfo)
}
