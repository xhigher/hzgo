package wxpay

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/go-pay/gopay"
	"github.com/go-pay/gopay/pkg/aes"
	"github.com/go-pay/gopay/pkg/xlog"
	"github.com/go-pay/gopay/wechat/v3"
	"net/http"
	"time"
	"github.com/xhigher/hzgo/config"
	"github.com/xhigher/hzgo/env"
	"github.com/xhigher/hzgo/logger"
	"github.com/xhigher/hzgo/utils"
)

type WxPayClient struct {
	client    *wechat.ClientV3
	notifyUrl string
	supportH5 bool
}

var (
	wxpayConfig *config.WxPayConfig
	wxpayClient *WxPayClient

	wxpayConfigs map[string]*config.WxPayConfig
	wxpayClients map[int]*WxPayClient
)

func Init(conf *config.WxPayConfig) {
	wxpayConfig = conf
	wxpayClient = newClient(conf.MchId, conf.SerialNo, conf.Apiv3Key, conf.PrivateKey, conf.ComplaintUrl)
	if wxpayClient != nil {
		wxpayClient.notifyUrl = conf.NotifyUrl
		wxpayClient.supportH5 = conf.SupportH5
	}
}

func InitMulti(configs []*config.WxPayConfig) {
	wxpayConfigs = make(map[string]*config.WxPayConfig)
	wxpayClients = make(map[int]*WxPayClient)
	for _, conf := range configs {
		wxpayConfigs[conf.MchId] = conf
		complaintUrl := conf.ComplaintUrl
		if conf.Id > 1 {
			complaintUrl = fmt.Sprintf("%s%d", complaintUrl, conf.Id)
		}
		client := newClient(conf.MchId, conf.SerialNo, conf.Apiv3Key, conf.PrivateKey, complaintUrl)
		if client != nil {
			//client.appId = conf.AppId
			if conf.Id > 1 {
				client.notifyUrl = fmt.Sprintf("%s%d", conf.NotifyUrl, conf.Id)
			} else {
				client.notifyUrl = conf.NotifyUrl
			}
			client.supportH5 = conf.SupportH5
		}
		wxpayClients[conf.Id] = client
	}
}

func newClient(mchId, serialNo, apiV3Key, privateKey, complaintUrl string) *WxPayClient {
	//  NewClientV3 初始化微信客户端 v3
	//	mchid：商户ID 或者服务商模式的 sp_mchid
	// 	serialNo：商户证书的证书序列号
	//	apiV3Key：apiV3Key，商户平台获取
	//	privateKey：私钥 apiclient_key.pem 读取后的内容
	client, err := wechat.NewClientV3(mchId, serialNo, apiV3Key, privateKey)
	if err != nil {
		xlog.Errorf("NewClientV3 error: %v", err)
		return nil
	}

	// 启用自动同步返回验签，并定时更新微信平台API证书
	err = client.AutoVerifySign()
	if err != nil {
		xlog.Errorf("AutoVerifySign error: %v", err)
		return nil
	}

	if env.IsProd() {
		if len(complaintUrl) > 0 {
			resp, err2 := client.V3ComplaintNotifyUrlUpdate(context.Background(), complaintUrl)
			if err2 != nil {
				xlog.Errorf("V3ComplaintNotifyUrlUpdate error: %v", err2)
			} else {
				xlog.Infof("V3ComplaintNotifyUrlCreate resp: %v", utils.JSONString(resp))
			}
		}
	}

	// 打开Debug开关，输出日志，默认是关闭的
	client.DebugSwitch = gopay.DebugOn
	return &WxPayClient{
		client: client,
	}
}

func getClient(mchNo int, userid int64) *WxPayClient {
	if wxpayClient != nil {
		return wxpayClient
	}
	clientNum := len(wxpayClients)
	if len(wxpayClients) > 0 {
		if mchNo == 0 {
			mchNo = int(userid)%clientNum + 1
		}
		logger.Infof("wxpay mchNo: %v", mchNo)
		return wxpayClients[mchNo]
	}
	return nil
}

func getClientById(mchNo int) *WxPayClient {
	if mchNo == 1 && wxpayClient != nil {
		return wxpayClient
	}
	logger.Infof("wxpay mchNo: %v", mchNo)
	return wxpayClients[mchNo]
}

func getSupportH5Client() *WxPayClient {
	if wxpayClient != nil {
		if wxpayClient.supportH5 {
			return wxpayClient
		}
	}
	if len(wxpayClients) > 0 {
		for _, client := range wxpayClients {
			if client.supportH5 {
				return client
			}
		}
	}
	return nil
}

func CreateTrade(appId string, mchNo int, orderId string, userid, money int64, description string) (payParams *wechat.AppPayParams, mchid string, err error) {
	if !env.IsProd() {
		money = 1
	}

	mClient := getClient(mchNo, userid)
	if mClient == nil {
		err = errors.New("wxpay client is not init")
		logger.Errorf("CreateTrade error: %v", err)
		return
	}

	// 初始化 BodyMap
	expireTime := time.Now().Add(2 * time.Minute).Format(time.RFC3339)
	bm := make(gopay.BodyMap)
	bm.Set("appid", appId).
		Set("mchid", mClient.client.Mchid).
		Set("out_trade_no", orderId).
		Set("attach", fmt.Sprintf("%d", userid)).
		Set("notify_url", mClient.notifyUrl).
		Set("time_expire", expireTime).
		Set("description", description).
		SetBodyMap("amount", func(bm gopay.BodyMap) {
			bm.Set("total", money).
				Set("currency", "CNY")
		})

	resp, err := mClient.client.V3TransactionApp(context.Background(), bm)
	if err != nil {
		return
	}

	if resp.Code != wechat.Success {
		err = fmt.Errorf("failed to create wxpay order, code: %d, msg: %s", resp.Code, resp.Error)
		return
	}

	payParams, err = mClient.client.PaySignOfApp(appId, resp.Response.PrepayId)
	mchid = mClient.client.Mchid

	return
}

