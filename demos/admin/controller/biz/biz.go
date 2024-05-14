package biz

import "github.com/xhigher/hzgo/server/admin"

func Modules(ctrl *admin.Controller) []admin.Module {
	return []admin.Module{
		UserModule{
			ctrl: ctrl,
		},
	}
}
