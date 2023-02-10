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

func(r Router) PermissionFunc(m Module) app.HandlerFunc {
	return func(ctx context.Context, c *app.RequestContext) {
		typ := CRUDRead
		if r.Method == consts.MethodPost {
			typ = CRUDWrite
		}
		myRoles := GetAudience(c)
		if len(myRoles) == 0 {
			resp.ReplyErrorPermission(c)
			return
		}
		allow := false
		logger.Infof("router.roles: %v, my.roles: %v", r.Roles, myRoles)
		for _, role := range myRoles {
			if r.containRole(role) {
				p := getModulePermission(role, m.Name())
				if typ == CRUDRead {
					if !p.read {
						logger.Warnf("permission error: %v, %v, %v, %v", role, m.Name(), typ, p.read)
						resp.ReplyErrorPermission(c)
						return
					}
				}else{
					if !p.write {
						logger.Warnf("permission error: %v, %v, %v, %v", role, m.Name(), typ, p.write)
						resp.ReplyErrorPermission(c)
						return
					}
				}
				allow = true
			}
		}
		if !allow {
			resp.ReplyErrorPermission(c)
			return
		}

		c.Next(ctx)
	}

}

func(r *Router) mergeRoles(roles []string) {
	if len(r.Roles) == 0 {
		r.Roles = roles
		return
	}
	roles2 := r.Roles
	rm := map[string]int{}
	for _, tr := range r.Roles {
		rm[tr] = 1
	}
	for _, tr := range roles {
		if _,ok := rm[tr]; !ok {
			roles2 = append(roles2, tr)
		}
	}
	r.Roles = roles2
}

func(r Router) containRole(role string) bool {
	for _, tr := range r.Roles {
		if tr == role {
			return true
		}
	}
	return false
}
