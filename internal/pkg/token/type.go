package token

import "github.com/golang-jwt/jwt"

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
