package auth

import (
	"context"
	"errors"
	"fmt"
	"github.com/xhigher/hzgo/resp"
	"github.com/xhigher/hzgo/utils"
	"strings"
	"time"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/golang-jwt/jwt/v4"
)

const (
	TokenHeaderKey = "Authorization"
	TokenValuePrefix = "Bearer"
	SigningAlgorithm = "HS256"
	DefaultRealm = "xhigher hzgo"
)

type HzgoToken struct {
	Token string `json:"token"`
	Et int64 `json:"et"`
}

type HzgoClaims struct {
	Audience string

	ExpiredAt int64

	IssuedAt int64

	TokenId string
}

type HzgoJWTMiddleware struct {

	Realm string

	Issuer string

	SecretKey string

	Timeout time.Duration

	MaxRefresh time.Duration

	Controller HzgoAuthController

	AuthenticationFunc func(ctx context.Context, c *app.RequestContext) (string, error)

	AuthorizationFunc func(data interface{}, ctx context.Context, c *app.RequestContext) bool

	LoginFunc func(ctx context.Context, c *app.RequestContext, claims HzgoClaims) bool

	LogoutFunc func(ctx context.Context, c *app.RequestContext, audience string) bool

	RefreshFunc func(ctx context.Context, c *app.RequestContext, claims HzgoClaims) bool

}

var (
	ErrMissingSecretKey = errors.New("secret key is required")

	ErrExpiredToken = errors.New("token is expired")

	ErrAuthenticationFuncNil = errors.New("authentication func is nil")

	ErrEmptyAuthHeader = errors.New("auth header is empty")

	ErrInvalidAuthHeader = errors.New("auth header is invalid")

	ErrInvalidSigningAlgorithm = errors.New("invalid signing algorithm")

)

// New for check error with HertzJWTMiddleware
func New(m *HzgoJWTMiddleware) (*HzgoJWTMiddleware, error) {
	if err := m.MiddlewareInit(); err != nil {
		return nil, err
	}

	return m, nil
}



// MiddlewareInit initialize jwt configs.
func (mw *HzgoJWTMiddleware) MiddlewareInit() error {
	if len(mw.Realm) == 0 {
		mw.Realm = DefaultRealm
	}
	if mw.Timeout == 0 {
		mw.Timeout = time.Hour
	}
	if mw.MaxRefresh == 0 {
		mw.MaxRefresh = time.Hour
	}
	if len(mw.SecretKey) == 0 {
		return ErrMissingSecretKey
	}

	if mw.Controller == nil {
		return ErrAuthenticationFuncNil
	}

	if mw.AuthenticationFunc == nil {
		return ErrAuthenticationFuncNil
	}

	if mw.LoginFunc == nil {
		mw.LoginFunc = func(ctx context.Context, c *app.RequestContext, claims HzgoClaims) bool {
			return true
		}
	}
	if mw.RefreshFunc == nil {
		mw.RefreshFunc = func(ctx context.Context, c *app.RequestContext, claims HzgoClaims) bool {
			return true
		}
	}
	if mw.LogoutFunc == nil {
		mw.LogoutFunc = func(ctx context.Context, c *app.RequestContext, audience string) bool {
			return true
		}
	}

	if mw.AuthorizationFunc == nil {
		mw.AuthorizationFunc = func(data interface{}, ctx context.Context, c *app.RequestContext) bool {
			return true
		}
	}

	return nil
}

func (mw *HzgoJWTMiddleware) MiddlewareFunc() app.HandlerFunc {
	return func(ctx context.Context, c *app.RequestContext) {
		mw.middlewareImpl(ctx, c)
	}
}

func (mw *HzgoJWTMiddleware) middlewareImpl(ctx context.Context, c *app.RequestContext) {
	claims, err := mw.GetClaimsFromToken(ctx, c)
	if err != nil {
		resp.ReplyErrorAuthorization(c)
		return
	}

	if claims.ExpiresAt == nil {
		resp.ReplyErrorIllegal(c)
		return
	}

	if !claims.VerifyExpiresAt(time.Now(), true) {
		resp.ReplyErrorAuthorization(c)
		return
	}

	if !mw.AuthorizationFunc(claims.Audience, ctx, c) {
		resp.ReplyErrorAuthorization(c)
		return
	}

	c.Next(ctx)
}

func (mw *HzgoJWTMiddleware) GetClaimsFromToken(ctx context.Context, c *app.RequestContext) (claims jwt.RegisteredClaims, err error) {
	token, err := mw.ParseToken(ctx, c)
	if err != nil {
		return
	}
	claims = token.Claims.(jwt.RegisteredClaims)
	return
}

func (mw *HzgoJWTMiddleware) LoginHandler(ctx context.Context, c *app.RequestContext) {
	if mw.AuthenticationFunc == nil {
		resp.ReplyErrorInternal(c)
		return
	}

	audience, err := mw.AuthenticationFunc(ctx, c)
	if err != nil {
		resp.ReplyErrorAuthorization(c)
		return
	}

	token, claims, err := mw.createToken(audience)
	if err != nil {
		resp.ReplyErrorAuthorization(c)
		return
	}

	if mw.LoginFunc != nil {
		if !mw.LoginFunc(ctx, c, claims) {
			resp.ReplyErrorAuthorization(c)
			return
		}
	}

	resp.ReplyData(c, HzgoToken{
		Token: token,
		Et: claims.ExpiredAt,
	})
}

