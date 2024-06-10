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
	Name    string
	Path    string
	NoAuth  bool
	Handler app.HandlerFunc
}

func (r Router) PermissionFunc() app.HandlerFunc {
	return func(ctx context.Context, c *app.RequestContext) {
		resp := resp.Responder{Ctx: c}
		roles := GetAudience(c)
		if len(roles) == 0 {
			resp.ReplyErrorPermission()
			return
		}
		allow := false
		logger.Infof("staff.roles: %v", roles)
		for _, role := range roles {
			if !CheckPermission(role, r.Path) {
				logger.Warnf("permission error: %v, %v", role, r.Path)
				resp.ReplyErrorPermission()
				return
			}
			allow = true
		}
		if !allow {
			resp.ReplyErrorPermission()
			return
		}

		c.Next(ctx)
	}

}

var (
	pathRegex = regexp.MustCompile(`^/(\w+)\/v(\d+)\/(\w+)$`)
)

func (r Router) FullPath(module string) string {
	return fmt.Sprintf("/%s/v%d/%s", module, r.Version, r.Name)
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
