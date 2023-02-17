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
			Path:    "config_list",
			Handler: ctrl.ConfigList,
		},
		{
			Version: 1,
			Path:    "config_info",
			Handler: ctrl.ConfigInfo,
		},
		{
			Version: 1,
			Path:    "config_reload",
			Handler: ctrl.ConfigReload,
		},
		{
			Version: 1,
			Path:    "banner_list",
			Handler: ctrl.BannerList,
		},
		{
			Version: 1,
			Path:    "banner_reload",
			Handler: ctrl.BannerReload,
		},
	}
}
