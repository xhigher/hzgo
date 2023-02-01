package yzhpay

// OrderBaseInfo 订单基本信息
type OrderBaseInfo struct {
	OrderID   string // 商户订单号(必填, 保持唯一性,64个英文字符以内)
	RealName  string // 收款人姓名(必填)
	IDCard    string // 收款人手机号(必填)
	Pay       string // 打款金额(必填 单位:元)
	PayRemark string // 打款备注(选填, 最大20个字符,一个汉字占两个字符,不允许特殊字符)
	NotifyURL string // 回调地址(选填, 最大长度为200)
}

// BankOrderParam 银行卡订单信息
type BankOrderParam struct {
	OrderBaseInfo
	CardNo  string // 收款人银行卡号(必填)
	PhoneNo string // 收款人手机号(选填)
}

// AliOrderParam 支付宝订单信息
type AliOrderParam struct {
	OrderBaseInfo
	CardNo    string // 收款人支付宝号(必填)
	CheckName string // 校验支付宝账户姓名(固定值：Check)
}

// WxOrderParam 微信订单信息
type WxOrderParam struct {
	OrderBaseInfo
	OpenID    string // 商户AppID下用户OpenID(必填)
	WxAppID   string // 微信打款商户微信AppID(选填,最⼤长度为200)
	WxPayMode string // 微信打款模式(固定值：transfer)
}

// BaseResponse 基础响应信息
type BaseResponse struct {
	Code      string `json:"code"`       // 响应码
	Message   string `json:"message"`    // 响应信息
	RequestID string `json:"request_id"` // 请求ID
}

// BaseCheckResponse 基础校验响应信息
type BaseCheckResponse struct {
	BaseResponse
	Data struct {
		Ok bool `json:"ok"` // 是否成功
	} `json:"data"`
}

// CreateOrderResponse 创建订单响应信息
type CreateOrderResponse struct {
	BaseResponse
	Data struct {
		Pay     string `json:"pay"`      // 打款金额
		Ref     string `json:"ref"`      // 综合服务平台订单流水号
		OrderID string `json:"order_id"` // 商户订单流水号
	} `json:"data"`
}

// QueryOrderResponse 查询订单响应信息
type QueryOrderResponse struct {
	BaseResponse
	Data OrderInfo `json:"data"`
}

// OrderInfo 订单详细信息
type OrderInfo struct {
	BrokerAmount        string `json:"broker_amount"`
	BrokerBankBill      string `json:"broker_bank_bill"`
	BrokerFee           string `json:"broker_fee"`
	BrokerID            string `json:"broker_id"`
	CardNo              string `json:"card_no"`
	CreatedAt           string `json:"created_at"`
	DealerID            string `json:"dealer_id"`
	EncryData           string `json:"encry_data"`
	FeeAmount           string `json:"fee_amount"`
	FinishedTime        string `json:"finished_time"`
	IDCard              string `json:"id_card"`
	OrderID             string `json:"order_id"`
	Pay                 string `json:"pay"`
	PayRemark           string `json:"pay_remark"`
	PhoneNo             string `json:"phone_no"`
	RealName            string `json:"real_name"`
	Ref                 string `json:"ref"`
	Status              string `json:"status"`
	StatusDetail        string `json:"status_detail"`
	StatusDetailMessage string `json:"status_detail_message"`
	StatusMessage       string `json:"status_message"`
	UserFee             string `json:"user_fee"`          //用户服务费
	WithdrawPlatform    string `json:"withdraw_platform"` //结算渠道  bankpay：银行卡 alipay：支付宝 wxpay：微信
	BrokerRealFee       string `json:"broker_real_fee"`   //余额账户支出服务费
	BrokerDeductFee     string `json:"broker_deduct_fee"` //抵扣账户支出服务费
	BankName            string `json:"bank_name"`         //银行名称
}

