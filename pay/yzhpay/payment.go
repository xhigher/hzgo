package yzhpay

import (
	"encoding/json"
	"fmt"
	"log"
)

/**
 *  CreateBankOrder 银行卡下单
 * pay  必填 （注意是字符串类型） 支付金额  用户实际到账金额，单位：元，支持小数点后两位
 * pay_remark  非必填 支付备注 最大20个字符 不支持特殊字符 ' " & | @% ( ) -
 * 同步响应结果中code=0000不代表支付成功，只表示接单成功，支付结果以异步通知为准
 */
func (c *Client) CreateBankOrder(param *BankOrderParam) (ref string, err error) {
	fn := "CreateBankOrder"
	order := map[string]string{
		"order_id":   param.OrderID,
		"broker_id":  c.BrokerID,
		"dealer_id":  c.DealerID,
		"real_name":  param.RealName,
		"id_card":    param.IDCard,
		"card_no":    param.CardNo,
		"phone_no":   param.PhoneNo,
		"pay":        param.Pay,
		"pay_remark": param.PayRemark,
		"notify_url": param.NotifyURL,
	}

	params, err := BuildParams(order, c.Appkey, c.Des3Key, c.PrivateKey)
	if err != nil {
		log.Printf("@[%s] BuildParams failed, err=%v, order=%+v", fn, err, order)
		return
	}

	headers := BuildHeader(c.DealerID)

	log.Printf("@[%s] begin, params=%+v, header=%v", fn, params, headers)

	result, err := Post(c.Gateway+BankOrderURL, params, headers)
	if err != nil {
		log.Printf("@[%s] Post failed, err=%v, params=%+v, headers=%+v", fn, err, params, headers)
		return
	}
	log.Printf("@[%s] end, params=%+v, header=%v, result=%s", fn, params, headers, string(result))

	var res CreateOrderResponse
	err = json.Unmarshal(result, &res)
	if err != nil {
		log.Printf("@[%s] json.Unmarshal failed, err=%v, result=%s, params=%+v, header=%+v", fn, err, string(result), params, headers)
		return
	}
	// 请求成功
	if res.Code == "0000" {
		ref = res.Data.Ref
		// 支付订单接收成功，尚未处理
		// TODO 比如可记录云账户综合服务平台流水号ref
		return
	} else if res.Code == "2002" {
		// 原订单已下单成功
		// TODO 具体订单结果需等待结算异步通知，或主动调用订单查询接口
	} else {
		// TODO 根据返回的message处理订单，如需重试请使用原订单号
	}

	err = fmt.Errorf("Err:%s", res.Message)
	return
}

/**
* CreateAliOrder  支付宝下单
* pay  必填 （注意是字符串类型） 支付金额  用户实际到账金额，单位：元，支持小数点后两位
* check_name  非必填  校验支付宝姓名  校验支付宝姓名，固定值：Check
* pay_remark  非必填 支付备注 会显示在支付宝账单中的“理由”中，账单备注默认显示 为“云账户”；
   若需修改账单备注，请登录云账户综合服务平台，在【商户中 心→商户管理→提醒设置→个性设置】中配置。
* 同步响应结果中code=0000不代表支付成功，只表示接单成功，支付结果以异步通知为准
*/
func (c *Client) CreateAliOrder(param *AliOrderParam) (ref string, err error) {
	fn := "CreateAliOrder"
	order := map[string]string{
		"order_id":   param.OrderID,
		"broker_id":  c.BrokerID,
		"dealer_id":  c.DealerID,
		"real_name":  param.RealName,
		"id_card":    param.IDCard,
		"card_no":    param.CardNo,
		"pay":        param.Pay,
		"pay_remark": param.PayRemark,
		"check_name": param.CheckName,
		"notify_url": param.NotifyURL,
	}

	params, err := BuildParams(order, c.Appkey, c.Des3Key, c.PrivateKey)
	if err != nil {
		log.Printf("@[%s] BuildParam failed, err=%v, order=%+v", fn, err, order)
		return
	}

	headers := BuildHeader(c.DealerID)

	log.Printf("@[%s] begin, params=%+v, header=%v", fn, params, headers)

	result, err := Post(c.Gateway+AliOrderURL, params, headers)
	if err != nil {
		log.Printf("@[%s] Post failed, err=%v, params=%+v, headers=%+v", fn, err, params, headers)
		return
	}
	log.Printf("@[%s] end, data=%+v, header=%v, result=%s", fn, params, headers, string(result))

	var res CreateOrderResponse
	err = json.Unmarshal(result, &res)
	if err != nil {
		log.Printf("@[%s] json.Unmarshal failed, err=%v, params=%s, data=%+v, header=%+v", fn, err, string(result), params, headers)
		return
	}

	// 请求成功
	if res.Code == "0000" {
		ref = res.Data.Ref
		// 支付订单接收成功，尚未处理
		// TODO 比如可记录云账户综合服务平台流水号ref
		return
	} else if res.Code == "2002" {
		// 原订单已下单成功
		// TODO 具体订单结果需等待结算异步通知，或主动调用订单查询接口
	} else {
		// TODO 根据返回的message处理订单，如需重试请使用原订单号
	}

	err = fmt.Errorf("Err:%s", res.Message)
	return
}

