package auth

import (
	"fmt"

	"github.com/Hajdudev/invoice-flow/internal/env"
	"github.com/golang-jwt/jwt/v5"
)

type Authenticator interface {
	GenerateToken(claims jwt.Claims) (string, error)
	ValidateToken(token string) (*jwt.Token, error)
}

var (
	secret = env.GetString("JWT_SECRET", "xxxx")
	aud    = env.GetString("JWT_AUD", "aud")
	iss    = env.GetString("ISS", "iss")
)

type JWTAuthenticator struct {
	secret string
	aud    string
	iss    string
}

func NewJWTAuthenticator() *JWTAuthenticator {
	return &JWTAuthenticator{secret, iss, aud}
}

func (a *JWTAuthenticator) GenerateToken(claims jwt.Claims) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString([]byte(a.secret))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func (a *JWTAuthenticator) ValidateToken(token string) (*jwt.Token, error) {
	return jwt.Parse(token, func(t *jwt.Token) (any, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method %v", t.Header["alg"])
		}

		return []byte(a.secret), nil
	},
		jwt.WithExpirationRequired(),
		jwt.WithAudience(a.aud),
		jwt.WithIssuer(a.aud),
		jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Name}),
	)
}
