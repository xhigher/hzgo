package yzhpay

// 基础订单状态
const (
	OrderDelete    = -1 // 订单删除(被标记为删除的订单，只有通过 Web 页面提交批量 订单的情况才会出现（最终态，不会回调），API 接口打款不会出现此状态)
	OrderAccept    = 0  // 订单已受理(支付订单接收成功，尚未支付（中间态，不会回调）)
	OrderSuccess   = 1  // 订单已打款(订单提交到支付网关成功（中间状态，会回调）)
	OrderFailed    = 2  // 订单已失败(主要表示订单数据校验不通过（最终态，会回调）)
	OrderPending   = 4  // 订单待打款(暂停处理，满足条件后会继续支付，例如账户余额 不足，充值后可以继续打款（中间态，会回调）)
	OrderSending   = 5  // 订单打款中(调用支付网关超时等状态异常情况导致，处于等待 交易查证的中间状态（中间态，不会回调）)
	OrderReadySend = 8  // 订单待打款(订单结算限额检查和风控判断完毕，等待执行打款 的状态（中间态，不会回调）)
	OrderReturned  = 9  // 订单已退汇(支付被退回（最终态，会回调）)
	OrderCancel    = 15 // 订单取消(表示待打款（暂停处理）订单数据被商户主动取消 （最终态，会回调）)
)

// 路由信息
const (
	BaseURL = "https://api-service.yunzhanghu.com" // 基础url

	// 实时结算接口url
	BankOrderURL        = "/api/payment/v1/order-bankpay"       // 银行卡下单接口url
	AliOrderURL         = "/api/payment/v1/order-alipay"        // 支付宝下单接口url
	WxOrderURL          = "/api/payment/v1/order-wxpay"         // 微信下单接口url
	QueryOrderURL       = "/api/payment/v1/query-order"         // 查单接口url
	CancelOrderURL      = "/api/payment/v1/order/fail"          // 取消订单url
	QueryAccountURL     = "/api/payment/v1/query-accounts"      // 查询账户信息url
	QueryVaAccountURL   = "/api/payment/v1/va-account"          // 查询商户VA账户url
	QueryReceiptFileURL = "/api/payment/v1/receipt/file"        // 查询电子回单URL
	QueryRechargeURL    = "/api/dataservice/v2/recharge-record" // 查询充值记录url

	// 数据接口url
	DownloadOrderURL        = "/api/dataservice/v1/order/downloadurl" // 下载日订单url
	DownloadBillURL         = "/api/dataservice/v2/bill/downloadurl"  // 下载日流水url
	QueryDayOrdersDataUrl   = "/api/dataservice/v1/orders"            // 查询日订单数据
	QueryDayOrdersFileUrl   = "/api/dataservice/v1/order/day/url"     // 查询⽇订单⽂件 (结算和退款订单)
	QueryDayBillUrl         = "/api/dataservice/v1/bills"             // 查询⽇流⽔数据
	QueryDailyStatementsUrl = "/api/dataservice/v1/statements-daily"  // 查询余额日账单数据

	// 用户信息验证接口url
	UploadUserURL      = "/api/payment/v1/user/exempted/info"           // 上传用户免验证名单url
	CheckExistUserURL  = "/api/payment/v1/user/white/check"             // 校验免验证用户名单是否存在url
	Element4RequestURL = "/authentication/verify-request"               // 银行卡四要素鉴权发送短信url
	Element4ConfirmURL = "/authentication/verify-confirm"               // 银行卡四要素鉴权提交验证码url
	Element4URL        = "/authentication/verify-bankcard-four-factor"  // 银行卡四要素鉴权url
	Element3URL        = "/authentication/verify-bankcard-three-factor" // 银行卡三要素鉴权url
	IDCheckURL         = "/authentication/verify-id"                    // 实名制二要素鉴权url
	BankCardInfoURL    = "/api/payment/v1/card"                         // 银行卡信息查询url

	// 发票接口url
	QueryInvoiceURL       = "/api/payment/v1/invoice-stat"           // 查询发票接口
	QueryInvoiceAmountUrl = "/api/invoice/v2/invoice-amount"         //查询可开票额度和开票信息
	InvoiceApplyUrl       = "/api/invoice/v2/apply"                  // 开票申请
	InvoiceStatusUrl      = "/api/invoice/v2/invoice/invoice-status" // 查询开票申请状态
	InvoicePdfUrl         = "/api/invoice/v2/invoice/invoice-pdf"    // 下载发票 PDF
	InvoiceEmailUrl       = "/api/invoice/v2/invoice/reminder/email" // 发送发票扫描件压缩包下载链接邮件

	//个税扣缴明细表下载接口url
	QueryTaxUserCrossUrl = "/api/tax/v1/user/cross"       // 查询纳税人是否为跨集团用户
	TaxFileDownloadUrl   = "/api/tax/v1/taxfile/download" // 下载个税扣缴明细表

)