/**
 * CreateWxOrder 微信实时下单（企业付款至用户微信零钱）
 * pay  必填 （注意是字符串类型） 支付金额  用户实际到账金额，单位：元，支持小数点后两位
 * wx_app_id  非必填 商户微信AppID  若商户配置了多个 AppID，此处需指定支付对应的AppID
 * wxpay_mode  必填 微信支付模式   固定值：transfer
 * 同步响应结果中code=0000不代表支付成功，只表示接单成功，支付结果以异步通知为准
 */
func (c *Client) CreateWxOrder(param *WxOrderParam) (ref string, err error) {
	fn := "CreateWxOrder"
	order := map[string]string{
		"order_id":   param.OrderID,
		"broker_id":  c.BrokerID,
		"dealer_id":  c.DealerID,
		"real_name":  param.RealName,
		"id_card":    param.IDCard,
		"openid":     param.OpenID,
		"pay":        param.Pay,
		"pay_remark": param.PayRemark,
		"notify_url": param.NotifyURL,
		"wx_app_id":  param.WxAppID,
		"wxpay_mode": param.WxPayMode,
	}

	params, err := BuildParams(order, c.Appkey, c.Des3Key, c.PrivateKey)
	if err != nil {
		log.Printf("@[%s] BuildParam failed, err=%v, order=%+v", fn, err, order)
		return
	}

	headers := BuildHeader(c.DealerID)

	log.Printf("@[%s] begin, params=%+v, header=%v", fn, params, headers)

	result, err := Post(c.Gateway+WxOrderURL, params, headers)
	if err != nil {
		log.Printf("@[%s] Post failed, err=%v, params=%+v, headers=%+v", fn, err, params, headers)
		return
	}
	log.Printf("@[%s] end, data=%+v, header=%v, result=%s", fn, params, headers, string(result))

	var res CreateOrderResponse
	err = json.Unmarshal(result, &res)
	if err != nil {
		log.Printf("@[%s] json.Unmarshal failed, err=%v, params=%s, data=%+v, header=%+v", fn, err, string(result), params, headers)
		return
	}
	// 请求成功
	if res.Code == "0000" {
		ref = res.Data.Ref
		// 支付订单接收成功，尚未处理
		// TODO 比如可记录云账户综合服务平台流水号ref
		return
	} else if res.Code == "2002" {
		// 原订单已下单成功
		// TODO 具体订单结果需等待结算异步通知，或主动调用订单查询接口
	} else {
		// TODO 根据返回的message处理订单，如需重试请使用原订单号
	}

	err = fmt.Errorf("Err:%s", res.Message)
	return
}

/**
 * 查询订单信息
 * data_type 非必填 数据类型 如果为 encryption，则对返回的 data 进行加密
 * ①各个通道支付，仅需要根据订单状态码status即可判断订单是否支付成功：1 成功，2或9 失败，15 取消。
 * ②对于状态1（已支付），在无退汇情况下是最终状态（退汇存在于银行卡通道，用户银行卡为II/III类银行卡收款超限额，导致支付先成功后退汇，此类情况较少）
 */
