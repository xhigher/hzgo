package admin

import (
	"context"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/xhigher/hzgo/consts"
	"github.com/xhigher/hzgo/logger"
	"github.com/xhigher/hzgo/resp"
)

type Router struct {
	Method  consts.HttpMethod
	Version int
	Path    string
	NoAuth bool
	Roles []string
	Handler app.HandlerFunc
}

func(r Router) PermissionFunc(module string) app.HandlerFunc {
	return func(ctx context.Context, c *app.RequestContext) {
		typ := CRUDRead
		if r.Method == consts.MethodPost {
			typ = CRUDWrite
		}
		myRoles := GetAudience(c)
		for _, role := range myRoles {
			if r.containRole(role) {
				p := getModulePermission(role, module)
				if typ == CRUDRead {
					if !p.read {
						logger.Warnf("permission error: %v, %v, %v, %v", role, module, typ, p.read)
						resp.ReplyErrorPermission(c)
						return
					}
				}else{
					if !p.write {
						logger.Warnf("permission error: %v, %v, %v, %v", role, module, typ, p.write)
						resp.ReplyErrorPermission(c)
						return
					}
				}
			}
		}

		c.Next(ctx)
	}

}

func(r Router) containRole(role string) bool {
	for _, tr := range r.Roles {
		if tr == role {
			return true
		}
	}
	return false
}
