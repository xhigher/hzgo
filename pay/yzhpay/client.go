package yzhpay

// Client 客户端
type Client struct {
	BrokerID     string // 代征主体ID
	DealerID     string // 商户ID
	Gateway      string // 路由
	Appkey       string // 商户appkey
	Des3Key      string // 商户des3key
	PrivateKey   string // 商户秘钥
	PublicKey    string //商户公钥
	YunPublicKey string // 云账户公钥
}

// 接口配置信息
const (
	productionGateway = "https://api-service.yunzhanghu.com"
	sandboxGateway    = "https://api-service.yunzhanghu.com/sandbox"
)

// New 新建客户端
func NewClient(brokerID, dealerID, appKey, des3Key, privateKey, publicKey, yunPublicKey string) *Client {
	gateway := productionGateway
	//if !env.IsProd() {
	//	gateway = sandboxGateway
	//}
	return &Client{
		BrokerID:     brokerID,
		DealerID:     dealerID,
		Gateway:      gateway,
		Appkey:       appKey,
		Des3Key:      des3Key,
		PrivateKey:   privateKey,
		PublicKey:    publicKey,
		YunPublicKey: yunPublicKey,
	}
}