func (c *Client) QueryOrder(orderID, channel, dataType string) (dest OrderInfo, err error) {
	fn := "QueryOrder"
	data := map[string]string{
		"order_id":  orderID,
		"channel":   channel,
		"data_type": dataType,
	}

	params, err := BuildParams(data, c.Appkey, c.Des3Key, c.PrivateKey)
	if err != nil {
		log.Printf("@[%s] BuildParam failed, err=%v, data=%+v", fn, err, data)
		return
	}

	headers := BuildHeader(c.DealerID)
	log.Printf("@[%s] begin, params=%+v, header=%v", fn, params, headers)

	result, err := Get(c.Gateway+QueryOrderURL, params, headers)
	if err != nil {
		log.Printf("@[%s] Get failed, err=%v, params=%+v, headers=%+v", fn, err, params, headers)
		return
	}
	log.Printf("@[%s] end, data=%+v, header=%v, result=%s", fn, params, headers, string(result))

	var res QueryOrderResponse
	err = json.Unmarshal(result, &res)
	if err != nil {
		log.Printf("@[%s] json.Unmarshal failed, err=%v, params=%s, data=%+v, header=%+v", fn, err, string(result), params, headers)
		return
	}
	log.Printf(" res-----=%+v, status=%v ", res, res.Data.Status)
	// 请求成功
	if res.Code == "0000" {
		dest = res.Data
		var status = dest.Status
		// 主动查询订单订单状态时status值分析：
		// status == "-1" 订单被删除，只有通过Web页面支付的情况才会出现（最终状态，此状态无需做处理）
		// status == "0" 支付订单接收成功，尚未处理（中间状态，此状态无需做处理）
		// status == "1" 支付成功，但是银行卡可能存在退汇情况，会导致支付先成功后退汇（对于银行卡支付是中间状态，对于支付宝和微信支付是最终状态，此状态无需做处理）
		// status == "2" 支付失败，订单数据校验不通过（最终状态，比如银行卡三要素不通过时需要去核验用户信息或卡信息）
		// status == "4" 支付挂单（中间状态，需要查看挂单具体原因，比如由于余额不足导致的挂单，在72小时内补足余额后可自动继续支付）
		// status == "5" 支付中，调用支付网关超时等状态异常情况导致，处于等待交易查证的中间状态（中间状态，此状态无需做处理）
		// status == "8" 待支付，订单支付限额检查和风控判断完毕，等待执行支付的状态（中间状态，此状态无需做处理）
		// status == "9" 退汇， 只有银行卡提现时可能会出现（银行卡支付的最终状态，需核实用户卡信息）
		// status == "15" 取消支付，表示待支付（暂停处理）订单数据被商户主动取消（最终状态，此状态无需做处理）

		if status == "0" {
			// 支付订单接收成功，尚未处理（中间状态）
			// TODO 比如可更新当前订单状态为“已接单”,记录云账户综合服务平台流水号ref
		} else if status == "1" {
			// 支付成功，但是银行卡可能存在退汇情况，会导致支付先成功后退汇（对于银行卡支付是中间状态，对于支付宝和微信支付是最终状态）
			// TODO 更新当前订单状态,记录云账户综合服务平台流水号ref
		} else if status == "2" {
			// 支付失败，订单数据校验不通过（最终状态，比如银行卡三要素不通过时需要去核验用户信息或卡信息）
			// TODO 最终状态，更新当前订单状态,记录云账户综合服务平台流水号ref
		} else if status == "4" {
			// 支付挂单（中间状态，需要查看挂单具体原因，比如由于余额不足导致的挂单，在72小时内补足余额后可自动继续支付）
			// TODO 比如可更新当前订单状态为“挂单”,记录云账户综合服务平台流水号ref
		} else if status == "5" {
			// 支付中，调用支付网关超时等状态异常情况导致，处于等待交易查证的中间状态（中间状态，此状态无需做处理）
		} else if status == "8" {
			// 待支付，订单支付限额检查和风控判断完毕，等待执行支付的状态（中间状态，此状态无需做处理）
		} else if status == "9" {
			// 退汇， 只有银行卡提现时可能会出现（银行卡支付的最终状态，需核实用户卡信息）
			// TODO 订单终态，更新当前订单状态为“退汇”
		} else if status == "15" {
			// 取消支付，最终状态 表示待支付（暂停处理）订单数据被商户主动取消
			// TODO 订单终态，比如更新当前订单状态为“已取消”
		}
		return
	} else if res.Code == "2018" {
		// 订单不存在，需要先检查入参是否正确（如channel是否为“支付宝”、“微信”、“银行卡”）
		// TODO 如有重新下单需求，务必使用原订单号order_id下单
	} else {
		// 其他响应码均为异常状态，异常状态代表订单不确定是否存在，请稍后重新查询或联系云账户人工查询确认订单状态
	}
	err = fmt.Errorf("Err:%s", res.Message)
	return
}

