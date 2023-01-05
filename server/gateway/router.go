package gateway

import (
	"github.com/cloudwego/hertz/pkg/app"
)

type HttpMethod int

const (
	MethodPost   HttpMethod = 1
	MethodGet    HttpMethod = 2
	MethodPut    HttpMethod = 3
	MethodDelete HttpMethod = 4
)

type Router struct {
	Method  HttpMethod
	Version int
	Path    string
	Auth    bool
	Handler app.HandlerFunc
}

type RouterManager interface {
	Routers() []Router
	Name() string
}
