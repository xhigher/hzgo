package admin

import (
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/xhigher/hzgo/resp"
)

type Controller struct {
	Auth *Auth
	LogSaver resp.TraceLogSaver
}

func (ctrl Controller) Resp(c *app.RequestContext) resp.Responder {
	return resp.Responder{Ctx:c, LogSaver: ctrl.LogSaver, LogOut: true}
}

func (ctrl Controller) Uid(c *app.RequestContext) string {
	return GetSubject(c)
}

func (ctrl Controller) Token(c *app.RequestContext) string {
	return GetToken(c)
}

func (ctrl Controller) CreateToken(c *app.RequestContext, uid string, roles []string) (tokenValue string, claims *Claims, err error) {
	return ctrl.Auth.CreateToken(c, uid, roles)
}

func (ctrl Controller) Roles(c *app.RequestContext) []string {
	return GetAudience(c)
}
func (ctrl Controller) isRoleType(c *app.RequestContext, typ string) bool{
	roles := GetAudience(c)
	for _, role := range roles {
		if role == typ {
			return true
		}
	}
	return false
}

func (ctrl Controller) IsRoleMaintainer(c *app.RequestContext) bool{
	return ctrl.isRoleType(c, RoleMaintainer)
}

func (ctrl Controller) IsRoleDeveloper(c *app.RequestContext) bool{
	return ctrl.isRoleType(c, RoleDeveloper)
}

func (ctrl Controller) IsRoleOperator(c *app.RequestContext) bool{
	return ctrl.isRoleType(c, RoleOperator)
}

func (ctrl Controller) IsRoleTreasurer(c *app.RequestContext) bool{
	return ctrl.isRoleType(c, RoleTreasurer)
}

func (ctrl Controller) IsRoleCustomer(c *app.RequestContext) bool{
	return ctrl.isRoleType(c, RoleCustomer)
}


