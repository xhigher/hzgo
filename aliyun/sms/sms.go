package sms

import (
	"encoding/json"
	"errors"
	"fmt"
	openapi "github.com/alibabacloud-go/darabonba-openapi/client"
	dysmsapi "github.com/alibabacloud-go/dysmsapi-20170525/v2/client"
	"github.com/alibabacloud-go/tea/tea"
)

type SmsConfig struct {
	Endpoint        string
	AccessKeyId     string
	AccessKeySecret string
	SignName        string
}

var (
	smsConfig *SmsConfig
	smsClient *dysmsapi.Client
)

func Init(conf *SmsConfig) (err error) {
	if conf == nil {
		err = errors.New("error config nil")
		return
	}
	if len(conf.Endpoint) == 0 {
		err = errors.New("error Endpoint nil")
		return
	}
	if len(conf.AccessKeyId) == 0 {
		err = errors.New("error AccessKeyId nil")
		return
	}
	if len(conf.AccessKeySecret) == 0 {
		err = errors.New("error AccessKeySecret nil")
		return
	}
	smsConfig = conf
	config := &openapi.Config{
		// 您的AccessKey ID。
		AccessKeyId: tea.String(smsConfig.AccessKeyId),
		// 您的AccessKey Secret。
		AccessKeySecret: tea.String(smsConfig.AccessKeySecret),
		Endpoint:        tea.String(smsConfig.Endpoint),
		// 设置HTTP代理。
		// HttpProxy: tea.String("http://xx.xx.xx.xx:xxxx"),
		// 设置HTTPS代理。
		// HttpsProxy: tea.String("https://xx.xx.xx.xx:xxxx"),
	}
	smsClient, err = dysmsapi.NewClient(config)
	if err != nil {
		err = errors.New("error client new failed")
		return
	}

	return
}

func SendSmsCode(tplCode, phoneno, code string) (err error) {
	data := make(map[string]interface{})
	data["code"] = code
	dataBytes, err := json.Marshal(data)
	if err != nil {
		return
	}
	request := &dysmsapi.SendSmsRequest{
		PhoneNumbers:  tea.String(phoneno),
		SignName:      tea.String(smsConfig.SignName),
		TemplateCode:  tea.String(tplCode),
		TemplateParam: tea.String(string(dataBytes)),
	}
	_, err = smsClient.SendSms(request)
	if err != nil {
		return
	}
	return
}

type SmsAdvertTarget struct {
	Phoneno string
	Data    map[string]interface{}
}

func SendSmsAdvert(tpl string, targets []SmsAdvertTarget) (err error) {
	total := len(targets)
	phonenos := make([]string, total)
	allData := make([]map[string]interface{}, total)
	for i, item := range targets {
		phonenos[i] = item.Phoneno
		allData[i] = item.Data
	}
	phonenosBytes, err := json.Marshal(phonenos)
	if err != nil {
		return
	}
	allDataBytes, err := json.Marshal(allData)
	if err != nil {
		return
	}
	request := &dysmsapi.SendBatchSmsRequest{
		PhoneNumberJson:   tea.String(string(phonenosBytes)),
		SignNameJson:      tea.String(smsConfig.SignName),
		TemplateCode:      tea.String(tpl),
		TemplateParamJson: tea.String(string(allDataBytes)),
	}
	_, err = smsClient.SendBatchSms(request)
	if err != nil {
		return
	}
	return
}

func CreateShortUrl(url string, days int) (shortUrl string, err error) {
	request := &dysmsapi.AddShortUrlRequest{
		SourceUrl:     tea.String(url),
		ShortUrlName:  tea.String(""),
		EffectiveDays: tea.String(fmt.Sprintf("%d", days)),
	}
	_, err = smsClient.AddShortUrl(request)
	if err != nil {
		return
	}
	return
}
