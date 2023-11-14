package jwtx

import (
	"time"

	"github.com/golang-jwt/jwt"
)

type Claims struct {
	jwt.StandardClaims

	Mark    string      `json:"mark,omitempty"`
	Content interface{} `json:"content,omitempty"`
}

// NewClaims creates a new Claims with defined content, mark, isSer, expire & secret
func NewClaims(content interface{}, mark, isSer string, expire time.Duration, secret Secret) *Claims {
	return &Claims{
		Mark:    mark,
		Content: content,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(expire).Unix(),
			Issuer:    isSer,
		},
	}
}

// GetToken gets token string with defined signingMethod, secret & claims
func (c *Claims) GetToken(method jwt.SigningMethod, secret Secret) (string, error) {
	tokenClaims := jwt.NewWithClaims(method, c)
	return tokenClaims.SignedString(secret())
}

// ParseToken parses token string with defined callback
func (c *Claims) ParseToken(token string, callback KeyFunc) (*Claims, error) {
	tc, err := jwt.ParseWithClaims(token, c, func(t *jwt.Token) (interface{}, error) {
		return callback((*Token)(t))
	})
	if err != nil {
		return nil, err
	}

	if tc == nil {
		return nil, ErrorInvalidToken
	}

	if claims, ok := tc.Claims.(*Claims); ok && tc.Valid {
		return claims, nil
	}

	return nil, ErrorInvalidToken
}
