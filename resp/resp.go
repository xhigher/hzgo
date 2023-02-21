package resp

import (
	"encoding/json"
	"fmt"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/xhigher/hzgo/logger"
	"github.com/xhigher/hzgo/utils"
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
	ErrorUserExists   = NewMsg(203, "您的账号已注册")
	ErrorUserCanceled    = NewMsg(204, "您的账号已注销")
	ErrorUserBlocked     = NewMsg(205, "您的账号已封禁")
	ErrorUserLogout    = NewMsg(206, "您的账号被踢出登录")

)

type BaseResp struct {
	Code int32       `json:"code"`
	Msg  string      `json:"msg"`
	Data json.RawMessage `json:"data"`
}

func NewCode(code int32) BaseResp {
	return NewData(code, "", nil)
}

func NewMsg(code int32, msg string) BaseResp {
	return NewData(code, msg, nil)
}

func NewData(code int32, msg string, data interface{}) BaseResp {
	raw, _ := json.Marshal(data)
	return BaseResp{
		Code: code,
		Msg:  msg,
		Data: raw,
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

func (e BaseResp) GetData(v interface{}) error {
	return json.Unmarshal(e.Data, &v)
}

type TraceLogSaver interface {
	AddLog(c *app.RequestContext, resp BaseResp)
}

type Responder struct {
	Ctx *app.RequestContext
	LogSaver TraceLogSaver
	LogOut bool
}

func (r Responder) ReplyOK() {
	r.abortWithJSON(OK)
}

func (r Responder) ReplyNOK() {
	r.abortWithJSON(NOK)
}

func (r Responder) ReplyErrorInternal() {
	r.abortWithJSON(ErrorInternal)
}

func (r Responder) ReplyErrorParam() {
	r.abortWithJSON(ErrorParam)
}

func (r Responder) ReplyErrorParam2(msg string) {
	r.abortWithJSON(BaseResp{
		Code: ErrorParam.Code,
		Msg:  fmt.Sprintf("%s:%s", ErrorParam.Msg, msg),
	})
}

func (r Responder) ReplyErrorPermission() {
	r.abortWithJSON(ErrorPermission)
}

func (r Responder) ReplyErrorIllegal() {
	r.abortWithJSON(ErrorIllegal)
}

func (r Responder) ReplyErrorAuthorization() {
	r.abortWithJSON(ErrorAuthorization)
}

func (r Responder) ReplyErrorUserNull() {
	r.abortWithJSON(ErrorUserNull)
}

func (r Responder) ReplyErrorUserExists() {
	r.abortWithJSON(ErrorUserExists)
}

func (r Responder) ReplyErrorUserBlocked() {
	r.abortWithJSON(ErrorUserBlocked)
}

func (r Responder) ReplyErrorUserCanceled() {
	r.abortWithJSON(ErrorUserCanceled)
}

func (r Responder) ReplyErrMsg(msg string) {
	r.abortWithJSON(BaseResp{
		Code: NOK.Code,
		Msg:  msg,
	})
}

func (r Responder) ReplyErr(err BaseResp) {
	r.abortWithJSON(BaseResp{
		Code: err.Code,
		Msg:  err.Msg,
	})
}

func (r Responder) ReplyData(data interface{}) {
	raw, _ := json.Marshal(data)
	r.abortWithJSON(BaseResp{
		Code: OK.Code,
		Msg:  OK.Msg,
		Data: raw,
	})
}

func (r Responder) Reply(code int32, msg string, data interface{}) {
	raw, _ := json.Marshal(data)
	r.abortWithJSON(BaseResp{
		Code: code,
		Msg:  msg,
		Data: raw,
	})
}


func (r Responder) abortWithJSON(resp BaseResp){
	if r.LogOut {
		params := map[string]interface{}{}
		r.Ctx.Bind(&params)
		logger.Infof("path: %v, params: %v, resp: %v", r.Ctx.FullPath(), utils.JSONString(params), utils.JSONString(resp))
	}
	if r.LogSaver != nil {
		r.LogSaver.AddLog(r.Ctx, resp)
	}
	r.Ctx.AbortWithStatusJSON(http.StatusOK,resp)
}