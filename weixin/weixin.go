package weixin

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/xhigher/hzgo/httpcli"
	"github.com/xhigher/hzgo/logger"
	"github.com/xhigher/hzgo/utils"
	"time"
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
	Timeout time.Duration `json:"timeout"`
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

type UserInfo struct {
	OpenId    string `json:"openId"`
	NickName  string `json:"nickName"`
	Gender    int    `json:"gender"`
	City      string `json:"city"`
	Province  string `json:"province"`
	Country   string `json:"country"`
	AvatarUrl string `json:"avatarUrl"`
	UnionId   string `json:"unionId"`
	Watermark struct {
		Appid     string      `json:"appid"`
		Timestamp interface{} `json:"timestamp"`
	} `json:"watermark"`
}

type WeixinManager struct {
	conf *Configs
}

var mgr *WeixinManager

func Init(conf *Configs) {
	mgr = &WeixinManager{
		conf: conf,
	}
}

func Code2Session(code string) (resp *Code2SessionResp, err error) {
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

	err = httpcli.GetJSONWithTimeout(code2SessionUrl, data, mgr.conf.Timeout, &resp)
	if err != nil {
		logger.Errorf("get data: %v, error: %v", data, err)
		return
	}
	return
}

func DecryptMiniAppUserData(sessionKey, encryptedData, iv string) (data *UserInfo, err error) {
	decodeBytes, err := base64.StdEncoding.DecodeString(encryptedData)
	if err != nil {
		logger.Errorf("error: %v", err)
		return
	}
	sessionKeyBytes, errKey := base64.StdEncoding.DecodeString(sessionKey)
	if errKey != nil {
		logger.Errorf("error: %v", err)
		return
	}
	ivBytes, errIv := base64.StdEncoding.DecodeString(iv)
	if errIv != nil {
		logger.Errorf("error: %v", err)
		return
	}
	dataBytes, err := utils.AesDecrypt(decodeBytes, sessionKeyBytes, ivBytes)
	if err != nil {
		logger.Errorf("error: %v", err)
		return
	}

	err = json.Unmarshal(dataBytes, &data)
	if err != nil {
		logger.Errorf("error: %v", err)
		return
	}
	return
}

func GetOauth2AccessToken(typ int, code string) (resp *OAuth2AccessTokenResp, err error) {
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

	err = httpcli.GetJSONWithTimeout(oauth2AccessTokenUrl, data, mgr.conf.Timeout, &resp)
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

func GetOauth2UserInfo(token, openid string) (resp *OAuth2UserInfoResp, err error) {
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

	err = httpcli.GetJSONWithTimeout(oauth2UserInfoUrl, data, mgr.conf.Timeout, &resp)
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
	err = httpcli.GetJSONWithTimeout(ticketUrl, data, mgr.conf.Timeout, resp)
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

func GetServerAccessToken() (token string, err error) {
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
	err = httpcli.GetJSONWithTimeout(serverAccessTokenUrl, data, mgr.conf.Timeout, resp)
	if err != nil {
		logger.Errorf("get data: %v, error: %v", data, err)
		return
	}
	if resp.Errcode != 0 {
		err = errors.New(fmt.Sprintf("%d:%s", resp.Errcode, resp.Errmsg))
		return
	}
	token = resp.AccessToken
	return
}

func GetServerUserInfo(token, openid string) (resp *ServerUserInfoResp, err error) {
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

	err = httpcli.GetJSONWithTimeout(serverUserInfoUrl, data, mgr.conf.Timeout, &resp)
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
