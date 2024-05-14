package controller

import (
	"context"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/xhigher/hzgo/defines"
	"github.com/xhigher/hzgo/demo/service/misc/logic"
	"github.com/xhigher/hzgo/logger"
)

func (ctrl Controller) ConfigList(ctx context.Context, c *app.RequestContext) {
	resp := ctrl.Resp(c)
	params := defines.ConfigReq{}
	if err := c.Bind(&params); err != nil {
		resp.ReplyErrorParam()
		return
	}

	sum, data, be := logic.GetConfigList()
	if be != nil {
		logger.Errorf("error: %v", be.String())
		resp.ReplyErr(be.ToResp())
		return
	}
	if params.Sum == sum {
		resp.ReplyOK()
		return
	}

	resp.ReplyData(&defines.ConfigListData{
		Sum:  sum,
		Data: data,
	})
}

func (ctrl Controller) ConfigInfo(ctx context.Context, c *app.RequestContext) {
	resp := ctrl.Resp(c)
	params := defines.ConfigReq{}
	if err := c.Bind(&params); err != nil {
		resp.ReplyErrorParam()
		return
	}

	data, be := logic.GetConfigInfo(params.Id)
	if be != nil {
		logger.Errorf("error: %v", be.String())
		resp.ReplyErr(be.ToResp())
		return
	}

	resp.ReplyData(data)
}

func (ctrl Controller) ConfigReload(ctx context.Context, c *app.RequestContext) {
	resp := ctrl.Resp(c)
	params := defines.ConfigReq{}
	if err := c.Bind(&params); err != nil {
		resp.ReplyErrorParam()
		return
	}

	be := logic.ReloadConfig(params.Id)
	if be != nil {
		logger.Errorf("error: %v", be.String())
		resp.ReplyErr(be.ToResp())
		return
	}

	resp.ReplyOK()
}
