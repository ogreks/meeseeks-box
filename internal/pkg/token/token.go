package token

import (
	"context"
	"errors"
	"time"

	"github.com/golang-jwt/jwt"
)

var (
	ErrTokenNotFound = errors.New("token not found")
	ErrInvalidKey    = errors.New("key is invalid")
	ErrTokenTimeout  = errors.New("token is expired")
)

//go:generate mockgen -source=./token.go -package=tkmocks -destination=mocks/token.mock.go Token
type Token[T Type, V Val] interface {
	// CreateToken creates a new token.
	CreateToken(ctx context.Context, claim jwt.Claims, expire ...time.Duration) (T, error)

	// RefreshToken creates a new refresh token.
	// The old token is revoked or set token to Expire.
	RefreshToken(ctx context.Context, token T, claim jwt.Claims, expire ...time.Duration) (T, error)

	// Validate validates a token.
	Validate(token T) (jwt.Claims, error)

	// Store returns the store used by the token.
	Store() Store[T]
}

type Option[T Type, V Val] func(Token[T, V])

type DefaultToken[T string, F Fun] struct {
	store Store[T]

	f F

	Expire time.Duration

	claims jwt.Claims
}

func WithStore[T string, F Fun](store Store[T]) Option[T, F] {
	return func(t Token[T, F]) {
		t.(*DefaultToken[T, F]).store = store
	}
}

func WithExpire[T string, F Fun](expire time.Duration) Option[T, F] {
	return func(t Token[T, F]) {
		t.(*DefaultToken[T, F]).Expire = expire
	}
}

func WithFun[T string, F Fun](f F) Option[T, F] {
	return func(t Token[T, F]) {
		t.(*DefaultToken[T, F]).f = f
	}
}

func WithClaims[T string, F Fun](claims jwt.Claims) Option[T, F] {
	return func(t Token[T, F]) {
		t.(*DefaultToken[T, F]).claims = claims
	}
}

// CreateToken creates a new token.
// jwt claim token to string
// `f` is a function that returns jwt.SigningMethod, []byte, jwt.Claims
// not `f` is struct dt.f
func (dt *DefaultToken[T, F]) CreateToken(ctx context.Context, claim jwt.Claims, expire ...time.Duration) (T, error) {
	if len(expire) <= 0 {
		expire = []time.Duration{dt.Expire}
	}
	method, secret := dt.f()
	token := jwt.NewWithClaims(method, claim)

	signedString, err := token.SignedString(secret)
	if err != nil {
		return "", err
	}

	err = dt.store.Set(T(signedString), expire[0])
	if err != nil {
		return "", err
	}

	return T(signedString), nil
}

// RefreshToken creates a new refresh token.
// The old token is revoked or set token to Expire.
// `token` is the old token.
func (dt *DefaultToken[T, F]) RefreshToken(ctx context.Context, token T, claim jwt.Claims, expire ...time.Duration) (T, error) {
	_, err := dt.Validate(token)
	if err != nil {
		if !errors.Is(err, ErrTokenTimeout) && !errors.Is(err, ErrTokenNotFound) {
			return "", err
		}
	}

	err = dt.Store().Delete(token)
	if err != nil {
		return "", err
	}

	return dt.CreateToken(ctx, claim, expire...)
}

// Validate validates a token.
// `f` is a function that returns jwt.SigningMethod, []byte, jwt.Claims
func (dt *DefaultToken[T, F]) Validate(token T) (jwt.Claims, error) {
	if !dt.store.Exists(token) {
		return nil, ErrTokenNotFound
	}

	_, secret := dt.f()

	t, err := jwt.ParseWithClaims(string(token), dt.claims, func(token *jwt.Token) (interface{}, error) {
		return secret, nil
	})

	if errors.Is(err, jwt.ErrInvalidKey) {
		return nil, ErrInvalidKey
	}

	if err := t.Claims.Valid(); err != nil {
		var jvl *jwt.ValidationError
		if errors.As(err, &jvl) {
			switch jvl.Errors {
			case jwt.ValidationErrorExpired:
				return nil, ErrTokenTimeout
			default:
				return nil, err
			}
		}

		return nil, err
	}

	return t.Claims, err
}

func (dt *DefaultToken[T, F]) Store() Store[T] {
	return dt.store
}

func NewDefaultToken[T string, F Fun](opts ...Option[T, F]) Token[T, F] {
	t := &DefaultToken[T, F]{}

	for _, opt := range opts {
		opt(t)
	}

	return t
}
