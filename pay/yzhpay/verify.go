package yzhpay

import (
	"encoding/json"
	"fmt"
	"log"
)

/**
 * UploadUserInfo 上传用户免验证名单信息
 *
 * 此接⼝可以修改未审核的信息。当发现上传信息有误时，可⽤正确的信息再次进⾏请求，
 * 若旧信息未审核，则直接覆盖；若审核已被拒绝，则写⼊新记录；若审核已通过，则不允许修改，请联系商务⼈员进⾏⼿⼯修改。修改操作需使⽤新的流⽔号。
 *
 * card_type证件类型码:
 *  passport 护照
 *  mtphkm 港澳居民来往内地通行证
 *  mtpt 台湾居民来往大陆通行证（台胞证）
 *  rphkm 中华人民共和国港澳居民居住证
 *  rpt 中华人民共和国台湾居民居住证
 *  fpr 外国人永久居留身份证
 *  ffwp 中华人民共和国外国人就业许可证书
 */
func (c *Client) UploadUserInfo(param *UserInfoParam) (ok bool, err error) {
	fn := "UploadUserInfo"
	data := map[string]interface{}{
		"broker_id":     c.BrokerID,         // 综合服务主体(必填)
		"dealer_id":     c.DealerID,         // 商户ID(必填)
		"ref":           param.Ref,          // 请求流水号，唯一，回调时会返回  (必填)
		"id_card":       param.IDCard,       // 护照、港澳台居⺠居住证等证件号(必填)
		"real_name":     param.RealName,     // 姓名(必填)
		"card_type":     param.CardType,     // 证件类型码: (必填)
		"country":       param.Country,      // 国别（地区）代码，见文档附录  (必填)
		"birthday":      param.Birthday,     // 出生日期 (必填)
		"gender":        param.Gender,       // 性别 (必填)
		"user_images":   param.UserImages,   // 人员信息图片需要base64编码，支持jpg、png、jpeg格式，正反面各一张；护照类型存在只有一面的情况，可以上传一张。图片大小不超过20M。
		"comment_apply": param.CommentApply, // 申请备注  (必填)
		"notify_url":    param.NotifyURL,    // 回调地址 (必填)
	}

	params, err := BuildParams(data, c.Appkey, c.Des3Key, c.PrivateKey)
	if err != nil {
		log.Printf("@[%s] BuildParam failed, err=%v, data=%+v", fn, err, data)
		return
	}

	headers := BuildHeader(c.DealerID)
	log.Printf("@[%s] begin, params=%+v, header=%v", fn, params, headers)

	result, err := Post(c.Gateway+UploadUserURL, params, headers)
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

// CheckUserExist 校验免验证用户是否存在
func (c *Client) CheckUserExist(idCard, realName string) (ok bool, err error) {
	fn := "CheckUserExist"
	data := map[string]string{
		"id_card":   idCard,
		"real_name": realName,
	}

	params, err := BuildParams(data, c.Appkey, c.Des3Key, c.PrivateKey)
	if err != nil {
		log.Printf("@[%s] BuildParam failed, err=%v, data=%+v", fn, err, data)
		return
	}

	headers := BuildHeader(c.DealerID)
	log.Printf("@[%s] begin, params=%+v, header=%v", fn, params, headers)

	result, err := Post(c.Gateway+CheckExistUserURL, params, headers)
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

// ElementVerifyRequest 银行卡四要素请求鉴权
func (c *Client) ElementVerifyRequest(idCard, realName, cardNo, mobile string) (ref string, err error) {
	fn := "ElementVerifyRequest"
	data := map[string]string{
		"id_card":   idCard,
		"real_name": realName,
		"card_no":   cardNo,
		"mobile":    mobile,
	}

	params, err := BuildParams(data, c.Appkey, c.Des3Key, c.PrivateKey)
	if err != nil {
		log.Printf("@[%s] BuildParam failed, err=%v, data=%+v", fn, err, data)
		return
	}

	headers := BuildHeader(c.DealerID)
	log.Printf("@[%s] begin, params=%+v, header=%v", fn, params, headers)

	result, err := Post(c.Gateway+Element4RequestURL, params, headers)
	if err != nil {
		log.Printf("@[%s] Post failed, err=%v, params=%+v, headers=%+v", fn, err, params, headers)
		return
	}
	log.Printf("@[%s] end, data=%+v, header=%v, result=%s", fn, params, headers, string(result))

	var res ElementVerifyResponse
	err = json.Unmarshal(result, &res)
	if err != nil {
		log.Printf("@[%s] json.Unmarshal failed, err=%v, params=%s, data=%+v, header=%+v", fn, err, string(result), params, headers)
		return
	}

	// 请求成功
	if res.Code == "0000" {
		ref = res.Data.Ref
		return
	}
	err = fmt.Errorf("Err:%s", res.Message)
	return
}

// ElementVerifyConfirm 银行卡四要素确认鉴权
func (c *Client) ElementVerifyConfirm(idCard, realName, cardNo, mobile, ref, captcha string) (ok bool, err error) {
	fn := "ElementVerifyConfirm"
	data := map[string]string{
		"id_card":   idCard,
		"real_name": realName,
		"card_no":   cardNo,
		"mobile":    mobile,
		"ref":       ref,
		"captcha":   captcha,
	}

	params, err := BuildParams(data, c.Appkey, c.Des3Key, c.PrivateKey)
	if err != nil {
		log.Printf("@[%s] BuildParam failed, err=%v, data=%+v", fn, err, data)
		return
	}

	headers := BuildHeader(c.DealerID)
	log.Printf("@[%s] begin, params=%+v, header=%v", fn, params, headers)

	result, err := Post(c.Gateway+Element4ConfirmURL, params, headers)
	if err != nil {
		log.Printf("@[%s] Post failed, err=%v, params=%+v, headers=%+v", fn, err, params, headers)
		return
	}
	log.Printf("@[%s] end, data=%+v, header=%v, result=%s", fn, params, headers, string(result))

	var res BaseResponse
	err = json.Unmarshal(result, &res)
	if err != nil {
		log.Printf("@[%s] json.Unmarshal failed, err=%v, params=%s, data=%+v, header=%+v", fn, err, string(result), params, headers)
		return
	}

	// 请求成功
	if res.Code == "0000" {
		ok = true
		return
	}
	err = fmt.Errorf("Err:%s", res.Message)
	return
}

// Element4Check 银行卡四要素鉴权
func (c *Client) Element4Check(idCard, realName, cardNo, mobile string) (ok bool, err error) {
	fn := "Element4Check"
	data := map[string]string{
		"id_card":   idCard,
		"real_name": realName,
		"card_no":   cardNo,
		"mobile":    mobile,
	}

	params, err := BuildParams(data, c.Appkey, c.Des3Key, c.PrivateKey)
	if err != nil {
		log.Printf("@[%s] BuildParam failed, err=%v, data=%+v", fn, err, data)
		return
	}

	headers := BuildHeader(c.DealerID)
	log.Printf("@[%s] begin, params=%+v, header=%v", fn, params, headers)

	result, err := Post(c.Gateway+Element4URL, params, headers)
	if err != nil {
		log.Printf("@[%s] Post failed, err=%v, params=%+v, headers=%+v", fn, err, params, headers)
		return
	}
	log.Printf("@[%s] end, data=%+v, header=%v, result=%s", fn, params, headers, string(result))

	var res BaseResponse
	err = json.Unmarshal(result, &res)
	if err != nil {
		log.Printf("@[%s] json.Unmarshal failed, err=%v, params=%s, data=%+v, header=%+v", fn, err, string(result), params, headers)
		return
	}

	// 请求成功
	if res.Code == "0000" {
		ok = true
		return
	}
	err = fmt.Errorf("Err:%s", res.Message)
	return
}

// Element3Check 银行卡三要素鉴权
func (c *Client) Element3Check(idCard, realName, cardNo string) (ok bool, err error) {
	fn := "Element3Check"
	data := map[string]string{
		"id_card":   idCard,
		"real_name": realName,
		"card_no":   cardNo,
	}

	params, err := BuildParams(data, c.Appkey, c.Des3Key, c.PrivateKey)
	if err != nil {
		log.Printf("@[%s] BuildParam failed, err=%v, data=%+v", fn, err, data)
		return
	}

	headers := BuildHeader(c.DealerID)
	log.Printf("@[%s] begin, params=%+v, header=%v", fn, params, headers)

	result, err := Post(c.Gateway+Element3URL, params, headers)
	if err != nil {
		log.Printf("@[%s] Post failed, err=%v, params=%+v, headers=%+v", fn, err, params, headers)
		return
	}
	log.Printf("@[%s] end, data=%+v, header=%v, result=%s", fn, params, headers, string(result))

	var res BaseResponse
	err = json.Unmarshal(result, &res)
	if err != nil {
		log.Printf("@[%s] json.Unmarshal failed, err=%v, params=%s, data=%+v, header=%+v", fn, err, string(result), params, headers)
		return
	}

	// 请求成功
	if res.Code == "0000" {
		ok = true
		return
	}
	err = fmt.Errorf("Err:%s", res.Message)
	return
}

// IDCheck 实名制二要素鉴权接口
func (c *Client) IDCheck(idCard, realName string) (ok bool, err error) {
	fn := "IDCheck"
	data := map[string]string{
		"id_card":   idCard,
		"real_name": realName,
	}

	params, err := BuildParams(data, c.Appkey, c.Des3Key, c.PrivateKey)
	if err != nil {
		log.Printf("@[%s] BuildParam failed, err=%v, data=%+v", fn, err, data)
		return
	}

	headers := BuildHeader(c.DealerID)
	log.Printf("@[%s] begin, params=%+v, header=%v", fn, params, headers)

	result, err := Post(c.Gateway+IDCheckURL, params, headers)
	if err != nil {
		log.Printf("@[%s] Post failed, err=%v, url=%+v, params=%+v, headers=%+v", c.Gateway+IDCheckURL, fn, err, params, headers)
		return
	}
	log.Printf("@[%s] end, data=%+v, header=%v, result=%s", fn, params, headers, string(result))

	var res BaseResponse
	err = json.Unmarshal(result, &res)
	if err != nil {
		log.Printf("@[%s] json.Unmarshal failed, err=%v, params=%s, data=%+v, header=%+v", fn, err, string(result), params, headers)
		return
	}

	// 请求成功
	if res.Code == "0000" {
		ok = true
		return
	}
	err = fmt.Errorf("Err:%s", res.Message)
	return
}

// QueryBankCardInfo 查询银行卡信息
func (c *Client) QueryBankCardInfo(cardNo, bankName string) (cardInfo BankCardInfo, err error) {
	fn := "QueryBankCardInfo"
	data := map[string]string{
		"card_no":   cardNo,
		"bank_name": bankName,
	}

	params, err := BuildParams(data, c.Appkey, c.Des3Key, c.PrivateKey)
	if err != nil {
		log.Printf("@[%s] BuildParam failed, err=%v, data=%+v", fn, err, data)
		return
	}

	headers := BuildHeader(c.DealerID)
	log.Printf("@[%s] begin, params=%+v, header=%v", fn, params, headers)

	result, err := Get(c.Gateway+BankCardInfoURL, params, headers)
	if err != nil {
		log.Printf("@[%s] Get failed, err=%v, params=%+v, headers=%+v", fn, err, params, headers)
		return
	}
	log.Printf("@[%s] end, data=%+v, header=%v, result=%s", fn, params, headers, string(result))

	var res QueryBankCardResponse
	err = json.Unmarshal(result, &res)
	if err != nil {
		log.Printf("@[%s] json.Unmarshal failed, err=%v, params=%s, data=%+v, header=%+v", fn, err, string(result), params, headers)
		return
	}

	// 请求成功
	if res.Code == "0000" {
		cardInfo = res.Data
		return
	}
	err = fmt.Errorf("Err:%s", res.Message)
	return
}
