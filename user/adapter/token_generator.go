package adapter

import (
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/stockfolioofficial/back-editfolio/domain"
)

type tokenGenerator struct {
	secret []byte
}

type customClaims struct {
	jwt.StandardClaims
	Roles []string `json:"roles"`
}

func NewTokenGenerateAdapter(secret []byte) domain.TokenGenerateAdapter {
	return &tokenGenerator{
		secret: secret,
	}
}

func (t *tokenGenerator) Generate(u domain.User) (string, error) {
	now := time.Now()
	return jwt.NewWithClaims(jwt.SigningMethodHS256, customClaims{
		StandardClaims: jwt.StandardClaims{
			Subject:  u.Id.String(),
			IssuedAt: now.Unix(),
			// Issuer: , tobe defined
		},
		Roles: []string{string(u.Role)},
	}).SignedString([]byte(t.secret))
}
