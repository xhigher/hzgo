package weixin

import (
	"errors"
	"fmt"
	"time"
	"github.com/xhigher/hzgo/httpcli"
	"github.com/xhigher/hzgo/logger"
)

const (
	baseUrl = "https://api.weixin.qq.com"

	//网页授权
	//https://developers.weixin.qq.com/doc/offiaccount/OA_Web_Apps/Wechat_webpage_authorization.html
	oauth2AccessTokenUrl = baseUrl + "/sns/oauth2/access_token"
	oauth2UserInfoUrl    = baseUrl + "/sns/userinfo"

	//小程序登录
	//https://developers.weixin.qq.com/miniprogram/dev/OpenApiDoc/user-login/code2Session.html
	code2SessionUrl = baseUrl + "/sns/jscode2session"

	//公众平台的 API 调用所需的access_token
	//https://developers.weixin.qq.com/doc/offiaccount/Basic_Information/Get_access_token.html
	serverAccessTokenUrl = baseUrl + "/cgi-bin/token"
	serverUserInfoUrl    = baseUrl + "/cgi-bin/user/info"

	//jsapi_ticket是公众号用于调用微信 JS 接口的临时票据
	//https://developers.weixin.qq.com/doc/offiaccount/OA_Web_Apps/JS-SDK.html
	ticketUrl = baseUrl + "/cgi-bin/ticket/getticket"
)

type AppConf struct {
	Appid  string `json:"appid"`
	Secret string `json:"secret"`
}

type Configs struct {
	MiniApp *AppConf `json:"mini_app"`
	MpApp   *AppConf `json:"mp_app"`
}

type Code2SessionResp struct {
	SessionKey string `json:"session_key"` //会话密钥
	Unionid    string `json:"unionid"`     //用户在开放平台的唯一标识符，若当前小程序已绑定到微信开放平台帐号下会返回，详见 UnionID 机制说明。
	Openid     string `json:"openid"`      //用户唯一标识

	Errcode int32  `json:"errcode"` //错误码
	Errmsg  string `json:"errmsg"`  //错误信息
}

type OAuth2AccessTokenResp struct {
	AccessToken    string `json:"access_token"`
	ExpiresIn      int    `json:"expires_in"`
	RefreshToken   string `json:"refresh_token"`
	Openid         string `json:"openid"`
	Scope          string `json:"scope"`
	IsSnapshotuser int    `json:"is_snapshotuser"`
	Unionid        string `json:"unionid"`

	Errcode int32  `json:"errcode"`
	Errmsg  string `json:"errmsg"`
}

type OAuth2UserInfoResp struct {
	Openid     string   `json:"openid"`
	Nickname   string   `json:"nickname"`
	Sex        int      `json:"sex"`
	Province   string   `json:"province"`
	City       string   `json:"city"`
	Country    string   `json:"country"`
	Headimgurl string   `json:"headimgurl"`
	Unionid    string   `json:"unionid"`
	Privilege  []string `json:"privilege"`

	Errcode int32  `json:"errcode"`
	Errmsg  string `json:"errmsg"`
}

type ServerAccessTokenResp struct {
	AccessToken string `json:"access_token"`
	ExpiresIn   int    `json:"expires_in"`

	Errcode int32  `json:"errcode"`
	Errmsg  string `json:"errmsg"`
}

type TicketResp struct {
	Ticket    string `json:"ticket"`
	ExpiresIn int    `json:"expires_in"`

	Errcode int    `json:"errcode"`
	Errmsg  string `json:"errmsg"`
}

type ServerUserInfoResp struct {
	Subscribe      int    `json:"subscribe"`
	Openid         string `json:"openid"`
	Language       string `json:"language"`
	SubscribeTime  int    `json:"subscribe_time"`
	Unionid        string `json:"unionid"`
	Remark         string `json:"remark"`
	Groupid        int    `json:"groupid"`
	TagidList      []int  `json:"tagid_list"`
	SubscribeScene string `json:"subscribe_scene"`
	QrScene        int    `json:"qr_scene"`
	QrSceneStr     string `json:"qr_scene_str"`

	Errcode int    `json:"errcode"`
	Errmsg  string `json:"errmsg"`
}

type WeixinManager struct {
	cli  *httpcli.HttpCli
	conf *Configs
}

var mgr *WeixinManager

func Init(conf *Configs) {
	mgr = &WeixinManager{
		cli:  httpcli.New(5 * time.Second),
		conf: conf,
	}
}

