package service

import (
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/xhigher/hzgo/req"
	"github.com/xhigher/hzgo/resp"
)

type Controller struct {
	Name string
}

func (ctrl Controller) Userid(c *app.RequestContext) (userid string, ok bool) {
	params := vars.UseridReq{}
	if err := c.Bind(&params); err != nil {
		resp.ReplyErrorParam(c)
		return
	}
	userid = params.Userid
	ok = true
	return
}

