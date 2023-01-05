package svcmgr

import (
	"github.com/xhigher/hzgo/req"
	"github.com/xhigher/hzgo/resp"
)

const UserSvcName = "svc-user"

func UserClient() *SvcClient {
	return svcClients[UserSvcName]
}

func UserRegister(data interface{}) resp.BaseResp {
	return BaseAction{
		SvcName: UserSvcName,
		Name:    "register",
		Ver:     1,
		Data:    data,
	}.Do()
}

func UserProfile(userid string) resp.BaseResp {
	return BaseAction{
		SvcName: UserSvcName,
		Name:    "profile",
		Ver:     1,
		Data:    req.UseridReq{
			Userid: userid,
		},
	}.Do()
}

