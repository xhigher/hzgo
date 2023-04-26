package alipay

import (
	"context"
	"github.com/go-pay/gopay"
	"github.com/go-pay/gopay/alipay"
	"net/http"
	"time"
	"github.com/xhigher/hzgo/config"
	"github.com/xhigher/hzgo/env"
	"github.com/xhigher/hzgo/logger"
	"github.com/xhigher/hzgo/utils"
)

type AliPayClient struct {
	client *alipay.Client
}

var (
	alipayConfig *config.AliPayConfig
	alipayClient *AliPayClient
)

func Init(conf *config.AliPayConfig) {
	alipayConfig = conf
	alipayClient = newClient(conf.AppId, conf.AppPrivateKey, conf.PublicKey, conf.IsProd, conf.NotifyUrl, conf.ReturnUrl)
}

func newClient(appId, appPrivateKey, publicKey string, isProd bool, notifyUrl, returnUrl string) *AliPayClient {
	// 初始化支付宝客户端
	//    appId：应用ID
	//    privateKey：应用私钥，支持PKCS1和PKCS8
	//    isProd：是否是正式环境
	client, err := alipay.NewClient(appId, appPrivateKey, isProd)
	if err != nil {
		logger.Errorf("error %v", err)
		return nil
	}
	// 打开Debug开关，输出日志，默认关闭
	client.DebugSwitch = gopay.DebugOff

	// 设置支付宝请求 公共参数
	//    注意：具体设置哪些参数，根据不同的方法而不同，此处列举出所有设置参数
	client.SetLocation(alipay.LocationShanghai). // 设置时区，不设置或出错均为默认服务器时间
							SetCharset(alipay.UTF8).  // 设置字符编码，不设置默认 utf-8
							SetSignType(alipay.RSA2). // 设置签名类型，不设置默认 RSA2
							SetReturnUrl(returnUrl).  // 设置返回URL
							SetNotifyUrl(notifyUrl).  // 设置异步通知URL
							SetAppAuthToken("")       // 设置第三方应用授权

	// 自动同步验签（只支持证书模式）
	// 传入 alipayCertPublicKey_RSA2.crt 内容
	client.AutoVerifySign([]byte(publicKey))

	// 公钥证书模式，需要传入证书，以下两种方式二选一
	// 证书路径
	//err = client.SetCertSnByPath("appCertPublicKey.crt", "alipayRootCert.crt", "alipayCertPublicKey_RSA2.crt")
	// 证书内容
	//certContent := []byte("appCertPublicKey bytes")
	//rootCertContent := []byte("aliPayRootCertContent bytes")
	//pubCertContent := []byte("aliPayPublicCertContent bytes")
	//err = client.SetCertSnByContent(certContent, rootCertContent, pubCertContent)
	//if err != nil {
	//	xlog.Error(err)
	//	return nil
	//}
	return &AliPayClient{
		client: client,
	}
}

func CreateTrade(orderId string, money int64, subject string) (param string, err error) {
	if !env.IsProd() {
		money = 1
	}
	// 初始化 BodyMap
	expireTime := time.Now().Add(2 * time.Minute).Format(utils.TimeFormatYMDHMS)
	bm := make(gopay.BodyMap)
	bm.Set("subject", subject).
		Set("out_trade_no", orderId).
		Set("total_amount", utils.FormatMoney(money)).
		Set("time_expire", expireTime)
	param, err = alipayClient.client.TradeAppPay(context.Background(), bm)
	if err != nil {
		logger.Errorf("err: %v", err)
		return
	}

	return
}

func ParseNotifyRequest(req *http.Request) (result gopay.BodyMap, err error) {
	result, err = alipay.ParseNotifyToBodyMap(req)
	if err != nil {
		return
	}

	_, err = alipay.VerifySign(alipayConfig.PublicKey, result)
	if err != nil {
		return
	}
	return
}

func QueryOrder(orderId string) (result *alipay.TradeQuery, err error) {
	bm := make(gopay.BodyMap)
	bm.Set("out_trade_no", orderId)
	resp, err := alipayClient.client.TradeQuery(context.Background(), bm)
	if err != nil {
		if resp != nil {
			if resp.Response.Code == "40004" && resp.Response.SubCode == "ACQ.TRADE_NOT_EXIST" {
				err = nil
				result = &alipay.TradeQuery{
					TradeStatus: "TRADE_NOT_EXIST",
				}
			}
		}
		return
	}
	result = resp.Response
	return
}

func QueryOrderState(orderId string) (state string, err error) {
	orderInfo, err := QueryOrder(orderId)
	if err != nil {
		return
	}
	state = orderInfo.TradeStatus
	return
}
