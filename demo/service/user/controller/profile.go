package controller

import (
	"context"
	"github.com/cloudwego/hertz/pkg/app"
	model_user "github.com/xhigher/hzgo/demo/service/user/model/user"
	"github.com/xhigher/hzgo/resp"
)

func (ctrl Controller) Profile(ctx context.Context, c *app.RequestContext) {
	userid, ok := ctrl.Userid(c)
	if !ok {
		return
	}
	userInfo, err := model_user.GetUserById(userid)
	if err != nil {
		resp.ReplyErrorInternal(c)
		return
	}
	if userInfo == nil {
		resp.ReplyErr(c, resp.ErrorUserExisted)
		return
	}

	resp.ReplyData(c, userInfo)
}
