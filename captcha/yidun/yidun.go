package yidun

import (
	"fmt"
	"sort"
	"strings"
	"time"
	"github.com/xhigher/hzgo/httpcli"
	"github.com/xhigher/hzgo/logger"
	"github.com/xhigher/hzgo/utils"
)

const (
	verifyUrl = "https://c.dun.163.com/api/v2/verify"

	errCodeOK      = 0   // 无异常
	errCodeSign    = 415 // 签名校验错误
	errCodeParam   = 419 // 参数校验错误，例如参数类型错误、参数值错误、必填项为空等
	errCodeQPS     = 430 // qps超限
	errCodeVersion = 421 // 验证码版本不匹配

)

type VerifyParam struct {
	Version   string `json:"version"`   //是	版本信息，固定值v2
	Nonce     string `json:"nonce"`     //是	用于防重放。最长32个字符
	Timestamp string `json:"timestamp"` //是	发起此查询操作的时间。UNIX时间戳，单位：毫秒
	SecretId  string `json:"secretId"`  //是	产品密钥id，由易盾验证码服务分配
	Signature string `json:"signature"` //是	此次请求的签名，用来验证请求的合法性。具体算法见 接口鉴权
	CaptchaId string `json:"captchaId"` //是	32	验证码id
	Validate  string `json:"validate"`  //是 不限制长度，建议1024 提交二次校验的验证数据，即NECaptchaValidate值。只能校验成功一次，重复校验会返回校验不通过。有效时长默认20分钟，可在官网配置为1~20分钟
	User      string `json:"user"`      //是	32	用户信息，该字段必传，值可为空
}
type VerifyResult struct {
	Result      bool   `json:"result"`
	Error       int    `json:"error"`
	Msg         string `json:"msg"`
	Phone       string `json:"phone"`
	ExtraData   string `json:"extraData"`
	CaptchaType int    `json:"captchaType"`
	Token       string `json:"token"`
	SdkReduce   bool   `json:"sdkReduce"`
}

type YiDunCaptcha struct {
	cli       *httpcli.HttpCli
	secretId  string
	secretKey string
}

var sp *YiDunCaptcha

func Init(secretId, secretKey string) {
	sp = &YiDunCaptcha{
		cli:       httpcli.New(5 * time.Second),
		secretId:  secretId,
		secretKey: secretKey,
	}
}

func Verify(code, validate, user string) (err error) {
	data := map[string]string{
		"version":   "v2",
		"nonce":     utils.RandString(20),
		"timestamp": fmt.Sprintf("%d", utils.NowTimeMillis()),
		"secretId":  sp.secretId,
		"signature": "",
		"captchaId": code,
		"validate":  validate,
		"user":      user,
	}

	signature, err := sp.createSignature(data)
	if err != nil {
		logger.Errorf("Verify signature data:%v, error: %v", data, err)
		return
	}
	data["signature"] = signature

	result := &VerifyResult{}
	err = sp.cli.PostForm(verifyUrl, data, result)
	if err != nil {
		logger.Errorf("Verify post data:%v, error: %v", data, err)
		return
	}
	return
}

func (s *YiDunCaptcha) createSignature(params map[string]string) (signature string, err error) {
	keys := make([]string, 0)
	for k := range params {
		if k == "signature" {
			continue
		}
		keys = append(keys, k)
	}
	sort.Strings(keys)
	signItems := make([]string, 0)
	for _, key := range keys {
		item := fmt.Sprintf("%s%v", key, params[key])
		signItems = append(signItems, item)
	}
	signStr := strings.Join(signItems, "")
	signature = utils.MD5(fmt.Sprintf("%s%s", signStr, s.secretKey))
	return
}
