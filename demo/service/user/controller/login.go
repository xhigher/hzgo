package controller

import (
	"context"
	"github.com/cloudwego/hertz/pkg/app"
	userlogic "github.com/xhigher/hzgo/demo/service/user/logic/user"
	usermodel "github.com/xhigher/hzgo/demo/service/user/model/user"
	"github.com/xhigher/hzgo/req"
	"github.com/xhigher/hzgo/resp"
)

func (ctrl Controller) LoginCheck(ctx context.Context, c *app.RequestContext) {
	params := req.LoginReq{}
	if err := c.Bind(&params); err != nil {
		resp.ReplyErrorParam(c)
		return
	}
	userid, err := userlogic.CheckUser(params.Username, params.Password)
	if err != nil {
		resp.ReplyErrorInternal(c)
		return
	}

	userInfo := &usermodel.UserInfo{
		Userid: userid,
	}

	resp.ReplyData(c, userInfo)
}


