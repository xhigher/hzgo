package controller

import (
	"github.com/xhigher/hzgo/router"
)

func New() Controller{
	return Controller{}
}

type Controller struct {

}

func (ctrl Controller) Routers() []router.Router{
	return []router.Router{
		{
			Method: router.MethodPost,
			Version: 1,
			Path: "register",
			Auth: false,
			Handler: ctrl.Register,
		},
	}
}