// OrderDetailInfo 回调通知订单信息
type OrderDetailInfo struct {
	BrokerAmount        string `json:"broker_amount"`
	BrokerBankBill      string `json:"broker_bank_bill"` //支付交易流水号
	BrokerFee           string `json:"broker_fee"`
	BrokerID            string `json:"broker_id"`
	CardNo              string `json:"card_no"`
	CreatedAt           string `json:"created_at"`
	DealerID            string `json:"dealer_id"`
	FinishedTime        string `json:"finished_time"`
	IDCard              string `json:"id_card"`
	OrderID             string `json:"order_id"`
	Pay                 string `json:"pay"`
	PayRemark           string `json:"pay_remark"`
	PhoneNo             string `json:"phone_no"`
	RealName            string `json:"real_name"`
	Ref                 string `json:"ref"` //综合服务平台流水号，唯一
	Status              string `json:"status"`
	StatusDetail        string `json:"status_detail"`
	StatusDetailMessage string `json:"status_detail_message"`
	StatusMessage       string `json:"status_message"`
	UserFee             string `json:"user_fee"`          //用户服务费
	WithdrawPlatform    string `json:"withdraw_platform"` //结算渠道  bankpay：银行卡 alipay：支付宝 wxpay：微信
	BrokerRealFee       string `json:"broker_real_fee"`   //余额账户支出服务费
	BrokerDeductFee     string `json:"broker_deduct_fee"` //抵扣账户支出服务费
	BankName            string `json:"bank_name"`         //银行名称
}

// OrderCallBackResponse 订单回调信息
type OrderCallBackResponse struct {
	NotifyID   string          `json:"notify_id"`
	NotifyTime string          `json:"notify_time"`
	Data       OrderDetailInfo `json:"data"`
}

// QueryAccountBalanceResponse 查询商户余额响应信息
type QueryAccountBalanceResponse struct {
	BaseResponse
	Data struct {
		DealerInfos []AccountBalance `json:"dealer_infos"`
	} `json:"data"`
}

// QueryVaAccountResponse 查询商户Va账户响应信息
type QueryVaAccountResponse struct {
	BaseResponse
	Data VaAccount `json:"data"`
}

// VaAccount 商户Va账户响应信息
type VaAccount struct {
	AcctName       string `json:"acct_name"`        //账户名称
	AcctNo         string `json:"acct_no"`          //专属账号
	BankName       string `json:"bank_name"`        //银行名称
	DealerAcctName string `json:"dealer_acct_name"` //付款账户
}

// AccountBalance 账户余额信息
type AccountBalance struct {
	BrokerID         string `json:"broker_id"`          // 代征主体ID
	BankCardBalance  string `json:"bank_card_balance"`  // 银行卡余额
	AlipayBalance    string `json:"alipay_balance"`     // ⽀付宝余额
	WxpayBalance     string `json:"wxpay_balance"`      // 微信余额
	IsBankCard       bool   `json:"is_bank_card"`       // 是否开通银行卡通道
	IsAlipay         bool   `json:"is_alipay"`          // 是否开通付宝通道
	IsWxpay          bool   `json:"is_wxpay"`           // 是否开通微信通道
	RebateFeeBalance string `json:"rebate_fee_balance"` // 服务费返点余额
	AcctBalance      string `json:"acct_balance"`       // 余额账户余额
	TotalBalance     string `json:"total_balance"`      // 总余额
}

// QueryReceiptFileResponse 查询电子回单响应信息
type QueryReceiptFileResponse struct {
	BaseResponse
	Data OrderReceiptFile `json:"data"`
}

// OrderReceiptFile 电子回单信息
type OrderReceiptFile struct {
	ExpireTime string `json:"expire_time"` // 过期时间
	FileName   string `json:"file_name"`   // 文件名称
	URL        string `json:"url"`         // 下载地址
}

// DownloadOrderFileResponse 下载日订单响应信息
type DownloadOrderFileResponse struct {
	BaseResponse
	Data struct {
		OrderDownloadURL string `json:"order_download_url"` // url地址
	} `json:"data"`
}

