package service

import (
	"github.com/cloudwego/hertz/pkg/app"
)

type Router struct {
	Version int
	Path    string
	Handler app.HandlerFunc
}

type RouterManager interface {
	Routers() []Router
	Name() string
}
