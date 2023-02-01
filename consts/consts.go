package consts

import "time"

const (
	YES = 1
	NO  = 0

	StatusEditing = 0
	StatusOnline  = 1
	StatusOffline = 2

	UserStatusNormal    = 1
	UserStatusBlocked   = 2
	UserStatusCancelled = 3

	UserCertifyRealFace int32 = 1
	UserCertifyRealName int32 = 2
)

type MediaType int

const (
	MediaTypeImage MediaType = 1 // 图片
	MediaTypeVideo MediaType = 2 // 视频
	MediaTypeAudio MediaType = 3 // 音频
	MediaTypeData MediaType = 4 // 数据文件
)

type MediaBiz struct {
	Dir     string
	Exp time.Duration
}

const (
	WithdrawalAccountBankpay = 1
	WithdrawalAccountAlipay  = 2
	WithdrawalAccountWxpay   = 3

	WithdrawalStatusCreated  = 0
	WithdrawalStatusDeducted = 1
	WithdrawalStatusPaying   = 2
	WithdrawalStatusSuccess  = 3
	WithdrawalStatusFailed   = 4
	WithdrawalStatusCanceled = 5
	WithdrawalStatusRefunded = 6
	WithdrawalStatusWaiting  = 7

	WithdrawalCheckUnwanted = 0
	WithdrawalCheckWaiting  = 1
	WithdrawalCheckPassed   = 2
	WithdrawalCheckRejected = 3
)