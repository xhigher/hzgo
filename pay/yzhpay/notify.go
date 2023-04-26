package yzhpay

import (
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"github.com/xhigher/hzgo/consts"
	"github.com/xhigher/hzgo/logger"
)

type OrderNotifyParams struct {
	Data      string `form:"data"`
	Mess      string `form:"mess"`
	Timestamp string `form:"timestamp"`
	Sign      string `form:"sign"`
}

type OrderNotifyResult struct {
	Success bool
	Message string
	OrderId string
	OutId   string
	Money   int64
	Status  int32
	Remark  string
}

// OrderCallBack 订单回调
func (c *Client) OrderCallBack(params *OrderNotifyParams) (result *OrderNotifyResult, err error) {
	result = &OrderNotifyResult{}
	// 获取参数值
	encryptData := params.Data
	mess := params.Mess
	timestamp := params.Timestamp
	sign := params.Sign
	plaintext := fmt.Sprintf("data=%s&mess=%s&timestamp=%s&key=%s", encryptData, mess, timestamp, c.Appkey)
	ok, err := VerifySign(plaintext, sign, c.YunPublicKey)
	if err != nil {
		logger.Errorf("VerifySign failed, data=%s, mess=%s, timestamp=%s, sign=%s", encryptData, mess, timestamp, sign)
		result.Message = "VerifySign failed"
		return
	}
	if !ok {
		logger.Errorf("sign mismatch, data=%s, mess=%s, timestamp=%s, sign=%s", encryptData, mess, timestamp, sign)
		err = errors.New("sign mismatch")
		result.Message = "sign mismatch"
		return
	}

	originData, err := Decrypt(encryptData, c.Des3Key)
	if err != nil {
		logger.Errorf("Decrypt failed, err=%v, data=%s", err, encryptData)
		result.Message = "Decrypt failed"
		return
	}
	logger.Infof("originData=%s", originData)
	var data *OrderCallBackResponse
	err = json.Unmarshal(originData, &data)
	if err != nil {
		logger.Errorf("json.Unmarshal failed, err=%v, data=%s", err, string(originData))
		result.Message = "json.Unmarshal failed"
		return
	}
	logger.Infof("data=%+v", data)

	result.OrderId = data.Data.OrderID
	result.OutId = data.Data.Ref
	result.Money, _ = strconv.ParseInt(data.Data.Pay, 10, 64)

	// 接收订单回调消息通知时订单状态status值分析：
	// status == "1" 支付成功，但是银行卡可能存在退汇情况，会导致支付先成功后退汇（对于银行卡支付是中间状态，对于支付宝和微信支付是最终状态，此状态无需做处理）
	// status == "2" 支付失败，订单数据校验不通过（最终状态，比如银行卡三要素不通过时需要去核验用户信息或卡信息）
	// status == "4" 支付挂单（中间状态，需要查看挂单具体原因，比如由于余额不足导致的挂单，在72小时内补足余额后可自动继续支付）
	// status == "9" 退汇， 只有银行卡提现时可能会出现（银行卡支付的最终状态，需核实用户卡信息）
	// status == "15" 取消支付，表示待支付（暂停处理）订单数据被商户主动取消（最终状态，此状态无需做处理）
	info := "其它原因"
	status, _ := strconv.Atoi(data.Data.Status)
	if status == OrderSuccess {
		// 支付成功，但是银行卡可能存在退汇情况，会导致支付先成功后退汇（对于银行卡支付是中间状态，对于支付宝和微信支付是最终状态）
		result.Status = consts.WithdrawalStatusSuccess
		info = "支付成功"
	} else if status == OrderFailed {
		// 支付失败，订单数据校验不通过（最终状态，比如银行卡三要素不通过时需要去核验用户信息或卡信息）
		result.Status = consts.WithdrawalStatusFailed
		info = "支付失败，订单数据校验不通过"
	} else if status == OrderPending {
		// 支付挂单（中间状态，需要查看挂单具体原因，比如由于余额不足导致的挂单，在72小时内补足余额后可自动继续支付）
		result.Status = consts.WithdrawalStatusWaiting
		info = "支付挂单中，商户余额可能不足"
	} else if status == OrderReturned {
		// 退汇， 只有银行卡提现时可能会出现（银行卡支付的最终状态，需核实用户卡信息）
		result.Status = consts.WithdrawalStatusRefunded
		info = "支付失败，退汇"
	} else if status == OrderCancel {
		// 取消支付，最终状态 表示待支付（暂停处理）订单数据被商户主动取消
		result.Status = consts.WithdrawalStatusCanceled
		info = "支付失败，取消支付"
	}
	result.Remark = fmt.Sprintf("%d:%s", status, info)
	result.Success = true
	result.Message = "success"
	return
}
