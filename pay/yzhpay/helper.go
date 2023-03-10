package yzhpay

import (
	"encoding/json"
	"fmt"
	"time"
)

// generateData 生成data
func generateData(v interface{}, des3key string) (string, error) {
	originData, err := json.Marshal(v)
	if err != nil {
		return "", err
	}
	return Encrypt(originData, des3key)
}

// generateMess 生成随机数
func generateMess() string {
	return fmt.Sprint(time.Now().Nanosecond())
}

// generateTimestamp 生成时间戳
func generateTimestamp() string {
	return fmt.Sprint(time.Now().Unix())
}

// generateRequestID 生成请求ID
func generateRequestID() string {
	return fmt.Sprint(time.Now().UnixNano())
}

// BuildParams 封装请求信息
func BuildParams(v interface{}, appKey, des3key, privateKey string) (map[string]string, error) {
	data, err := generateData(v, des3key)
	if err != nil {
		return nil, err
	}
	mess := generateMess()
	timestamp := generateTimestamp()
	plaintext := fmt.Sprintf("data=%s&mess=%s&timestamp=%s&key=%s", data, mess, timestamp, appKey)
	sign, err := Sign(plaintext, privateKey)
	if err != nil {
		return nil, err
	}
	return map[string]string{
		"data":      data,
		"mess":      mess,
		"timestamp": timestamp,
		"sign":      sign,
		"sign_type": "rsa",
	}, nil
}

// BuildHeader 封装请求头
func BuildHeader(dealerID string) map[string]string {
	return map[string]string{
		"dealer-id":  dealerID,
		"request-id": generateRequestID(),
	}
}
