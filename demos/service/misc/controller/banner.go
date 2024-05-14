package controller

import (
	"context"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/xhigher/hzgo/bizerr"
	"github.com/xhigher/hzgo/defines"
	"github.com/xhigher/hzgo/demo/service/misc/logic"
	"github.com/xhigher/hzgo/logger"
)

func (ctrl Controller) BannerList(ctx context.Context, c *app.RequestContext) {
	resp := ctrl.Resp(c)
	params := defines.BannerReq{}
	if err := c.Bind(&params); err != nil {
		resp.ReplyErrorParam()
		return
	}

	var data interface{}
	var be *bizerr.Error
	if len(params.Site) > 0 {
		data, be = logic.GetSiteBannerList(params.Site)
	} else {
		data, be = logic.GetAllBannerList()
	}
	if be != nil {
		logger.Errorf("error: %v", be.String())
		resp.ReplyErr(be.ToResp())
		return
	}

	resp.ReplyData(data)
}

func (ctrl Controller) BannerReload(ctx context.Context, c *app.RequestContext) {
	resp := ctrl.Resp(c)
	be := logic.ReloadBanner()
	if be != nil {
		logger.Errorf("error: %v", be.String())
		resp.ReplyErr(be.ToResp())
		return
	}

	resp.ReplyOK()
}
