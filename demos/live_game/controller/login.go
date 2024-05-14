package controller

import (
	"context"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/xhigher/hzgo/games/live_game/logic"
)

type WechatLoginReq struct {
	Code          string `form:"code" json:"code" query:"code"`
	Iv            string `form:"iv" json:"iv" query:"iv"`
	EncryptedData string `form:"encrypted_data" json:"encrypted_data" query:"encrypted_data"`
	Signature     string `form:"signature" json:"signature" query:"signature"`
	Latitude      string `form:"latitude" json:"latitude" query:"latitude"`
	Longitude     string `form:"longitude" json:"longitude" query:"longitude"`
	Channel       string `form:"channel" json:"channel" query:"channel"`
	ParentId      string `form:"parentId" json:"parentId" query:"parentId"`
}

//https://hifun.yunwan.tech/game/gameUser/weChatLogin?
//code=0c3cHa100IdPJP1x4Q300CotbM1cHa1A&
//iv=d%2BnFkPoGiYjSisd%2Btr%2Ba0w%3D%3D&
//encryptedData=2XLisRABaKyMfszoo%2BhiAEgtNP%2Bliu2ld3X7DEZH1z9Nzgc82ZstvJWjK6AHWNLLISFwysalzH0NZTNYfKVwHO3zFu5ZnxpF6miKLj68O8Ej64PEnaykdRjh5Xmaq02QT9clewvWYy5c%2BzikzamKf5anxmfnVGZEfMivFEVEnIbV4Wme0cZW2f6f1dkJpnkVsbhhvZwKoW8kLAsK%2FmEAluX7x7OMrOZqwLiTeLfXVnccNXvVXa7u6W27sIJToEFiN9p20sGMsd0LkfgztsYRApXUR8Snp7og2DZo9R5aSk96hsnxeToyQ8z%2FA4hkA8XcotVKrDOHpNuwRboHpF%2FJYmJgYhC9ySRMNU60Oxt4nRoKFz3r2tOCDg9l2fFwCJ8JaMKQJbSTIYzfvT0LZ1S2bBqI2cDEeleU7DJU8VEgGvEVQM45SKjWaOJH2TOSBWXiZAV4IeunkKuKlwPfyyMplw%3D%3D
//&signature=06b7ae23b8f62da1bda80783676f5db252d37726&
//latitude=999&
//longitude=999&
//parentId=
//&channel=
//&loginType=1

func (c Controller) WechatLogin(ctx context.Context, reqCtx *app.RequestContext) {
	resp := c.Resp(reqCtx)
	params := WechatLoginReq{}
	if err := reqCtx.Bind(&params); err != nil {
		resp.ReplyErrorParam()
		return
	}

	id := "fsbi6y"
	_, userInfo := logic.GetUser(id)

	resp.ReplyData(userInfo)
}
