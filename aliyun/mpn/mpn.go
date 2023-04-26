package mpn

import (
	"errors"
	openapi "github.com/alibabacloud-go/darabonba-openapi/v2/client"
	dypnsapi "github.com/alibabacloud-go/dypnsapi-20170525/v2/client"
	teautils "github.com/alibabacloud-go/tea-utils/v2/service"
	"github.com/alibabacloud-go/tea/tea"
	"github.com/xhigher/hzgo/logger"
)

type MpnConfig struct {
	Endpoint        string
	AccessKeyId     string
	AccessKeySecret string
	SignName        string
}

var (
	mpnConfig *MpnConfig
	mpnClient *dypnsapi.Client
)

func Init(conf *MpnConfig) (err error) {
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
	mpnConfig = conf
	config := &openapi.Config{
		// 您的AccessKey ID。
		AccessKeyId: tea.String(mpnConfig.AccessKeyId),
		// 您的AccessKey Secret。
		AccessKeySecret: tea.String(mpnConfig.AccessKeySecret),
		Endpoint:        tea.String(mpnConfig.Endpoint),
		// 设置HTTP代理。
		// HttpProxy: tea.String("http://xx.xx.xx.xx:xxxx"),
		// 设置HTTPS代理。
		// HttpsProxy: tea.String("https://xx.xx.xx.xx:xxxx"),
	}
	mpnClient, err = dypnsapi.NewClient(config)
	if err != nil {
		err = errors.New("error client new failed")
		return
	}
	return
}

func GetAuthToken() (token string, err error) {
	request := &dypnsapi.GetAuthTokenRequest{}
	response, err := mpnClient.GetAuthToken(request)
	if err != nil {
		logger.Errorf("GetAuthToken error: %v", err)
		return
	}
	if response == nil {
		logger.Errorf("GetAuthToken response nil")
		return
	}

	respStr := teautils.ToJSONString(response)
	if *response.StatusCode != 200 {
		logger.Errorf("GetAuthToken response: %v", *respStr)
		return
	}
	token = *response.Body.TokenInfo.AccessToken
	return
}

func GetMobilePhone() (phone string, err error) {
	request := &dypnsapi.GetPhoneWithTokenRequest{}
	// 复制代码运行请自行打印 API 的返回值
	response, err := mpnClient.GetPhoneWithToken(request)
	if err != nil {
		logger.Errorf("GetPhoneWithToken error: %v", err)
		return
	}
	if response == nil {
		logger.Errorf("GetPhoneWithToken response nil")
		return
	}

	respStr := teautils.ToJSONString(response)
	if *response.StatusCode != 200 {
		logger.Errorf("GetPhoneWithToken response: %v", *respStr)
		return
	}
	phone = *response.Body.Data.Mobile
	return
}
