package test

import (
	"fmt"
	"github.com/golang-jwt/jwt/v4"
	"github.com/xhigher/hzgo/utils"
	"testing"
	"time"
)

type TokenHelper struct {
	Issuer           string
	SecretKey        []byte
	SigningAlgorithm string
	Timeout          time.Duration
}

type TokenInfo struct {
	Uid string `json:"uid"`
	Did string `json:"did"`
}

func (h TokenHelper) CreateToken(subject, audience string) (tokenValue string, err error) {
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
	token := jwt.NewWithClaims(jwt.GetSigningMethod(h.SigningAlgorithm), regClaims)
	tokenValue, err = token.SignedString([]byte(h.SecretKey))
	if err != nil {
		return
	}
	return
}

func (h TokenHelper) ParseTokenInfo(token string) (claims *jwt.RegisteredClaims, err error) {
	jwtToken, err := jwt.ParseWithClaims(token, &jwt.RegisteredClaims{}, func(t *jwt.Token) (interface{}, error) {
		if jwt.GetSigningMethod(h.SigningAlgorithm) != t.Method {
			return nil, jwt.ValidationError{}
		}
		return h.SecretKey, nil
	})
	if err != nil {
		return
	}
	claims = jwtToken.Claims.(*jwt.RegisteredClaims)
	return
}

func TestToken(t *testing.T) {
	th := TokenHelper{
		Issuer:           "xhigher",
		SecretKey:        []byte("xhigher123"),
		SigningAlgorithm: "HS256",
		Timeout:          time.Duration(1) * time.Hour,
	}
	uid := "123456"
	app := "web"
	ti, err := th.CreateToken(uid, app)
	fmt.Println("ti=", ti, "err=", err)
}

func TestToken2(t *testing.T) {
	th := TokenHelper{
		Issuer:           "xhigher",
		SecretKey:        []byte("xhigher123"),
		SigningAlgorithm: "HS256",
		Timeout:          time.Duration(1) * time.Hour,
	}
	token := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpc3MiOiJ4aGlnaGVyIiwic3ViIjoiMTIzNDU2IiwiYXVkIjpbIndlYiJdLCJleHAiOjE3MTU1ODYwMTYsImlhdCI6MTcxNTU4MjQxNiwianRpIjoiYTgwNDZmZDk3YjRjZTM3NmU4MjM3NTFhYWJlN2MyNWQifQ.TvMIBZBtyhrzLB0hBxhatTYEg_SkXfIqx-E4YgCQ7GE"
	ti, err := th.ParseTokenInfo(token)

	fmt.Println("ti=", utils.JSONString(ti), "err=", err)
}
