package yzhpay

import (
	"encoding/json"
	"fmt"
	"log"
)

// DownloadOrderFile 下载日订单文件
func (c *Client) DownloadOrderFile(orderDate string) (url string, err error) {
	fn := "DownloadOrderFile"
	data := map[string]string{
		"order_date": orderDate,
	}

	params, err := BuildParams(data, c.Appkey, c.Des3Key, c.PrivateKey)
	if err != nil {
		log.Printf("@[%s] BuildParam failed, err=%v, data=%+v", fn, err, data)
		return
	}

	headers := BuildHeader(c.DealerID)
	log.Printf("@[%s] begin, params=%+v, header=%v", fn, params, headers)

	result, err := Get(c.Gateway+DownloadOrderURL, params, headers)
	if err != nil {
		log.Printf("@[%s] Get failed, err=%v, params=%+v, headers=%+v", fn, err, params, headers)
		return
	}
	log.Printf("@[%s] end, data=%+v, header=%v, result=%s", fn, params, headers, string(result))

	var res DownloadOrderFileResponse
	err = json.Unmarshal(result, &res)
	if err != nil {
		log.Printf("@[%s] json.Unmarshal failed, err=%v, params=%s, data=%+v, header=%+v", fn, err, string(result), params, headers)
		return
	}

	// 请求成功
	if res.Code == "0000" {
		url = res.Data.OrderDownloadURL
		return
	}
	err = fmt.Errorf("Err:%s", res.Message)
	return
}

// DownloadBillFile 下载日流水文件
func (c *Client) DownloadBillFile(billDate string) (url string, err error) {
	fn := "DownloadBillFile"
	data := map[string]string{
		"bill_date": billDate,
	}

	params, err := BuildParams(data, c.Appkey, c.Des3Key, c.PrivateKey)
	if err != nil {
		log.Printf("@[%s] BuildParam failed, err=%v, data=%+v", fn, err, data)
		return
	}

	headers := BuildHeader(c.DealerID)
	log.Printf("@[%s] begin, params=%+v, header=%v", fn, params, headers)

	result, err := Get(c.Gateway+DownloadBillURL, params, headers)
	if err != nil {
		log.Printf("@[%s] Get failed, err=%v, params=%+v, headers=%+v", fn, err, params, headers)
		return
	}
	log.Printf("@[%s] end, data=%+v, header=%v, result=%s", fn, params, headers, string(result))

	var res DownloadBillFileResponse
	err = json.Unmarshal(result, &res)
	if err != nil {
		log.Printf("@[%s] json.Unmarshal failed, err=%v, params=%s, data=%+v, header=%+v", fn, err, string(result), params, headers)
		return
	}

	// 请求成功
	if res.Code == "0000" {
		url = res.Data.BillDownloadURL
		return
	}
	err = fmt.Errorf("Err:%s", res.Message)
	return
}

// QueryDayOrdersData 查询日订单数据
func (c *Client) QueryDayOrdersData(queryDayOrdersParam *QueryDayOrdersParam) (dayOrderInfos []DayOrderInfo, err error) {
	fn := "QueryDayOrdersData"
	data := map[string]interface{}{
		"order_date": queryDayOrdersParam.OrderDate,
		"offset":     queryDayOrdersParam.Offset,
		"length":     queryDayOrdersParam.Length,
		"channel":    queryDayOrdersParam.Channel,
	}
	params, err := BuildParams(data, c.Appkey, c.Des3Key, c.PrivateKey)
	if err != nil {
		log.Printf("@[%s] BuildParam failed, err=%v, data=%+v", fn, err, data)
		return
	}

	headers := BuildHeader(c.DealerID)
	log.Printf("@[%s] begin, params=%+v, header=%v", fn, params, headers)

	result, err := Get(c.Gateway+QueryDayOrdersDataUrl, params, headers)
	if err != nil {
		log.Printf("@[%s] Get failed, err=%v, params=%+v, headers=%+v", fn, err, params, headers)
		return
	}
	log.Printf("@[%s] end, data=%+v, header=%v, result=%s", fn, params, headers, string(result))

	var res QueryDayOrdersResponse
	err = json.Unmarshal(result, &res)
	if err != nil {
		log.Printf("@[%s] json.Unmarshal failed, err=%v, params=%s, data=%+v, header=%+v", fn, err, string(result), params, headers)
		return
	}

	// 请求成功
	if res.Code == "0000" {
		dayOrderInfos = res.Data.DayOrderInfos
		return
	}
	err = fmt.Errorf("Err:%s", res.Message)
	return
}

