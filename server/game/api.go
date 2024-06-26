package game

import (
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/xhigher/hzgo/consts"
)

type Router struct {
	Method  consts.HttpMethod
	Name    string
	Handler app.HandlerFunc
}

type API interface {
	Routers() []Router
}
