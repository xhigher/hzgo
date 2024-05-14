package controller

import (
	"context"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/xhigher/hzgo/demo/api"
	"github.com/xhigher/hzgo/logger"
	"github.com/xhigher/hzgo/utils"
)

func (ctrl Controller) Profile(ctx context.Context, c *app.RequestContext) {
	resp := ctrl.Resp(c)
	userid := ctrl.Userid(c)
	baseParams := ctrl.BaseParams(c)

	logger.Infof("baseParams: %v", utils.JSONString(baseParams))

	result := api.User().Profile(userid)
	if result.NotOK() {
		resp.ReplyErr(result)
		return
	}

	resp.ReplyData(result.Data)
}