// QueryAccountBalance 查询账户余额
func (c *Client) QueryAccountBalance() (accounts []AccountBalance, err error) {
	fn := "QueryAccountBalance"
	data := map[string]string{
		"dealer_id": c.DealerID,
	}

	params, err := BuildParams(data, c.Appkey, c.Des3Key, c.PrivateKey)
	if err != nil {
		log.Printf("@[%s] BuildParam failed, err=%v, data=%+v", fn, err, data)
		return
	}

	headers := BuildHeader(c.DealerID)
	log.Printf("@[%s] begin, params=%+v, header=%v", fn, params, headers)

	result, err := Get(c.Gateway+QueryAccountURL, params, headers)
	if err != nil {
		log.Printf("@[%s] Get failed, err=%v, params=%+v, headers=%+v", fn, err, params, headers)
		return
	}
	log.Printf("@[%s] end, data=%+v, header=%v, result=%s", fn, params, headers, string(result))

	var res QueryAccountBalanceResponse
	err = json.Unmarshal(result, &res)
	if err != nil {
		log.Printf("@[%s] json.Unmarshal failed, err=%v, params=%s, data=%+v, header=%+v", fn, err, string(result), params, headers)
		return
	}

	// 请求成功
	if res.Code == "0000" {
		accounts = res.Data.DealerInfos
		return
	}
	err = fmt.Errorf("Err:%s", res.Message)
	return
}

/**
 * QueryReceiptFile 查询电子回单
 * @param {*} queryObj {order_id, ref}
 * order_id 商户订单号  与平台流水号不能同时为空
 * ref 综合服务平台流水号  与商户订单号不能同时为空
 */
func (c *Client) QueryReceiptFile(orderID, ref string) (file OrderReceiptFile, err error) {
	fn := "QueryOrderReceiptFile"
	data := map[string]string{
		"order_id": orderID,
		"ref":      ref,
	}

	params, err := BuildParams(data, c.Appkey, c.Des3Key, c.PrivateKey)
	if err != nil {
		log.Printf("@[%s] BuildParam failed, err=%v, data=%+v", fn, err, data)
		return
	}

	headers := BuildHeader(c.DealerID)
	log.Printf("@[%s] begin, params=%+v, header=%v", fn, params, headers)

	result, err := Get(c.Gateway+QueryReceiptFileURL, params, headers)
	if err != nil {
		log.Printf("@[%s] Get failed, err=%v, params=%+v, headers=%+v", fn, err, params, headers)
		return
	}
	log.Printf("@[%s] end, data=%+v, header=%v, result=%s", fn, params, headers, string(result))

	var res QueryReceiptFileResponse
	err = json.Unmarshal(result, &res)
	if err != nil {
		log.Printf("@[%s] json.Unmarshal failed, err=%v, params=%s, data=%+v, header=%+v", fn, err, string(result), params, headers)
		return
	}

	// 请求成功
	if res.Code == "0000" {
		file = res.Data
		return
	}
	err = fmt.Errorf("Err:%s", res.Message)
	return
}

