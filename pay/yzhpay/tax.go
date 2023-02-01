package yzhpay

import (
	"crypto/rand"
	"crypto/rsa"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
)

// TaxfileDowload 下载个税扣缴明细表
func (c *Client) TaxfileDowload(entId, yearMonth string) (fileInfo []TaxfileDowloadInfo, err error) {
	fn := "TaxfileDowload"
	data := map[string]interface{}{
		"dealer_id":  c.DealerID, // 商户ID(必填)
		"ent_id":     entId,      // 商户签约主体，其中天津：accumulus_tj， 上海：accumulus_sh(必填)
		"year_month": yearMonth,  // 所属期（必填）, 注意格式 yyyy-mm
	}

	params, err := BuildParams(data, c.Appkey, c.Des3Key, c.PrivateKey)
	if err != nil {
		log.Printf("@[%s] BuildParam failed, err=%v, data=%+v", fn, err, data)
		return
	}

	headers := BuildHeader(c.DealerID)
	log.Printf("@[%s] begin, params=%+v, header=%v", fn, params, headers)

	result, err := Post(c.Gateway+TaxFileDownloadUrl, params, headers)
	if err != nil {
		log.Printf("@[%s] Post failed, err=%v, params=%+v, headers=%+v", fn, err, params, headers)
		return
	}
	log.Printf("@[%s] end, data=%+v, header=%v, result=%s", fn, params, headers, string(result))

	var res TaxfileDowloadResonse
	err = json.Unmarshal(result, &res)
	if err != nil {
		log.Printf("@[%s] json.Unmarshal failed, err=%v, params=%s, data=%+v, header=%+v", fn, err, string(result), params, headers)
		return
	}

	// 请求成功
	if res.Code == "0000" {
		fileInfo = res.Data.FileInfo
		return
	}
	err = fmt.Errorf("Err:%s", res.Message)
	return
}

// TaxUserCross 查询纳税人是否为跨集团用户
func (c *Client) TaxUserCross(year, idCard, entId string) (IsCross bool, err error) {
	fn := "TaxfileDowload"
	data := map[string]interface{}{
		"dealer_id": c.DealerID, // 商户ID(必填)
		"year":      year,       // 用户报税所在年份
		"id_card":   idCard,     //所查询用户的身份证件号码
		"ent_id":    entId,      // 商户签约主体，其中天津：accumulus_tj， 上海：accumulus_sh(必填)
	}

	params, err := BuildParams(data, c.Appkey, c.Des3Key, c.PrivateKey)
	if err != nil {
		log.Printf("@[%s] BuildParam failed, err=%v, data=%+v", fn, err, data)
		return
	}

	headers := BuildHeader(c.DealerID)
	log.Printf("@[%s] begin, params=%+v, header=%v", fn, params, headers)

	result, err := Post(c.Gateway+QueryTaxUserCrossUrl, params, headers)
	if err != nil {
		log.Printf("@[%s] Post failed, err=%v, params=%+v, headers=%+v", fn, err, params, headers)
		return
	}
	log.Printf("@[%s] end, data=%+v, header=%v, result=%s", fn, params, headers, string(result))

	var res QueryTaxUserCrossResponse
	err = json.Unmarshal(result, &res)
	if err != nil {
		log.Printf("@[%s] json.Unmarshal failed, err=%v, params=%s, data=%+v, header=%+v", fn, err, string(result), params, headers)
		return
	}

	// 请求成功
	if res.Code == "0000" {
		IsCross = res.Data.IsCross
		return
	}
	err = fmt.Errorf("Err:%s", res.Message)
	return
}

/**
 *
 * RSA公钥加密
 * 测试用
 */
func (c *Client) PwdEncryption(ciphertext string) (encryptInfo string, err error) {
	fn := "PwdEncryption"
	pubKey, err := LoadPublicKey([]byte(c.PublicKey))
	if err != nil {
		return
	}
	//对字符串进行加密
	encryptText, err := rsa.EncryptPKCS1v15(rand.Reader, pubKey, []byte(ciphertext))
	if err != nil {
		log.Printf("@[%s] PwdEncryption failed, err=%v, pubKey=%+v", fn, err, pubKey)
		return
	}
	//返回密文，base64转码
	encryptInfo = base64.StdEncoding.EncodeToString(encryptText)
	return
}

/**
 * 个税文件解压缩密码解密:
 * ①使用 BASE64 算法对数据解码，得到密文。
 * ②使用 RSA 算法用私钥对密文进行解密，得到文件解压缩密码。
 */
func (c *Client) PwdDecryption(ciphertext string) (plainText string, err error) {
	fn := "PwdDecryption"

	priKey, err := LoadPrivateKey([]byte(c.PrivateKey))
	if err != nil {
		return
	}
	cipherBytes, err := base64.StdEncoding.DecodeString(ciphertext)
	//对密文进行解密
	plainByte, err := rsa.DecryptPKCS1v15(rand.Reader, priKey, cipherBytes)

	if err != nil {
		log.Printf("@[%s] PwdDecryption failed, err=%v", fn, err)
		return
	}
	plainText = string(plainByte)
	//返回明文
	return
}
