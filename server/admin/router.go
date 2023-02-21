package admin

import (
	"context"
	"fmt"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/xhigher/hzgo/consts"
	"github.com/xhigher/hzgo/logger"
	"github.com/xhigher/hzgo/resp"
	"regexp"
	"strconv"
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
		resp := resp.Responder{Ctx: c}
		typ := CRUDRead
		if r.Method == consts.MethodPost {
			typ = CRUDWrite
		}
		myRoles := GetAudience(c)
		if len(myRoles) == 0 {
			resp.ReplyErrorPermission()
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
						resp.ReplyErrorPermission()
						return
					}
				}else{
					if !p.write {
						logger.Warnf("permission error: %v, %v, %v, %v", role, m.Name(), typ, p.write)
						resp.ReplyErrorPermission()
						return
					}
				}
				allow = true
			}
		}
		if !allow {
			resp.ReplyErrorPermission()
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

var (
	pathRegex = regexp.MustCompile(`^/(\w+)\/v(\d+)\/(\w+)$`)
)


func(r Router) FullPath(module string) string {
	return fmt.Sprintf("/%s/v%d/%s", module, r.Version, r.Path)
}

func ParseRouterPath(path string) (module, action string, version int) {
	if len(path) > 0 {
		items := pathRegex.FindStringSubmatch(path)
		if len(items) > 0 {
			module = items[1]
			version, _ = strconv.Atoi(items[2])
			action = items[3]
		}
	}
	return
}
