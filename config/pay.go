package config

type PayConfig struct {
	AliPay *AliPayConfig
	WxPay  *WxPayConfig
	IapPay *IapPayConfig
}

type AliPayConfig struct {
	AppId         string `yaml:"app_id"`
	AppPrivateKey string `yaml:"app_private_key"`
	PublicKey     string `yaml:"public_key"`
	IsProd        bool   `yaml:"is_prod"`
	NotifyUrl     string `yaml:"notify_url"`
	ReturnUrl     string `yaml:"return_url"`
}

type WxPayConfig struct {
	Id int `yaml:"id"`
	MchId        string `yaml:"mch_id"`
	SerialNo     string `yaml:"serial_no"`
	Apiv3Key     string `yaml:"apiv3_key"`
	PrivateKey   string `yaml:"private_key"`
	NotifyUrl    string `yaml:"notify_url"`
	ComplaintUrl string `yaml:"complaint_url"`
	SupportH5    bool   `yaml:"support_h5"`
}

type IapPayConfig struct {
	Password string `yaml:"password"`
}
