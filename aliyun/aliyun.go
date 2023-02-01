package aliyun

import (
	"errors"
	"github.com/xhigher/hzgo/aliyun/certify"
	"github.com/xhigher/hzgo/aliyun/oss"
	"github.com/xhigher/hzgo/aliyun/sms"
	"github.com/xhigher/hzgo/logger"
)

type AliyunConfig struct {
	Sms        *sms.SmsConfig
	Oss *oss.OssConfig
	Certify *certify.CertifyConfig
}

func Init(conf *AliyunConfig) (err error){
	if conf == nil {
		err = errors.New("error config nil")
		logger.Errorf("aliyun init error: %v", err)
		return
	}
	if conf.Oss != nil {
		err = oss.Init(conf.Oss)
		if err != nil {
			logger.Errorf("aliyun oss init error: %v", err)
			return
		}
	}
	if conf.Sms != nil {
		err = sms.Init(conf.Sms)
		if err != nil {
			logger.Errorf("aliyun sms init error: %v", err)
			return
		}
	}
	if conf.Certify != nil {
		err = certify.Init(conf.Certify)
		if err != nil {
			logger.Errorf("aliyun certify init error: %v", err)
			return
		}
	}
	return
}