// CancelOrder 取消订单
func (c *Client) CancelOrder(orderID, ref, channel string) (ok bool, err error) {
	fn := "CancelOrder"
	data := map[string]string{
		"dealer_id": c.DealerID,
		"order_id":  orderID,
		"ref":       ref,
		"channel":   channel,
	}

	params, err := BuildParams(data, c.Appkey, c.Des3Key, c.PrivateKey)
	if err != nil {
		log.Printf("@[%s] BuildParam failed, err=%v, data=%+v", fn, err, data)
		return
	}

	headers := BuildHeader(c.DealerID)
	log.Printf("@[%s] begin, params=%+v, header=%v", fn, params, headers)

	result, err := Post(c.Gateway+CancelOrderURL, params, headers)
	if err != nil {
		log.Printf("@[%s] Post failed, err=%v, params=%+v, headers=%+v", fn, err, params, headers)
		return
	}
	log.Printf("@[%s] end, data=%+v, header=%v, result=%s", fn, params, headers, string(result))

	var res BaseCheckResponse
	err = json.Unmarshal(result, &res)
	if err != nil {
		log.Printf("@[%s] json.Unmarshal failed, err=%v, params=%s, data=%+v, header=%+v", fn, err, string(result), params, headers)
		return
	}

	// 请求成功
	if res.Code == "0000" {
		ok = res.Data.Ok
		return
	}
	err = fmt.Errorf("Err:%s", res.Message)
	return
}

// QueryRechargeRecord 查询充值记录
func (c *Client) QueryRechargeRecord(beginAt, endAt string) (records []RechargeRecord, err error) {
	fn := "QueryRechargeRecord"
	data := map[string]string{
		"begin_at": beginAt,
		"end_at":   endAt,
	}

	params, err := BuildParams(data, c.Appkey, c.Des3Key, c.PrivateKey)
	if err != nil {
		log.Printf("@[%s] BuildParam failed, err=%v, data=%+v", fn, err, data)
		return
	}

	headers := BuildHeader(c.DealerID)
	log.Printf("@[%s] begin, params=%+v, header=%v", fn, params, headers)

	result, err := Get(c.Gateway+QueryRechargeURL, params, headers)
	if err != nil {
		log.Printf("@[%s] Get failed, err=%v, params=%+v, headers=%+v", fn, err, params, headers)
		return
	}
	log.Printf("@[%s] end, data=%+v, header=%v, result=%s", fn, params, headers, string(result))

	var res QueryRechargeRecordResponse
	err = json.Unmarshal(result, &res)
	if err != nil {
		log.Printf("@[%s] json.Unmarshal failed, err=%v, params=%s, data=%+v, header=%+v", fn, err, string(result), params, headers)
		return
	}

	// 请求成功
	if res.Code == "0000" {
		records = res.Data
		return
	}
	err = fmt.Errorf("Err:%s", res.Message)
	return
}

/**
 *  查询商户VA账户信息
 * dealer_id 商户ID
 * broker_id 综合服务主体ID
 */
func (c *Client) QueryVaAccount() (vaAccount VaAccount, err error) {
	fn := "QueryVaAccount"
	data := map[string]string{
		"dealer_id": c.DealerID,
		"broker_id": c.BrokerID,
	}
	params, err := BuildParams(data, c.Appkey, c.Des3Key, c.PrivateKey)
	if err != nil {
		log.Printf("@[%s] BuildParam failed, err=%v, data=%+v", fn, err, data)
		return
	}

	headers := BuildHeader(c.DealerID)
	log.Printf("@[%s] begin, params=%+v, header=%v", fn, params, headers)

	result, err := Get(c.Gateway+QueryVaAccountURL, params, headers)
	if err != nil {
		log.Printf("@[%s] Get failed, err=%v, params=%+v, headers=%+v", fn, err, params, headers)
		return
	}
	log.Printf("@[%s] end, data=%+v, header=%v, result=%s", fn, params, headers, string(result))

	var res QueryVaAccountResponse
	err = json.Unmarshal(result, &res)
	if err != nil {
		log.Printf("@[%s] json.Unmarshal failed, err=%v, params=%s, data=%+v, header=%+v", fn, err, string(result), params, headers)
		return
	}

	// 请求成功
	if res.Code == "0000" {
		vaAccount = res.Data
		return
	}
	err = fmt.Errorf("Err:%s", res.Message)
	return
}
