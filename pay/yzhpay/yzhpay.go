package yzhpay

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/xhigher/hzgo/logger"
	"github.com/xhigher/hzgo/utils"
)

var (
	yzhConfig *YZHConfig
	yzhClient *Client
)

type YZHConfig struct {
	BrokerID     string // 代征主体ID
	DealerID     string // 商户ID
	Appkey       string // 商户appkey
	Des3Key      string // 商户des3key
	PrivateKey   string // 商户秘钥
	PublicKey    string //商户公钥
	YunPublicKey string // 云账户公钥
	NotifyURL    string
}

type BankAccount struct {
	Cardno   string `json:"cardno"`
	Phoneno  string `json:"phoneno"`
	Realname string `json:"realname"`
	Idcard   string `json:"idcard"`
	Bank     string `json:"bank"`
}

type AlipayAccount struct {
	Userno   string `json:"userno"`
	Name     string `json:"name"`
	Realname string `json:"realname"`
	Idcard   string `json:"idcard"`
}

func Init(conf *YZHConfig) {
	yzhConfig = conf
	yzhClient = NewClient(conf.BrokerID, conf.DealerID, conf.Appkey, conf.Des3Key, conf.PrivateKey, conf.PublicKey, conf.YunPublicKey)
}

/**
* 银行卡下单,响应云账户综合服务平台订单流水号
 */
func CreateBankOrder(account BankAccount, orderid string, money int64) (ref string, err error) {
	bankOrderParam := &BankOrderParam{
		CardNo:  account.Cardno,
		PhoneNo: account.Phoneno,
	}

	bankOrderParam.OrderID = orderid
	bankOrderParam.RealName = account.Realname
	bankOrderParam.IDCard = account.Idcard
	bankOrderParam.Pay = utils.FormatMoney(money)
	bankOrderParam.PayRemark = "银行卡打款"
	bankOrderParam.NotifyURL = yzhConfig.NotifyURL

	ref, err = yzhClient.CreateBankOrder(bankOrderParam)
	if err != nil {
		logger.Errorf("CreateBankOrder err: %s, %v", orderid, err)
		return
	}
	logger.Warnf("CreateBankOrder: %s, %s, %d, %s", account.Cardno, orderid, money, ref)
	return
}

/**
* 支付宝下单,响应云账户综合服务平台订单流水号
 */
func CreateAlipayOrder(account AlipayAccount, orderid string, money int64) (ref string, err error) {
	aliOrderParam := &AliOrderParam{
		CardNo:    account.Userno,
		CheckName: account.Realname,
	}
	aliOrderParam.RealName = account.Realname
	aliOrderParam.IDCard = account.Idcard
	aliOrderParam.OrderID = orderid
	aliOrderParam.Pay = utils.FormatMoney(money)
	aliOrderParam.PayRemark = "支付宝打款"
	aliOrderParam.NotifyURL = yzhConfig.NotifyURL

	ref, err = yzhClient.CreateAliOrder(aliOrderParam)
	if err != nil {
		logger.Errorf("CreateAlipayOrder err: %s, %v", orderid, err)
		return
	}
	logger.Warnf("CreateAlipayOrder: %s, %s, %d, %s", account.Userno, orderid, money, ref)
	return
}

/**
* 查询订单信息,响应云账户综合服务平台订单信息
 */
func QueryOrder(orderid string) (data OrderInfo, err error) {
	data, err = yzhClient.QueryOrder(orderid, "支付宝", "")
	if err != nil {
		logger.Errorf("QueryOrder err: %s, %v", orderid, err)
		return
	}
	return
}

/**
* 取消商户未打款订单,响应是否取消成功
 */
func CancelOrder(orderid string) (ok bool, err error) {
	ok, err = yzhClient.CancelOrder(orderid, "", "支付宝")
	if err != nil {
		logger.Errorf("CancelOrder err: %s, %v", orderid, err)
		return
	}
	return
}

/**
* 查询商户充值记录,响应商户充值记录信息
 */
func QueryRechargeRecords(startDate, endDate string) (records []RechargeRecord, err error) {
	records, err = yzhClient.QueryRechargeRecord(startDate, endDate)
	if err != nil {
		logger.Errorf("QueryRechargeRecords err: %s, %s, %v", startDate, endDate, err)
		return
	}
	return
}

/**
* 	银行卡四要素鉴权发送短信上行接口,响应云账户综合服务平台银行卡四要素流水号
 */
func VerifyBankAccountSendMsgCode(account BankAccount) (ref string, err error) {
	ref, err = yzhClient.ElementVerifyRequest(account.Idcard, account.Realname, account.Cardno, account.Phoneno)
	if err != nil {
		logger.Errorf("VerifyBankAccountSendMsgCode err: %v, %v", account, err)
		return
	}
	return
}

/**
* 	银行卡四要素鉴权提交验证码确认接口,响应银行卡四要素校验是否通过
 */
func VerifyUserBankAccountConfirm(account BankAccount, ref, code string) {
	ok, err := yzhClient.ElementVerifyConfirm(account.Idcard, account.Realname, account.Cardno, account.Phoneno, ref, code)
	if err != nil {
		logger.Errorf("VerifyUserBankAccountConfirm err: %v, %v", account, err)
		return
	}
	fmt.Println(ok)
}

/**
* 	银行卡四要素鉴权接口,响应银行卡四要素鉴权是否通过
* 银行卡三要素鉴权接口,响应银行卡三要素鉴权是否通过
 */
func VerifyBankAccountResult(account BankAccount, element4 bool) (ok bool, err error) {
	if element4 {
		ok, err = yzhClient.Element4Check(account.Idcard, account.Realname, account.Cardno, account.Phoneno)
	} else {
		ok, err = yzhClient.Element3Check(account.Idcard, account.Realname, account.Cardno)
	}
	if err != nil {
		logger.Errorf("VerifyBankAccountResult err: %v, %v, %v", account, element4, err)
		return
	}
	return
}

/**
* 实名制二要素鉴权接口,响应实名制二要素鉴权是否通过
 */
func CheckIDCard(idcard, realname string) (ok bool, err error) {
	ok, err = yzhClient.IDCheck(idcard, realname)
	if err != nil {
		logger.Errorf("CheckIDCard err: %s, %s, %v", idcard, realname, err)
		return
	}
	return
}

/**
* 	查询银行卡信息接口,响应银行卡具体信息以及云账户综合服务平台是否支持该银行卡打款
 */
func CheckBankAccount(cardno, bank string) (cardInfo BankCardInfo, err error) {
	cardInfo, err = yzhClient.QueryBankCardInfo(cardno, bank)
	if err != nil {
		logger.Errorf("CheckBankAccount err: %s, %s, %v", cardno, bank, err)
		return
	}
	return
}

func ParseNotifyRequest(ctx *gin.Context) (result *OrderNotifyResult, err error) {
	params := &OrderNotifyParams{}
	err = ctx.Bind(params)
	if err != nil {
		return
	}

	result, err = yzhClient.OrderCallBack(params)
	return
}
