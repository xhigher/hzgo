package base

import (
	"context"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/xhigher/hzgo/auth"
	usermodel "github.com/xhigher/hzgo/demo/model/user"
	userlogic "github.com/xhigher/hzgo/demo/logic/user"
)

type LoginReq struct {
	Username  string `form:"username" json:"username" query:"username" vd:"(len($) > 0 && len($) < 30); msg:'Illegal format'"`
	Password string `form:"password" json:"password" query:"password" vd:"(len($) > 0 && len($) < 30); msg:'Illegal format'"`

}

func NewAuthMiddleware() *auth.HzgoJWTMiddleware{
	return &auth.HzgoJWTMiddleware{
		SecretKey: "",
		LoginFunc: func(ctx context.Context, c *app.RequestContext, claims auth.HzgoClaims) bool {
			err := usermodel.SaveToken(claims.Audience, claims.TokenId, claims.ExpiredAt, claims.IssuedAt)
			if err != nil {
				return false
			}
			return true
		},
		LogoutFunc: func(ctx context.Context, c *app.RequestContext, audience string) bool {
			return true
		},
		RefreshFunc: func(ctx context.Context, c *app.RequestContext, claims auth.HzgoClaims) bool {
			err := usermodel.SaveToken(claims.Audience, claims.TokenId, claims.ExpiredAt, claims.IssuedAt)
			if err != nil {
				return false
			}
			return true
		},
		AuthenticationFunc: func(ctx context.Context, c *app.RequestContext) (string, error) {
			var req LoginReq
			if err := c.BindAndValidate(&req); err != nil {
				return "", err
			}
			userid, err := userlogic.CheckUser(req.Username, req.Password)
			if err != nil {
				return "", err
			}

			return userid, nil
		},
		AuthorizationFunc: func(ctx context.Context, c *app.RequestContext, claims auth.HzgoClaims) bool {
			ok, _ := usermodel.CheckToken(claims.Audience, claims.TokenId)
			return ok
		},
	}
}




