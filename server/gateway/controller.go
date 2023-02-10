package gateway

import (
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/xhigher/hzgo/defines"
	"github.com/xhigher/hzgo/resp"
	"github.com/xhigher/hzgo/server/gateway/middlewares"
)

type Controller struct {
	Name string
	Auth *middlewares.JWTAuth
}

func (ctrl Controller) Userid(c *app.RequestContext) string {
	return middlewares.GetAudience(c)
}

func (ctrl Controller) Token(c *app.RequestContext) string {
	return middlewares.GetToken(c)
}

func (ctrl Controller) BaseParams(c *app.RequestContext) *defines.BaseParams {
	return middlewares.GetBaseParams(c)
}

func (ctrl Controller) ParamUserid(c *app.RequestContext)  (userid string, ok bool) {
	params := defines.UseridReq{}
	if err := c.Bind(&params); err != nil {
		resp.ReplyErrorParam(c)
		return
	}
	userid = params.Userid
	ok = true
	return
}