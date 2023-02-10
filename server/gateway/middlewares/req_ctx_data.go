package middlewares

import (
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/xhigher/hzgo/defines"
)

func setToken(c *app.RequestContext, token string) {
	c.Set("JWT_TOKEN", token)
}

func GetToken(c *app.RequestContext) string {
	if token, ok := c.Get("JWT_TOKEN"); ok {
		return token.(string)
	}
	return ""
}

func setSubject(c *app.RequestContext, subject string) {
	c.Set("JWT_SUBJECT", subject)
}

func GetSubject(c *app.RequestContext) string {
	if subject, ok := c.Get("JWT_SUBJECT"); ok {
		return subject.(string)
	}
	return ""
}

func setAudience(c *app.RequestContext, audience string) {
	c.Set("JWT_AUDIENCE", audience)
}

func GetAudience(c *app.RequestContext) string {
	if audience, ok := c.Get("JWT_AUDIENCE"); ok {
		return audience.(string)
	}
	return ""
}

func setBaseParams(c *app.RequestContext, params *defines.BaseParams) {
	c.Set("BASE_PARAMS", params)
}

func GetBaseParams(c *app.RequestContext) *defines.BaseParams {
	if params, ok := c.Get("BASE_PARAMS"); ok {
		return params.(*defines.BaseParams)
	}
	return nil
}