// DownloadBillFileResponse 下载日流水响应信息
type DownloadBillFileResponse struct {
	BaseResponse
	Data struct {
		BillDownloadURL string `json:"bill_download_url"` // url地址
	} `json:"data"`
}

// QueryRechargeRecordResponse 充值记录响应信息
type QueryRechargeRecordResponse struct {
	BaseResponse
	Data []RechargeRecord `json:"data"`
}

// RechargeRecord 充值记录信息
type RechargeRecord struct {
	BrokerID          string `json:"broker_id"`           // 代征主体ID
	DealerID          string `json:"dealer_id"`           // 商户ID
	ActualAmount      int    `json:"actual_amount"`       // 实际到账金额
	Amount            int    `json:"amount"`              // 充值金额
	CreatedAt         string `json:"created_at"`          // 创建时间
	RechargeChannel   string `json:"recharge_channel"`    // 充值渠道
	RechargeID        string `json:"recharge_id"`         // 充值记录ID
	Remark            string `json:"remark"`              //备注
	RechargeAccountNo string `json:"recharge_account_no"` // 付款银行账号
}

// UserInfoParam 免验证用户名单信息
type UserInfoParam struct {
	RealName     string   `json:"real_name"`     // 姓名
	IDCard       string   `json:"id_card"`       // 证件号
	Birthday     string   `json:"birthday"`      // 出生日期
	CardType     string   `json:"card_type"`     // 证件类型
	Country      string   `json:"country"`       // 国别（地区）代码
	Gender       string   `json:"gender"`        // 性别
	NotifyURL    string   `json:"notify_url"`    // 回调地址
	Ref          string   `json:"ref"`           // 流水号(回调时附带)
	UserImages   []string `json:"user_images"`   // 证件照片
	CommentApply string   `json:"comment_apply"` // 申请备注
}

// UserCallBackInfo 通知用户上传信息
type UserCallBackInfo struct {
	BrokerID string `json:"broker_id"` // 代征主体ID
	DealerID string `json:"dealer_id"` // 商户ID
	Comment  string `json:"comment"`   // 备注
	RealName string `json:"real_name"` // 姓名
	IDCard   string `json:"id_card"`   // 证件号
	Ref      string `json:"ref"`       // 凭证(上传信息中)
	Status   string `json:"status"`    // 状态(pass: 通过 reject: 拒绝)
}

// QueryInvoiceResponse 查询发票响应信息
type QueryInvoiceResponse struct {
	BaseResponse
	Data InvoiceInfo `json:"data"`
}

// InvoiceInfo 发票信息
type InvoiceInfo struct {
	BrokerID    string `json:"broker_id"`    // 代征主体ID
	DealerID    string `json:"dealer_id"`    // 商户ID
	Invoiced    string `json:"invoiced"`     // 已开发票金额
	NotInvoiced string `json:"not_invoiced"` // 待开发票⾦额
}

// ElementVerifyResponse 银行卡四要素发送短信请求信息
type ElementVerifyResponse struct {
	BaseResponse
	Data struct {
		Ref string `json:"ref"` // 交易凭证
	} `json:"data"`
}

// QueryBankCardResponse 查询银行卡信息响应信息
type QueryBankCardResponse struct {
	BaseResponse
	Data BankCardInfo `json:"data"` // 银行卡信息
}

// BankCardInfo 银行卡信息
type BankCardInfo struct {
	BankCode  string `json:"bank_code"`  // 银行代码
	BankName  string `json:"bank_name"`  // 银行名称
	CardType  string `json:"card_type"`  // 银行卡类型
	IsSupport bool   `json:"is_support"` // 云账户综合服务平台是否支持该银行打款
}

// QueryTaxUserCrossInfo 查询纳税人是否为跨集团用户
type QueryTaxUserCrossInfo struct {
	Year   string `json:"year"`    //用户报税所在年份
	IdCard string `json:"id_card"` //所查询用户的身份证件号码
	EntId  string `json:"ent_id"`  // 商户签约主体 accumulus_tj：天津 accumulus_gs：甘肃
}

