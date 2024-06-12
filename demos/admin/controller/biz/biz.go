package biz

import "github.com/xhigher/hzgo/server/admin"

func Modules(ctrl *admin.Controller) []admin.Module {
	return []admin.Module{
		ConfigModule{
			ctrl: ctrl,
		},
		UserModule{
			ctrl: ctrl,
		},
	}
}
