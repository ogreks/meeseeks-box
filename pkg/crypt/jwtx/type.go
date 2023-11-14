package jwtx

import "github.com/golang-jwt/jwt"

type Token jwt.Token

type KeyFunc = func(*Token) (interface{}, error)

type Secret func() []byte

type SigningMethod jwt.SigningMethod
