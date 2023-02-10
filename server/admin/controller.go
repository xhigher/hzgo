package admin

import (
	"github.com/cloudwego/hertz/pkg/app"
)

type Controller struct {
	Auth *Auth
}

func (ctrl Controller) Uid(c *app.RequestContext) string {
	return GetSubject(c)
}

func (ctrl Controller) Token(c *app.RequestContext) string {
	return GetToken(c)
}

func (ctrl Controller) CreateToken(uid string, roles []string) (tokenValue string, claims *Claims, err error) {
	return ctrl.Auth.CreateToken(uid, roles)
}