func (mw *HzgoJWTMiddleware) createToken(audience string) (tokenValue string, claims HzgoClaims, err error) {
	regClaims := jwt.RegisteredClaims{
		Issuer: mw.Issuer,
		Audience: jwt.ClaimStrings{audience},
		ID: utils.MD5(fmt.Sprintf("%s-%s-%s", mw.Issuer, audience,utils.UUID())),
	}
	issuedAt := time.Now()
	expiresAt := time.Now().Add(mw.Timeout)
	regClaims.ExpiresAt = jwt.NewNumericDate(expiresAt)
	regClaims.IssuedAt = jwt.NewNumericDate(issuedAt)
	token := jwt.NewWithClaims(jwt.GetSigningMethod(SigningAlgorithm), regClaims)
	tokenValue, err = token.SignedString(mw.SecretKey)
	if err != nil {
		return
	}
	claims = mw.getClaims(regClaims)
	return
}

func (mw *HzgoJWTMiddleware) getClaims(claims jwt.RegisteredClaims) HzgoClaims{
	return HzgoClaims{
		Audience: claims.Audience[0],
		ExpiredAt: claims.ExpiresAt.Unix(),
		IssuedAt: claims.IssuedAt.Unix(),
		TokenId: claims.ID,
	}
}

func (mw *HzgoJWTMiddleware) LogoutHandler(ctx context.Context, c *app.RequestContext) {
	claims, err := mw.CheckIfTokenExpire(ctx, c)
	if err != nil {
		return
	}

	if mw.LogoutFunc != nil {
		if !mw.LogoutFunc(ctx, c, claims.Audience[0]) {
			resp.ReplyNOK(c)
			return
		}
	}

	resp.ReplyOK(c)
}

func (mw *HzgoJWTMiddleware) RefreshHandler(ctx context.Context, c *app.RequestContext) {
	token, claims, err := mw.RefreshToken(ctx, c)
	if err != nil {
		resp.ReplyErrorAuthorization(c)
		return
	}

	if mw.RefreshFunc != nil {
		if !mw.RefreshFunc(ctx, c, claims) {
			resp.ReplyErrorAuthorization(c)
			return
		}
	}

	resp.ReplyData(c, HzgoToken{
		Token: token,
		Et: claims.ExpiredAt,
	})
}

func (mw *HzgoJWTMiddleware) RefreshToken(ctx context.Context, c *app.RequestContext) (tokenValue string, claims HzgoClaims, err error) {
	regClaims, err := mw.CheckIfTokenExpire(ctx, c)
	if err != nil {
		return
	}

	issuedAt := time.Now()
	expiresAt := time.Now().Add(mw.Timeout)
	regClaims.ExpiresAt = jwt.NewNumericDate(expiresAt)
	regClaims.IssuedAt = jwt.NewNumericDate(issuedAt)
	newToken := jwt.NewWithClaims(jwt.GetSigningMethod(SigningAlgorithm), regClaims)
	tokenValue, err = newToken.SignedString(mw.SecretKey)
	if err != nil {
		return
	}
	claims = mw.getClaims(regClaims)
	return
}

func (mw *HzgoJWTMiddleware) CheckIfTokenExpire(ctx context.Context, c *app.RequestContext) (claims jwt.RegisteredClaims, err error) {
	token, err := mw.ParseToken(ctx, c)
	if err != nil {
		validationErr, ok := err.(*jwt.ValidationError)
		if !ok || validationErr.Errors != jwt.ValidationErrorExpired {
			return
		}
		err = nil
	}

	claims = token.Claims.(jwt.RegisteredClaims)
	if !claims.VerifyIssuedAt(time.Now().Add(-mw.MaxRefresh), true) {
		err = ErrExpiredToken
		return
	}

	return
}

func (mw *HzgoJWTMiddleware) getTokenFromHeader(ctx context.Context, c *app.RequestContext) (string, error) {
	authHeader := c.Request.Header.Get(TokenHeaderKey)

	if authHeader == "" {
		return "", ErrEmptyAuthHeader
	}

	parts := strings.SplitN(authHeader, " ", 2)
	if !(len(parts) == 2 && parts[0] == TokenValuePrefix) {
		return "", ErrInvalidAuthHeader
	}

	return parts[1], nil
}

func (mw *HzgoJWTMiddleware) ParseToken(ctx context.Context, c *app.RequestContext) (*jwt.Token, error) {
	token, err := mw.getTokenFromHeader(ctx, c)
	if err != nil {
		return nil, err
	}

	return jwt.Parse(token, func(t *jwt.Token) (interface{}, error) {
		if jwt.GetSigningMethod(SigningAlgorithm) != t.Method {
			return nil, ErrInvalidSigningAlgorithm
		}

		c.Set("JWT_TOKEN", token)

		return mw.SecretKey, nil
	})
}

// ParseTokenString parse jwt token string
func (mw *HzgoJWTMiddleware) ParseTokenString(token string) (*jwt.Token, error) {
	return jwt.Parse(token, func(t *jwt.Token) (interface{}, error) {
		if jwt.GetSigningMethod(SigningAlgorithm) != t.Method {
			return nil, ErrInvalidSigningAlgorithm
		}

		return mw.SecretKey, nil
	})
}