func CreateH5Trade(appId string, orderId string, userid, money int64, description, ip string) (payUrl, mchid string, err error) {
	if !env.IsProd() {
		money = 1
	}

	mClient := getSupportH5Client()
	if mClient == nil {
		err = errors.New("wxpay client is not init")
		logger.Errorf("CreateTrade error: %v", err)
		return
	}

	// 初始化 BodyMap
	bm := make(gopay.BodyMap)
	bm.Set("appid", appId).
		Set("mchid", mClient.client.Mchid).
		Set("out_trade_no", orderId).
		Set("attach", fmt.Sprintf("%d", userid)).
		Set("notify_url", mClient.notifyUrl).
		Set("description", description).
		SetBodyMap("amount", func(m1 gopay.BodyMap) {
			m1.Set("total", money).
				Set("currency", "CNY")
		}).
		SetBodyMap("scene_info", func(m1 gopay.BodyMap) {
			m1.Set("payer_client_ip", ip).
				SetBodyMap("h5_info", func(m2 gopay.BodyMap) {
					m2.Set("type", "Wap")
				})
		})

	resp, err := mClient.client.V3TransactionH5(context.Background(), bm)
	if err != nil {
		return
	}

	if resp.Code != wechat.Success {
		err = fmt.Errorf("failed to create wxpay order, code: %d, msg: %s", resp.Code, resp.Error)
		return
	}

	payUrl = resp.Response.H5Url

	return
}

func ParseNotifyRequest(req *http.Request, id int) (result *wechat.V3DecryptResult, err error) {
	notifyReq, err := wechat.V3ParseNotify(req)
	if err != nil {
		xlog.Infof("V3ParseNotify err: %v", err)
		return
	}

	mClient := getClientById(id)
	if mClient == nil {
		logger.Errorf("wxpay client is not init")
		return
	}

	err = notifyReq.VerifySignByPK(mClient.client.WxPublicKey())
	if err != nil {
		xlog.Infof("VerifySignByPK err: %v", err)
		return
	}

	result, err = notifyReq.DecryptCipherText(string(mClient.client.ApiV3Key))
	if err != nil {
		xlog.Infof("DecryptCipherText err: %v", err)
		return
	}

	return
}

type V3ComplaintResult struct {
	OutTradeNo           string `json:"out_trade_no"`
	ComplaintTime        string `json:"complaint_time"`
	Amount               int    `json:"amount"`
	PayerPhone           string `json:"payer_phone"`
	ComplaintDetail      string `json:"complaint_detail"`
	TransactionId        string `json:"transaction_id"`
	FrozenEndTime        string `json:"frozen_end_time"`
	SubMchid             string `json:"sub_mchid"`
	ComplaintHandleState string `json:"complaint_handle_state"`
	ActionType           string `json:"action_type"`
}

func ParseComplaintRequest(req *http.Request, id int) (result *V3ComplaintResult, mchId string, err error) {
	notifyReq, err := wechat.V3ParseNotify(req)
	if err != nil {
		xlog.Infof("V3ParseNotify err: %v", err)
		return
	}

	mClient := getClientById(id)
	if mClient == nil {
		logger.Errorf("wxpay client is not init")
		return
	}
	mchId = mClient.client.Mchid

	err = notifyReq.VerifySignByPK(mClient.client.WxPublicKey())
	if err != nil {
		xlog.Infof("VerifySignByPK err: %v", err)
		return
	}

	if notifyReq.EventType != "COMPLAINT.CREATE" {
		xlog.Warnf("wxpay complaint event_type: %v", notifyReq.EventType)
		return
	}

	cipherBytes, _ := base64.StdEncoding.DecodeString(notifyReq.Resource.Ciphertext)
	decrypt, err := aes.GCMDecrypt(cipherBytes, []byte(notifyReq.Resource.Nonce), []byte(notifyReq.Resource.AssociatedData), mClient.client.ApiV3Key)
	if err != nil {
		xlog.Infof("GCMDecrypt err: %v", err)
		return
	}
	result = &V3ComplaintResult{}
	if err = json.Unmarshal(decrypt, result); err != nil {
		xlog.Infof("Unmarshal err: %v", err)
		return
	}
	return
}

func QueryOrder(mchNo int, userid int64, orderId string) (orderInfo *wechat.QueryOrder, err error) {
	mClient := getClient(mchNo, userid)
	if mClient == nil {
		err = errors.New("wxpay client is not init")
		logger.Errorf("CreateTrade error: %v", err)
		return
	}
	resp, err := mClient.client.V3TransactionQueryOrder(context.Background(), wechat.OutTradeNo, orderId)
	if err != nil {
		return
	}
	if resp.Code != wechat.Success {
		err = fmt.Errorf("failed to query wxpay order, code: %d, msg: %s", resp.Code, resp.Error)
		return
	}
	orderInfo = resp.Response
	return
}

// 交易状态，枚举值：
//SUCCESS：支付成功,
//REFUND：转入退款,
//NOTPAY：未支付,
//CLOSED：已关闭,
//REVOKED：已撤销（付款码支付）,
//USERPAYING：用户支付中（付款码支付）,
//PAYERROR：支付失败(其他原因，如银行返回失败)
func QueryOrderState(mchNo int, userid int64, orderId string) (state string, err error) {
	orderInfo, err := QueryOrder(mchNo, userid, orderId)
	if err != nil {
		return
	}
	state = orderInfo.TradeState

	return
}
