package controller

import (
	"context"
	"github.com/cloudwego/hertz/pkg/app"
	usermodel "github.com/xhigher/hzgo/demo/service/user/model/user"
	"github.com/xhigher/hzgo/resp"
)

func (ctrl Controller) Logout(ctx context.Context, c *app.RequestContext) {
	userid := ctrl.Userid(c)

	err := usermodel.SaveToken(userid, "", 0, 0)
	if err != nil {
		resp.ReplyErrorInternal(c)
		return
	}

	resp.ReplyOK(c)
}
