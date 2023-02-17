package controller

import (
	"context"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/xhigher/hzgo/defines"
	"github.com/xhigher/hzgo/demo/service/misc/logic"
	"github.com/xhigher/hzgo/logger"
	"github.com/xhigher/hzgo/resp"
)

func (ctrl Controller) ConfigList(ctx context.Context, c *app.RequestContext) {
	params := defines.ConfigReq{}
	if err := c.Bind(&params); err != nil {
		resp.ReplyErrorParam(c)
		return
	}

	sum, data, be := logic.GetConfigList()
	if be != nil {
		logger.Errorf("error: %v", be.String())
		resp.ReplyErr(c, be.ToResp())
		return
	}
	if params.Sum == sum {
		resp.ReplyOK(c)
		return
	}

	resp.ReplyData(c, &defines.ConfigListData{
		Sum: sum,
		Data: data,
	})
}

func (ctrl Controller) ConfigInfo(ctx context.Context, c *app.RequestContext) {
	params := defines.ConfigReq{}
	if err := c.Bind(&params); err != nil {
		resp.ReplyErrorParam(c)
		return
	}

	data, be := logic.GetConfigInfo(params.Id)
	if be != nil {
		logger.Errorf("error: %v", be.String())
		resp.ReplyErr(c, be.ToResp())
		return
	}

	resp.ReplyData(c, data)
}

func (ctrl Controller) ConfigReload(ctx context.Context, c *app.RequestContext) {
	params := defines.ConfigReq{}
	if err := c.Bind(&params); err != nil {
		resp.ReplyErrorParam(c)
		return
	}

	be := logic.ReloadConfig(params.Id)
	if be != nil {
		logger.Errorf("error: %v", be.String())
		resp.ReplyErr(c, be.ToResp())
		return
	}

	resp.ReplyOK(c)
}