func Code2Session(code string) (err error) {
	if mgr == nil {
		logger.Errorf("error: not init")
		return
	}
	if mgr.conf.MiniApp == nil {
		logger.Errorf("error: mini app config nil")
		return
	}

	data := map[string]string{
		"appid":      mgr.conf.MiniApp.Appid,
		"secret":     mgr.conf.MiniApp.Secret,
		"grant_type": "authorization_code",
		"js_code":    code,
	}

	resp := &Code2SessionResp{}
	err = mgr.cli.GetJSON2(code2SessionUrl, data, resp)
	if err != nil {
		logger.Errorf("get data: %v, error: %v", data, err)
		return
	}
	return
}

func GetOauth2AccessToken(typ int, code string) (err error) {
	if mgr == nil {
		logger.Errorf("error: not init")
		return
	}
	if mgr.conf.MpApp == nil {
		logger.Errorf("error: mp app config nil")
		return
	}

	data := map[string]string{
		"appid":      mgr.conf.MpApp.Appid,
		"secret":     mgr.conf.MpApp.Secret,
		"grant_type": "authorization_code",
		"code":       code,
	}

	resp := &OAuth2AccessTokenResp{}
	err = mgr.cli.GetJSON2(oauth2AccessTokenUrl, data, resp)
	if err != nil {
		logger.Errorf("get data: %v, error: %v", data, err)
		return
	}
	if resp.Errcode != 0 {
		err = errors.New(fmt.Sprintf("%d:%s", resp.Errcode, resp.Errmsg))
		return
	}
	return
}

func GetOauth2UserInfo(token, openid string) (err error) {
	if mgr == nil {
		logger.Errorf("error: not init")
		return
	}
	if mgr.conf.MpApp == nil {
		logger.Errorf("error: mp app config nil")
		return
	}

	data := map[string]string{
		"access_token": token,
		"openid":       openid,
		"lang":         "zh_CN",
	}

	resp := &OAuth2UserInfoResp{}
	err = mgr.cli.GetJSON2(oauth2UserInfoUrl, data, resp)
	if err != nil {
		logger.Errorf("get data: %v, error: %v", data, err)
		return
	}
	if resp.Errcode != 0 {
		err = errors.New(fmt.Sprintf("%d:%s", resp.Errcode, resp.Errmsg))
		return
	}
	return
}

func GetTicket(token, typ string) (ticket string, err error) {
	if mgr == nil {
		logger.Errorf("error: not init")
		return
	}
	if mgr.conf.MpApp == nil {
		logger.Errorf("error: mp app config nil")
		return
	}

	data := map[string]string{
		"access_token": token,
		"type":         typ,
	}

	resp := &TicketResp{}
	err = mgr.cli.GetJSON2(ticketUrl, data, resp)
	if err != nil {
		logger.Errorf("get data: %v, error: %v", data, err)
		return
	}
	if resp.Errcode != 0 {
		err = errors.New(fmt.Sprintf("%d:%s", resp.Errcode, resp.Errmsg))
		return
	}
	ticket = resp.Ticket
	return
}

func GetServerAccessToken(token, openid string) (err error) {
	if mgr == nil {
		logger.Errorf("error: not init")
		return
	}
	if mgr.conf.MpApp == nil {
		logger.Errorf("error: mp app config nil")
		return
	}

	data := map[string]string{
		"appid":      mgr.conf.MpApp.Appid,
		"secret":     mgr.conf.MpApp.Secret,
		"grant_type": "client_credential",
	}

	resp := &ServerAccessTokenResp{}
	err = mgr.cli.GetJSON2(serverAccessTokenUrl, data, resp)
	if err != nil {
		logger.Errorf("get data: %v, error: %v", data, err)
		return
	}
	if resp.Errcode != 0 {
		err = errors.New(fmt.Sprintf("%d:%s", resp.Errcode, resp.Errmsg))
		return
	}
	return
}

func GetServerUserInfo(token, openid string) (err error) {
	if mgr == nil {
		logger.Errorf("error: not init")
		return
	}
	if mgr.conf.MpApp == nil {
		logger.Errorf("error: mp app config nil")
		return
	}

	data := map[string]string{
		"access_token": token,
		"openid":       openid,
		"lang":         "zh_CN",
	}

	resp := &ServerUserInfoResp{}
	err = mgr.cli.GetJSON2(serverUserInfoUrl, data, resp)
	if err != nil {
		logger.Errorf("get data: %v, error: %v", data, err)
		return
	}
	if resp.Errcode != 0 {
		err = errors.New(fmt.Sprintf("%d:%s", resp.Errcode, resp.Errmsg))
		return
	}
	return
}
