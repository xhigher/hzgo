package controller

import (
	"github.com/xhigher/hzgo/server/service"
)

func New(name string) Controller {
	return Controller{
		&service.Controller{
			Name: name,
		},
	}
}

type Controller struct {
	*service.Controller
}

func (ctrl Controller) Name() string {
	return ctrl.Controller.Name
}

func (ctrl Controller) Routers() []service.Router {
	return []service.Router{
		{
			Version: 1,
			Path:    "register",
			Handler: ctrl.Register,
		},
		{
			Version: 1,
			Path:    "login_check",
			Handler: ctrl.LoginCheck,
		},
		{
			Version: 1,
			Path:    "token_update",
			Handler: ctrl.TokenUpdate,
		},
		{
			Version: 1,
			Path:    "token_check",
			Handler: ctrl.TokenCheck,
		},
		{
			Version: 1,
			Path:    "profile",
			Handler: ctrl.Profile,
		},
	}
}
