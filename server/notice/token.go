package notice

import (
	"fmt"
	"github.com/golang-jwt/jwt/v4"
	"github.com/xhigher/hzgo/config"
	"github.com/xhigher/hzgo/utils"
	"time"
)

type TokenHelper struct {
	Issuer           string
	SecretKey        []byte
	SigningAlgorithm string
	Timeout          time.Duration
}

func newTokenHelper(conf *config.JWTConfig) TokenHelper {
	return TokenHelper{
		Issuer:           conf.Issuer,
		SecretKey:        []byte(conf.SecretKey),
		SigningAlgorithm: "HS256",
		Timeout:          time.Duration(conf.Timeout) * time.Hour,
	}
}

func (h TokenHelper) CreateToken(subject, audience string) (token string, err error) {
	regClaims := &jwt.RegisteredClaims{
		Issuer:   h.Issuer,
		Subject:  subject,
		Audience: jwt.ClaimStrings{audience},
		ID:       utils.MD5(fmt.Sprintf("%s-%s-%s", audience, subject, utils.UUID())),
	}
	issuedAt := time.Now()
	expiresAt := time.Now().Add(h.Timeout)
	regClaims.ExpiresAt = jwt.NewNumericDate(expiresAt)
	regClaims.IssuedAt = jwt.NewNumericDate(issuedAt)
	jwtToken := jwt.NewWithClaims(jwt.GetSigningMethod(h.SigningAlgorithm), regClaims)
	token, err = jwtToken.SignedString([]byte(h.SecretKey))
	if err != nil {
		return
	}
	return
}

func (h TokenHelper) ParseTokenInfo(token string) (subject, audience string, err error) {
	jwtToken, err := jwt.ParseWithClaims(token, &jwt.RegisteredClaims{}, func(t *jwt.Token) (interface{}, error) {
		if jwt.GetSigningMethod(h.SigningAlgorithm) != t.Method {
			return nil, jwt.ValidationError{}
		}
		return h.SecretKey, nil
	})
	if err != nil {
		return
	}
	claims := jwtToken.Claims.(*jwt.RegisteredClaims)
	subject = claims.Subject
	audience = claims.Audience[0]
	return
}
