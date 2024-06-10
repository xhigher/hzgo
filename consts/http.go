package consts

import hzconsts "github.com/cloudwego/hertz/pkg/protocol/consts"

type HttpMethod string

const (
	MethodGet    HttpMethod = hzconsts.MethodGet
	MethodPost   HttpMethod = hzconsts.MethodPost
	MethodPut    HttpMethod = hzconsts.MethodPut
	MethodDelete HttpMethod = hzconsts.MethodDelete
)
