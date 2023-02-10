package middlewares

import (
	"context"
	"errors"
	"fmt"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/golang-jwt/jwt/v4"
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

type AuthClaims struct {
	Subject string

	Audience string

	ExpiredAt int64

	IssuedAt int64

	TokenId string
}

type JWTAuth struct {
	Realm string

	Issuer string

	SecretKey []byte

	Timeout        time.Duration
	MaxRefreshTime time.Duration

	CheckTokenFunc func(ctx context.Context, c *app.RequestContext, claims *AuthClaims) bool
}

func NewJWTAuth(conf *config.JWTConfig) *JWTAuth {
	if len(conf.Realm) == 0 {
		panic(ErrMissingRealm)
	}
	if len(conf.SecretKey) == 0 {
		panic(ErrMissingSecretKey)
	}
	mw := &JWTAuth{
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

func (mw *JWTAuth) Authenticate() app.HandlerFunc {
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
			if !mw.CheckTokenFunc(ctx, c, mw.getClaims(claims)) {
				logger.Infof("AuthorizationFunc: false")
				resp.ReplyErrorAuthorization(c)
				return
			}
		}

		setSubject(c, claims.Subject)
		setAudience(c, claims.Audience[0])

		c.Next(ctx)
	}
}

func (mw *JWTAuth) getClaimsFromToken(ctx context.Context, c *app.RequestContext) (claims *jwt.RegisteredClaims, err error) {
	token, err := mw.parseToken(ctx, c)
	if err != nil {
		return
	}
	claims = token.Claims.(*jwt.RegisteredClaims)
	return
}

func (mw *JWTAuth) parseToken(ctx context.Context, c *app.RequestContext) (*jwt.Token, error) {
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

func (mw *JWTAuth) getTokenFromHeader(ctx context.Context, c *app.RequestContext) (string, error) {
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

func (mw *JWTAuth) getClaims(claims *jwt.RegisteredClaims) *AuthClaims {
	return &AuthClaims{
		Subject: claims.Subject,
		Audience:  claims.Audience[0],
		ExpiredAt: claims.ExpiresAt.Unix(),
		IssuedAt:  claims.IssuedAt.Unix(),
		TokenId:   claims.ID,
	}
}

func (mw *JWTAuth) CreateToken(subject, audience string) (tokenValue string, claims *AuthClaims, err error) {
	regClaims := &jwt.RegisteredClaims{
		Issuer:   mw.Issuer,
		Subject: subject,
		Audience: jwt.ClaimStrings{audience},
		ID:       utils.MD5(fmt.Sprintf("%s-%s-%s", audience, subject, utils.UUID())),
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

func (mw *JWTAuth) RenewalToken(c *app.RequestContext) (tokenValue string, claims *AuthClaims, err error){
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

