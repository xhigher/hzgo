package resp

import (
	"fmt"
	"github.com/cloudwego/hertz/pkg/app"
	"net/http"
)

var (
	OK  = NewCode(0)
	NOK = NewMsg(1, "请求失败")

	ErrorStop       = NewMsg(100, "系统维护中")
	ErrorInternal   = NewMsg(101, "服务器错误")
	ErrorParam      = NewMsg(102, "参数错误")
	ErrorPermission = NewMsg(103, "无权限访问")
	ErrorIllegal    = NewMsg(104, "非法请求")

	ErrorAuthorization = NewMsg(201, "您的账号未登录")
	ErrorUserNull      = NewMsg(202, "您的账号未注册")
	ErrorUserExisted      = NewMsg(203, "您的账号已注册")
	ErrorUserCancel    = NewMsg(204, "您的账号已注销")
	ErrorUserBlock     = NewMsg(205, "您的账号已封禁")
	ErrorUserLogout   = NewMsg(206, "您的账号被踢出登录")


)
type BaseResp struct {
	Code int32 `json:"code"`
	Msg  string `json:"msg"`
	Data interface{} `json:"data"`
}

func NewCode(code int32) BaseResp {
	return NewData(code, "", nil)
}

func NewMsg(code int32, msg string) BaseResp {
	return NewData(code, msg, nil)
}

func NewData(code int32, msg string, data interface{}) BaseResp {
	return BaseResp{
		Code: code,
		Msg:  msg,
		Data: data,
	}
}


func (e BaseResp) String() string {
	return fmt.Sprintf("error code: %d, msg: %s", e.Code, e.Msg)
}

func (e BaseResp) IsOK() bool {
	return OK.Code == e.Code
}

func (e BaseResp) NotOK() bool {
	return !e.IsOK()
}

func ReplyOK(ctx *app.RequestContext){
	ctx.JSON(http.StatusOK, OK)
}

func ReplyNOK(ctx *app.RequestContext){
	ctx.JSON(http.StatusOK, NOK)
}

func ReplyErrorInternal(ctx *app.RequestContext){
	ctx.JSON(http.StatusOK, ErrorInternal)
}

func ReplyErrorParam(ctx *app.RequestContext){
	ctx.JSON(http.StatusOK, ErrorParam)
}

func ReplyErrorPermission(ctx *app.RequestContext){
	ctx.JSON(http.StatusOK, ErrorPermission)
}

func ReplyErrorIllegal(ctx *app.RequestContext){
	ctx.JSON(http.StatusOK, ErrorIllegal)
}

func ReplyErrorAuthorization(ctx *app.RequestContext){
	ctx.JSON(http.StatusOK, ErrorAuthorization)
}

func ReplyErrMsg(ctx *app.RequestContext, msg string){
	ctx.JSON(http.StatusOK, BaseResp{
		Code: NOK.Code,
		Msg: msg,
	})
}

func ReplyErr(ctx *app.RequestContext, err BaseResp){
	ctx.JSON(http.StatusOK, BaseResp{
		Code: err.Code,
		Msg: err.Msg,
	})
}

func ReplyData(ctx *app.RequestContext, data interface{}){
	ctx.JSON(http.StatusOK, BaseResp{
		Code: OK.Code,
		Msg: OK.Msg,
		Data: data,
	})
}

func Reply(ctx *app.RequestContext, code int32, msg string, data interface{}){
	ctx.JSON(http.StatusOK, BaseResp{
		Code: code,
		Msg: msg,
		Data: data,
	})
}
