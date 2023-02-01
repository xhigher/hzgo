package iappay

import (
	"context"
	"github.com/go-pay/gopay/apple"
	"github.com/go-pay/gopay/pkg/xlog"
	"github.com/xhigher/hzgo/config"
)

type TradeData struct {
	TradeNo string `json:"trade_no"`
}

type IapPayClient struct {
	password string
	bundleId string
}

var iappayClient *IapPayClient

func Init(conf *config.IapPayConfig) {
	iappayClient = newClient(conf.Password)
}

func newClient(password string) *IapPayClient {
	return &IapPayClient{
		password: password,
	}
}

func VerifyOrder(receipt string, IsSandbox bool) (resp *apple.VerifyResponse, err error) {
	url := apple.UrlProd
	if IsSandbox {
		url = apple.UrlSandbox
	}
	resp, err = apple.VerifyReceipt(context.Background(), url, iappayClient.password, receipt)
	if err != nil {
		xlog.Error(err)
		return
	}
	if resp.Receipt != nil {
		xlog.Infof("VerifyOrder receipt:%+v", resp.Receipt)
	}

	return
}
