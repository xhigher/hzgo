package admin

import (
	"context"
	"errors"
	"fmt"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/golang-jwt/jwt/v4"
	"github.com/xhigher/hzgo/bizerr"
	"github.com/xhigher/hzgo/config"
	"github.com/xhigher/hzgo/logger"
	"github.com/xhigher/hzgo/resp"
	"github.com/xhigher/hzgo/utils"
	"strings"
	"time"
)

const (
	TokenHeaderKey   = "Authorization"
	TokenValuePrefix = "Bearer"
	SigningAlgorithm = "HS256"
)

var (
	ErrMissingRealm = errors.New("jwt auth realm is required")

	ErrMissingSecretKey = errors.New("jwt auth secret key is required")

	ErrExpiredToken = errors.New("jwt auth token is expired")

	ErrEmptyAuthHeader = errors.New("jwt auth  header is empty")

	ErrInvalidAuthHeader = errors.New("jwt auth  header is invalid")

	ErrInvalidSigningAlgorithm = errors.New("jwt auth invalid signing algorithm")
)

type Claims struct {
	Subject string

	Audience []string

	ExpiredAt int64

	IssuedAt int64

	TokenId string
}

type Auth struct {
	Realm string

	Issuer string

	SecretKey []byte

	Timeout        time.Duration
	MaxRefreshTime time.Duration

	CheckTokenFunc func(ctx context.Context, c *app.RequestContext, claims *Claims) (bool, *bizerr.Error)


}

func NewAuth(conf *config.JWTConfig) *Auth {
	if len(conf.Realm) == 0 {
		panic(ErrMissingRealm)
	}
	if len(conf.SecretKey) == 0 {
		panic(ErrMissingSecretKey)
	}
	mw := &Auth{
		Realm:     conf.Realm,
		SecretKey: []byte(conf.SecretKey),
	}
	if conf.Timeout == 0 {
		conf.Timeout = 1
	}
	mw.Timeout = time.Duration(conf.Timeout) * time.Hour

	if conf.MaxRefreshTime == 0 {
		conf.MaxRefreshTime = 1
	}
	mw.MaxRefreshTime = time.Duration(conf.MaxRefreshTime) * time.Hour

	return mw
}

func (mw *Auth) Handler() app.HandlerFunc {
	return func(ctx context.Context, c *app.RequestContext) {
		claims, err := mw.getClaimsFromToken(ctx, c)
		if err != nil {
			logger.Errorf("GetClaimsFromToken error: %v", err)
			resp.ReplyErrorAuthorization(c)
			return
		}
		logger.Infof("GetClaimsFromToken claims: %v", claims)
		if claims.ExpiresAt == nil {
			resp.ReplyErrorIllegal(c)
			return
		}

		if !claims.VerifyExpiresAt(time.Now(), true) {
			logger.Errorf("VerifyExpiresAt: false")
			resp.ReplyErrorAuthorization(c)
			return
		}

		if mw.CheckTokenFunc != nil {
			ok, be := mw.CheckTokenFunc(ctx, c, mw.getClaims(claims))
			if be != nil {
				logger.Errorf("error: %v", be.String())
				resp.ReplyErr(c, be.ToResp())
				return
			}
			if !ok {
				resp.ReplyErrorAuthorization(c)
				return
			}

		}

		setSubject(c, claims.Subject)
		setAudience(c, claims.Audience)

		c.Next(ctx)
	}
}

func (mw *Auth) getClaimsFromToken(ctx context.Context, c *app.RequestContext) (claims *jwt.RegisteredClaims, err error) {
	token, err := mw.parseToken(ctx, c)
	if err != nil {
		return
	}
	claims = token.Claims.(*jwt.RegisteredClaims)
	return
}

func (mw *Auth) parseToken(ctx context.Context, c *app.RequestContext) (*jwt.Token, error) {
	token, err := mw.getTokenFromHeader(ctx, c)
	if err != nil {
		return nil, err
	}

	return jwt.ParseWithClaims(token, &jwt.RegisteredClaims{}, func(t *jwt.Token) (interface{}, error) {
		if jwt.GetSigningMethod(SigningAlgorithm) != t.Method {
			return nil, ErrInvalidSigningAlgorithm
		}

		setToken(c, token)

		return mw.SecretKey, nil
	})
}

func (mw *Auth) getTokenFromHeader(ctx context.Context, c *app.RequestContext) (string, error) {
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

func (mw *Auth) getClaims(claims *jwt.RegisteredClaims) *Claims {
	return &Claims{
		Audience:  claims.Audience,
		ExpiredAt: claims.ExpiresAt.Unix(),
		IssuedAt:  claims.IssuedAt.Unix(),
		TokenId:   claims.ID,
	}
}

func (mw *Auth) CreateToken(subject string, audience []string) (tokenValue string, claims *Claims, err error) {
	regClaims := &jwt.RegisteredClaims{
		Issuer:   mw.Issuer,
		Subject: subject,
		Audience: audience,
		ID:       utils.MD5(fmt.Sprintf("%s-%s-%s", mw.Issuer, audience, utils.UUID())),
	}
	issuedAt := time.Now()
	expiresAt := time.Now().Add(mw.Timeout)
	regClaims.ExpiresAt = jwt.NewNumericDate(expiresAt)
	regClaims.IssuedAt = jwt.NewNumericDate(issuedAt)
	token := jwt.NewWithClaims(jwt.GetSigningMethod(SigningAlgorithm), regClaims)
	tokenValue, err = token.SignedString([]byte(mw.SecretKey))
	if err != nil {
		return
	}
	claims = mw.getClaims(regClaims)
	return
}

func (mw *Auth) RenewalToken(c *app.RequestContext) (tokenValue string, claims *Claims, err error){
	tokenString := GetToken(c)
	token, err := jwt.ParseWithClaims(tokenString, &jwt.RegisteredClaims{}, func(t *jwt.Token) (interface{}, error) {
		if jwt.GetSigningMethod(SigningAlgorithm) != t.Method {
			return nil, ErrInvalidSigningAlgorithm
		}

		return mw.SecretKey, nil
	})
	if err != nil {
		resp.ReplyErrorAuthorization(c)
		return
	}
	regClaims := token.Claims.(*jwt.RegisteredClaims)
	if !regClaims.VerifyIssuedAt(time.Now().Add(-mw.MaxRefreshTime), true) {
		resp.ReplyErrorAuthorization(c)
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

func setAudience(c *app.RequestContext, audience []string) {
	c.Set("JWT_AUDIENCE", audience)
}

func GetAudience(c *app.RequestContext) []string {
	if audience, ok := c.Get("JWT_AUDIENCE"); ok {
		return audience.([]string)
	}
	return nil
}