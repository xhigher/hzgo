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
			Name:    "config_list",
			Handler: ctrl.ConfigList,
		},
		{
			Name:    "config_info",
			Handler: ctrl.ConfigInfo,
		},
		{
			Name:    "config_reload",
			Handler: ctrl.ConfigReload,
		},
		{
			Name:    "banner_list",
			Handler: ctrl.BannerList,
		},
		{
			Name:    "banner_reload",
			Handler: ctrl.BannerReload,
		},
	}
}
