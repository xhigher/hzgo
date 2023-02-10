package gateway

import (
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/xhigher/hzgo/consts"
)

type Router struct {
	Method  consts.HttpMethod
	Version int
	Path    string
	Auth    bool
	Sign bool
	Handler app.HandlerFunc
}

type RouterManager interface {
	Routers() []Router
	Name() string
}