// QueryTaxUserCrossResponse 查询纳税人是否为跨集团用户
type QueryTaxUserCrossResponse struct {
	BaseResponse
	Data struct {
		IsCross bool `json:"is_cross"` // 用户是否跨集团标识 false：非跨集团 true：跨集团
	} `json:"data"`
}

// TaxfileDowloadResonse 查询发票信息
type TaxfileDowloadResonse struct {
	BaseResponse
	Data struct {
		FileInfo []TaxfileDowloadInfo `json:"file_info"`
	} `json:"data"`
}

// TaxfileDowloadInfo 发票信息
type TaxfileDowloadInfo struct {
	Name string `json:"name"` //文件名称
	Url  string `json:"url"`  //下载文件临时url
	Pwd  string `json:"pwd"`  //使用公钥加密后的文件解压缩密码
}

// QueryInvoiceAmount 查询可开票额度和开票信息
type QueryInvoiceAmountResponse struct {
	BaseResponse
	Data InvoiceAmount `json:"data"`
}

// 开票额度和开票信息
type InvoiceAmount struct {
	Amount            string                  `json:"amount"`              //可开票额度
	BankNameAccount   []BankNameAccountInfo   `json:"bank_name_account"`   //系统支持的开户行及 账号
	GoodsServicesName []GoodsServicesNameInfo `json:"goods_services_name"` //系统支持的货物或应 税劳务、服务名称
}

//BankNameAccountInfo系统支持的开户行及 账号
type BankNameAccountInfo struct {
	Item    string `json:"item"`    //开户行及账号
	Default bool   `json:"default"` //是否为默认值
}

//GoodsServicesNameInfo系统支持的货物或应 税劳务、服务名称
type GoodsServicesNameInfo struct {
	Item    string `json:"item"`    //货 物 或 应 税 劳 务、服务名称
	Default bool   `json:"default"` //是否为默认值
}

// QueryDayOrdersParam查询日订单数据请求参数
type QueryDayOrdersParam struct {
	OrderDate string `json:"order_date"` //订单查询日期
	Offset    int    `json:"offset"`     //偏移量
	Length    int    `json:"length"`     //条数
	Channel   string `json:"channel"`    //渠道名称
	DataType  string `json:"data_type"`  //数据类型
}

// QueryDayOrdersResponse 查询日订单数据响应信息
type QueryDayOrdersResponse struct {
	BaseResponse
	Data struct {
		TotalNum      int            `json:"total_num"` //总条数
		DayOrderInfos []DayOrderInfo `json:"list"`      //条目信息
	} `json:"data"`
}

// 日订单数据
type DayOrderInfo struct {
	BrokerID            string `json:"broker_id"`             // 代征主体ID
	DealerID            string `json:"dealer_id"`             // 商户ID
	OrderId             string `json:"order_id"`              //商户订单号
	Ref                 string `json:"ref"`                   //流水号
	BatchId             string `json:"batch_id"`              //批次号
	RealName            string `json:"real_name"`             //姓名
	CardNo              string `json:"card_no"`               //收款账号
	BrokerAmount        string `json:"broker_amount"`         //综合服务主体订单金
	BrokerFee           string `json:"broker_fee"`            //综合服务主体服务费
	Bill                string `json:"bill"`                  //渠道流水号
	Status              string `json:"status"`                //订单状态码
	StatusDetail        string `json:"status_detail"`         //订单详细状态码
	StatusMessage       string `json:"status_message"`        //订单状态码描述
	StatusDetailMessage string `json:"status_detail_message"` //订单详细状态码描述
	StatmentId          string `json:"statment_id"`           //短周期授信账单号
	FeeStatmentId       string `json:"fee_statment_id"`       //服务费账单号
	BalStatmentId       string `json:"bal_statment_id"`       //余额账单号
	Channel             string `json:"channel"`               //支付渠道  银行卡、支付宝、微信
	CreatedAt           string `json:"created_at"`            //订单接收时间
	FinishedTime        string `json:"finished_time"`         //订单完成时间
}