// QueryDayOrdersFile 查询日订单文件（支付和退款订单）
func (c *Client) QueryDayOrdersFile(orderDate string) (url string, err error) {
	fn := "QueryDayOrdersFile"
	data := map[string]string{
		"order_date": orderDate,
	}
	params, err := BuildParams(data, c.Appkey, c.Des3Key, c.PrivateKey)
	if err != nil {
		log.Printf("@[%s] BuildParam failed, err=%v, data=%+v", fn, err, data)
		return
	}

	headers := BuildHeader(c.DealerID)
	log.Printf("@[%s] begin, params=%+v, header=%v", fn, params, headers)

	result, err := Get(c.Gateway+QueryDayOrdersFileUrl, params, headers)
	if err != nil {
		log.Printf("@[%s] Get failed, err=%v, params=%+v, headers=%+v", fn, err, params, headers)
		return
	}
	log.Printf("@[%s] end, data=%+v, header=%v, result=%s", fn, params, headers, string(result))

	var res QueryDayOrdersFileResponse
	err = json.Unmarshal(result, &res)
	if err != nil {
		log.Printf("@[%s] json.Unmarshal failed, err=%v, params=%s, data=%+v, header=%+v", fn, err, string(result), params, headers)
		return
	}

	// 请求成功
	if res.Code == "0000" {
		url = res.Data.Url
		return
	}
	err = fmt.Errorf("Err:%s", res.Message)
	return
}

// QueryDayBills 查询日流水数据
func (c *Client) QueryDayBills(queryDayBillsParam *QueryDayBillsParam) (dayBillInfos []DayBillInfo, err error) {
	fn := "QueryDayBills"
	data := map[string]interface{}{
		"bill_date": queryDayBillsParam.BillDate,
		"offset":    queryDayBillsParam.Offset,
		"length":    queryDayBillsParam.Length,
	}
	params, err := BuildParams(data, c.Appkey, c.Des3Key, c.PrivateKey)
	if err != nil {
		log.Printf("@[%s] BuildParam failed, err=%v, data=%+v", fn, err, data)
		return
	}

	headers := BuildHeader(c.DealerID)
	log.Printf("@[%s] begin, params=%+v, header=%v", fn, params, headers)

	result, err := Get(c.Gateway+QueryDayBillUrl, params, headers)
	if err != nil {
		log.Printf("@[%s] Get failed, err=%v, params=%+v, headers=%+v", fn, err, params, headers)
		return
	}
	log.Printf("@[%s] end, data=%+v, header=%v, result=%s", fn, params, headers, string(result))

	var res QueryDayBillsResponse
	err = json.Unmarshal(result, &res)
	if err != nil {
		log.Printf("@[%s] json.Unmarshal failed, err=%v, params=%s, data=%+v, header=%+v", fn, err, string(result), params, headers)
		return
	}

	// 请求成功
	if res.Code == "0000" {
		dayBillInfos = res.Data.DayBillInfos
		return
	}
	err = fmt.Errorf("Err:%s", res.Message)
	return
}

// QueryDayBills 查询日流水数据
func (c *Client) QueryDailyStatements(statementDate string) (dailyStatements []DailyStatement, err error) {
	fn := "QueryDailyStatements"
	data := map[string]string{
		"statement_date": statementDate,
	}
	params, err := BuildParams(data, c.Appkey, c.Des3Key, c.PrivateKey)
	if err != nil {
		log.Printf("@[%s] BuildParam failed, err=%v, data=%+v", fn, err, data)
		return
	}

	headers := BuildHeader(c.DealerID)
	log.Printf("@[%s] begin, params=%+v, header=%v", fn, params, headers)

	result, err := Get(c.Gateway+QueryDailyStatementsUrl, params, headers)
	if err != nil {
		log.Printf("@[%s] Get failed, err=%v, params=%+v, headers=%+v", fn, err, params, headers)
		return
	}
	log.Printf("@[%s] end, data=%+v, header=%v, result=%s", fn, params, headers, string(result))

	var res QueryDailyStatementsResponse
	err = json.Unmarshal(result, &res)
	if err != nil {
		log.Printf("@[%s] json.Unmarshal failed, err=%v, params=%s, data=%+v, header=%+v", fn, err, string(result), params, headers)
		return
	}

	// 请求成功
	if res.Code == "0000" {
		dailyStatements = res.Data.DailyStatements
		return
	}
	err = fmt.Errorf("Err:%s", res.Message)
	return
}
