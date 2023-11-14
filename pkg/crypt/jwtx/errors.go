package jwtx

import (
	"errors"

	"github.com/golang-jwt/jwt"
)

var (
	ErrInvalidKey     = jwt.ErrInvalidKey
	ErrorInvalidToken = errors.New("invalid token")
	ErrorInvalidCliam = errors.New("invalid claim")
)