//  查询日订单文件（支付和退款订单）响应信息
type QueryDayOrdersFileResponse struct {
	BaseResponse
	Data struct {
		Url string `json:"url"` //下载地址
	} `json:"data"`
}

// QueryDayOrdersParam 查询日流水数据
type QueryDayBillsParam struct {
	BillDate string `json:"bill_date"` //流水查询日期
	Offset   int    `json:"offset"`    //偏移量
	Length   int    `json:"length"`    //条数
	DataType string `json:"data_type"` //数据类型
}

// QueryDayOrdersResponse 查询日流水数据据响应信息
type QueryDayBillsResponse struct {
	BaseResponse
	Data struct {
		TotalNum     int           `json:"total_num"` //总条数
		DayBillInfos []DayBillInfo `json:"list"`      //条目信息
	} `json:"data"`
}

// DayBillInfo 日流水数据
type DayBillInfo struct {
	BrokerID          string `json:"broker_id"`           // 代征主体ID
	DealerID          string `json:"dealer_id"`           // 商户ID
	OrderId           string `json:"order_id"`            //商户订单号
	Ref               string `json:"ref"`                 //流水号
	BrokerProductName string `json:"broker_product_name"` //综合服务主体名称
	DealerProductName string `json:"dealer_product_name"` //商户名称
	BizRef            string `json:"biz_ref"`             //业务订单流水号
	AcctType          string `json:"acct_type"`           //账户类型
	Amount            string `json:"amount"`              //入账金额
	Balance           string `json:"balance"`             //账户余额
	BusinessCategory  string `json:"business_category"`   // 业务分类
	BusinessType      string `json:"business_type"`       //业务类型
	ConsumptionType   string `json:"consumption_type"`    //收支类型
	CreatedAt         string `json:"created_at"`          //订单接收时间
	Remark            string `json:"remark"`              //备注
}

// QueryDayOrdersResponse 查询日流水数据据响应信息
type QueryDailyStatementsResponse struct {
	BaseResponse
	Data struct {
		DailyStatements []DailyStatement `json:"list"` //条目信息
	} `json:"data"`
}

// DailyStatement余额日账单数据
type DailyStatement struct {
	StatementId           string `json:"statement_id"`             // 账单
	StatementDate         string `json:"statement_date"`           // 账单日期
	BrokerID              string `json:"broker_id"`                // 代征主体ID
	DealerID              string `json:"dealer_id"`                // 商户ID
	BrokerProductName     string `json:"broker_product_name"`      //综合服务主体名称
	DealerProductName     string `json:"dealer_product_name"`      //商户名称
	BizType               string `json:"biz_type"`                 //业务类型
	TotalMoney            string `json:"total_money"`              //账单金额
	Amount                string `json:"amount"`                   //订单金额
	ReexAmount            string `json:"reex_amount"`              //退汇金额
	FeeAmount             string `json:"fee_amount"`               // 服务费金额
	DeductRebateFeeAmount string `json:"deduct_rebate_fee_amount"` //服务费抵扣金额
	MoneyAdjust           string `json:"money_adjust"`             //冲补金额
	Status                string `json:"status"`                   //账单状态
	InvoiceStatus         string `json:"invoice_status"`           //开票状态
}

// InvoiceApplyParam 开票申请
type InvoiceApplyParam struct {
	InvoiceApplyId    string `json:"invoice_apply_id"`    //发票申请编号
	Amount            string `json:"amount"`              //申请开票金额
	InvoiceType       string `json:"invoice_type"`        //发票类型
	BankNameAccount   string `json:"bank_name_account"`   //开户行及账号
	GoodsServicesName string `json:"goods_services_name"` //货物或应税劳务、服务名称
	Remark            string `json:"remark"`              //发票备注
}

//开票申请响应信息
type InvoiceApplyResponse struct {
	BaseResponse
	Data InvoiceApply `json:"data"`
}

// 开票申请
type InvoiceApply struct {
	ApplicationId string `json:"application_id"` //发票申请单
	Count         string `json:"count"`          //本次开票申请，发票张数
}
