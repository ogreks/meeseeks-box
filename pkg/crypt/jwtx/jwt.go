package jwtx

import (
	"time"

	"github.com/golang-jwt/jwt"
)

type JWT struct {
	SigningMethod jwt.SigningMethod
	Secret        Secret
	Claims        *Claims
}

// NewJWT creates a new JWT with defined signingMethod, secret & claims
func NewJWT(signingMethod SigningMethod, claims ...*Claims) *JWT {
	j := &JWT{
		SigningMethod: signingMethod,
	}

	if len(claims) > 0 {
		j.Claims = claims[0]
	}

	return j
}

// SetClaims sets claims
func (j *JWT) SetClaims(claims *Claims) *JWT {
	j.Claims = claims
	return j
}

// BuildClaims builds claims with defined content, mark, isSer, expire & secret
func (j *JWT) BuildClaims(content interface{}, mark, isSer string, expire time.Duration) *JWT {
	j.Claims = NewClaims(content, mark, isSer, expire, j.Secret)

	return j
}

// BuildSecret builds secret with defined secret
func (j *JWT) BuildSecret(secret Secret) *JWT {
	j.Secret = secret
	return j
}

// GetToken gets token string with defined signingMethod, secret & claims
func (j *JWT) GetToken() (string, error) {
	if j.Claims == nil {
		return "", ErrorInvalidCliam
	}

	return j.Claims.GetToken(j.SigningMethod, j.Secret)
}

// ParseToken parses token string with defined callback
func (j *JWT) ParseToken(token string, keyFunc KeyFunc) (*Claims, error) {
	if j.Claims == nil {
		return nil, ErrorInvalidCliam
	}

	return j.Claims.ParseToken(token, keyFunc)
}
