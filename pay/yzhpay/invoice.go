package yzhpay

import (
	"encoding/json"
	"fmt"
	"log"
)

// QueryInvoice 查询发票信息
func (c *Client) QueryInvoice(year int) (invoice InvoiceInfo, err error) {
	fn := "QueryInvoice"
	data := map[string]interface{}{
		"broker_id": c.BrokerID,
		"dealer_id": c.DealerID,
		"year":      year,
	}

	params, err := BuildParams(data, c.Appkey, c.Des3Key, c.PrivateKey)
	if err != nil {
		log.Printf("@[%s] BuildParam failed, err=%v, data=%+v", fn, err, data)
		return
	}

	headers := BuildHeader(c.DealerID)
	log.Printf("@[%s] begin, params=%+v, header=%v", fn, params, headers)

	result, err := Get(c.Gateway+QueryInvoiceURL, params, headers)
	if err != nil {
		log.Printf("@[%s] Get failed, err=%v, params=%+v, headers=%+v", fn, err, params, headers)
		return
	}
	log.Printf("@[%s] end, data=%+v, header=%v, result=%s", fn, params, headers, string(result))

	var res QueryInvoiceResponse
	err = json.Unmarshal(result, &res)
	if err != nil {
		log.Printf("@[%s] json.Unmarshal failed, err=%v, params=%s, data=%+v, header=%+v", fn, err, string(result), params, headers)
		return
	}

	// 请求成功
	if res.Code == "0000" {
		invoice = res.Data
		return
	}
	err = fmt.Errorf("Err:%s", res.Message)
	return
}

// QueryInvoice 查询可开票额度和开票信息
func (c *Client) QueryInvoiceAmount() (invoiceAmount InvoiceAmount, err error) {
	fn := "QueryInvoiceAmount"
	data := map[string]interface{}{
		"broker_id": c.BrokerID,
		"dealer_id": c.DealerID,
	}

	params, err := BuildParams(data, c.Appkey, c.Des3Key, c.PrivateKey)
	if err != nil {
		log.Printf("@[%s] BuildParam failed, err=%v, data=%+v", fn, err, data)
		return
	}

	headers := BuildHeader(c.DealerID)
	log.Printf("@[%s] begin, params=%+v, header=%v", fn, params, headers)

	result, err := Post(c.Gateway+QueryInvoiceAmountUrl, params, headers)
	if err != nil {
		log.Printf("@[%s] Post failed, err=%v, params=%+v, headers=%+v", fn, err, params, headers)
		return
	}
	log.Printf("@[%s] end, data=%+v, header=%v, result=%s", fn, params, headers, string(result))

	var res QueryInvoiceAmountResponse
	err = json.Unmarshal(result, &res)
	if err != nil {
		log.Printf("@[%s] json.Unmarshal failed, err=%v, params=%s, data=%+v, header=%+v", fn, err, string(result), params, headers)
		return
	}

	// 请求成功
	if res.Code == "0000" {
		invoiceAmount = res.Data
		return
	}
	err = fmt.Errorf("Err:%s", res.Message)
	return
}

// InvoiceApply 开票申请
func (c *Client) InvoiceApply(invoiceApplyParam *InvoiceApplyParam) (invoiceApply InvoiceApply, err error) {
	fn := "InvoiceApply"
	data := map[string]interface{}{
		"broker_id":           c.BrokerID,
		"dealer_id":           c.DealerID,
		"invoice_apply_id":    invoiceApplyParam.InvoiceApplyId,
		"amount":              invoiceApplyParam.Amount,
		"invoice_type":        invoiceApplyParam.InvoiceType,
		"bank_name_account":   invoiceApplyParam.BankNameAccount,
		"goods_services_name": invoiceApplyParam.GoodsServicesName,
		"remark":              invoiceApplyParam.Remark,
	}

	params, err := BuildParams(data, c.Appkey, c.Des3Key, c.PrivateKey)
	if err != nil {
		log.Printf("@[%s] BuildParam failed, err=%v, data=%+v", fn, err, data)
		return
	}

	headers := BuildHeader(c.DealerID)
	log.Printf("@[%s] begin, params=%+v, header=%v", fn, params, headers)

	result, err := Post(c.Gateway+QueryInvoiceAmountUrl, params, headers)
	if err != nil {
		log.Printf("@[%s] Post failed, err=%v, params=%+v, headers=%+v", fn, err, params, headers)
		return
	}
	log.Printf("@[%s] end, data=%+v, header=%v, result=%s", fn, params, headers, string(result))

	var res InvoiceApplyResponse
	err = json.Unmarshal(result, &res)
	if err != nil {
		log.Printf("@[%s] json.Unmarshal failed, err=%v, params=%s, data=%+v, header=%+v", fn, err, string(result), params, headers)
		return
	}

	// 请求成功
	if res.Code == "0000" {
		invoiceApply = res.Data
		return
	}
	err = fmt.Errorf("Err:%s", res.Message)
	return
}

// QueryInvoiceStatus 查询开票申请状态
func (c *Client) QueryInvoiceStatus(invoiceApplyId, applicationId string) (invoiceApply InvoiceApply, err error) {
	fn := "InvoiceApply"
	data := map[string]interface{}{
		"invoice_apply_id": invoiceApplyId,
		"application_id":   applicationId,
	}

	params, err := BuildParams(data, c.Appkey, c.Des3Key, c.PrivateKey)
	if err != nil {
		log.Printf("@[%s] BuildParam failed, err=%v, data=%+v", fn, err, data)
		return
	}

	headers := BuildHeader(c.DealerID)
	log.Printf("@[%s] begin, params=%+v, header=%v", fn, params, headers)

	result, err := Post(c.Gateway+QueryInvoiceAmountUrl, params, headers)
	if err != nil {
		log.Printf("@[%s] Post failed, err=%v, params=%+v, headers=%+v", fn, err, params, headers)
		return
	}
	log.Printf("@[%s] end, data=%+v, header=%v, result=%s", fn, params, headers, string(result))

	var res InvoiceApplyResponse
	err = json.Unmarshal(result, &res)
	if err != nil {
		log.Printf("@[%s] json.Unmarshal failed, err=%v, params=%s, data=%+v, header=%+v", fn, err, string(result), params, headers)
		return
	}

	// 请求成功
	if res.Code == "0000" {
		invoiceApply = res.Data
		return
	}
	err = fmt.Errorf("Err:%s", res.Message)
	return
}
