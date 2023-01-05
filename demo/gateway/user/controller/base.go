package controller

import (
	"github.com/xhigher/hzgo/server/gateway"
	"github.com/xhigher/hzgo/server/gateway/middlewares"
)

func New(name string) Controller {
	return Controller{
		&gateway.Controller{
			Name: name,
		},
	}
}

func NewWithAuth(name string, auth *middlewares.JWTAuth) Controller {
	return Controller{
		&gateway.Controller{
			Name: name,
			Auth: auth,
		},
	}
}

type Controller struct {
	*gateway.Controller
}

func (ctrl Controller) Name() string {
	return ctrl.Controller.Name
}

func (ctrl Controller) Routers() []gateway.Router {
	return []gateway.Router{
		{
			Method:  gateway.MethodPost,
			Version: 1,
			Path:    "register",
			Auth:    false,
			Handler: ctrl.Register,
		},
		{
			Method:  gateway.MethodPost,
			Version: 1,
			Path:    "login",
			Auth:    false,
			Handler: ctrl.Login,
		},
		{
			Method:  gateway.MethodGet,
			Version: 1,
			Path:    "logout",
			Auth:    true,
			Handler: ctrl.Logout,
		},
		{
			Method:  gateway.MethodGet,
			Version: 1,
			Path:    "renewal",
			Auth:    true,
			Handler: ctrl.Renewal,
		},
		{
			Method:  gateway.MethodGet,
			Version: 1,
			Path:    "profile",
			Auth:    true,
			Handler: ctrl.Profile,
		},
	}
}
