package api

import (
	"github.com/xhigher/hzgo/defines"
	"github.com/xhigher/hzgo/resp"
	"github.com/xhigher/hzgo/svcmgr"
)

type SvcUser struct {
	name string
}

var svcUser = SvcUser{
	name: "svc-user",
}

func User() SvcUser {
	return svcUser
}

func (s SvcUser) Register(data interface{}) resp.BaseResp {
	return svcmgr.BaseAction{
		SvcName: s.name,
		Name:    "register",
		Ver:     1,
		Data:    data,
	}.Do()
}

func (s SvcUser) Profile(userid string) resp.BaseResp {
	return svcmgr.BaseAction{
		SvcName: s.name,
		Name:    "profile",
		Ver:     1,
		Data: defines.UseridReq{
			Userid: userid,
		},
	}.Do()
}

func (s SvcUser) LoginCheck(data interface{}) resp.BaseResp {
	return svcmgr.BaseAction{
		SvcName: s.name,
		Name:    "login_check",
		Ver:     1,
		Data:    data,
	}.Do()
}

func (s SvcUser) TokenUpdate(data interface{}) resp.BaseResp {
	return svcmgr.BaseAction{
		SvcName: s.name,
		Name:    "token_update",
		Ver:     1,
		Data:    data,
	}.Do()
}

func (s SvcUser) TokenCheck(data interface{}) resp.BaseResp {
	return svcmgr.BaseAction{
		SvcName: s.name,
		Name:    "token_check",
		Ver:     1,
		Data:    data,
	}.Do()
}

func (s SvcUser) Logout(userid string) resp.BaseResp {
	return svcmgr.BaseAction{
		SvcName: s.name,
		Name:    "logout",
		Ver:     1,
		Data: defines.UseridReq{
			Userid: userid,
		},
	}.Do()
}
