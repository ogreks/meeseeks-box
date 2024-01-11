package token

import (
	"errors"

	"github.com/golang-jwt/jwt"
)

var (
	ErrNotString = errors.New("type is not string")
)

type (
	Type interface {
		string | ~int | ~uint | float64 | float32
	}

	Val interface {
		interface{}
	}

	Fun interface {
		func() (jwt.SigningMethod, []byte, jwt.Claims)
	}
)
