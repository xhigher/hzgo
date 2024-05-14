package controller

import (
	"context"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/xhigher/hzgo/demos/live_game/logic"
	"github.com/xhigher/hzgo/logger"
)

type TokenReq struct {
	Id    string `form:"id" json:"id" query:"id"`
	Token string `form:"token" json:"token" query:"token"`
}

func (c Controller) UserInfo(ctx context.Context, reqCtx *app.RequestContext) {
	resp := c.Resp(reqCtx)
	params := TokenReq{}
	if err := reqCtx.Bind(&params); err != nil {
		resp.ReplyErrorParam()
		return
	}

	if !c.checkToken(params.Id, params.Token) {
		resp.ReplyErrorAuthorization()
		return
	}

	logger.Infof("GetUser: %s", params.Id)
	_, userInfo := logic.GetUser(params.Id)

	resp.ReplyData(userInfo)
}

//https://hifun.yunwan.tech/game/gameUser/?route=getMessage&_id=fsbi6y&token=IKFQDBKFBGZLLNVRNMFP
//https://hifun.yunwan.tech/game/gameUser/?childName=store&route=getChild&_id=fsbi6y&token=IKFQDBKFBGZLLNVRNMFP
//https://hifun.yunwan.tech/game/gameUser/?lastMessageId=&channel=World&hint=1&route=updateLobby&_id=fsbi6y&token=IKFQDBKFBGZLLNVRNMFP
//https://hifun.yunwan.tech/game/gameUser/?type=0&v=0.0.7&route=ppt.getRandomServer&_id=fsbi6y&token=IKFQDBKFBGZLLNVRNMFP
