package config

type PayConfig struct {
	AliPay *AliPayConfig
	WxPay  *WxPayConfig
	IapPay *IapPayConfig
}

type AliPayConfig struct {
	AppId         string
	AppPrivateKey string
	PublicKey     string
	IsProd        bool
	NotifyUrl     string
	ReturnUrl     string
}

type WxPayConfig struct {
	Id           int
	MchId        string
	SerialNo     string
	Apiv3Key     string
	PrivateKey   string
	NotifyUrl    string
	ComplaintUrl string
	SupportH5    bool
}

type IapPayConfig struct {
	Password string
}
