package middleware

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"github.com/ogreks/meeseeks-box/configs"
	"net/http"
	"strings"
	"time"
)

// ErrorInvalidKey is the error for invalid key
var (
	ErrorInvalidKey = errors.New("invalid key")
)

type JwtConfig struct {
	SigningMethod jwt.SigningMethod
	JWTHeaderKey  string
}

var defaultJwtConfig = &JwtConfig{
	SigningMethod: jwt.SigningMethodHS512,
	JWTHeaderKey:  configs.GetConfig().Jwt.HeaderKey,
}

// Claims is the interface for claims
type Claims interface {
	jwt.Claims
	// GetContent get content
	GetContent() any
}

type GlobalJWT struct {
	jwt.StandardClaims

	Content any `json:"content"`
}

func NewGlobalJWT(content any, expire time.Duration) *GlobalJWT {
	return &GlobalJWT{
		StandardClaims: jwt.StandardClaims{
			Issuer:    configs.GetConfig().Jwt.Issuer,
			ExpiresAt: time.Now().Add(expire).Unix(),
		},
		Content: content,
	}
}

// GetContent get content
func (m *GlobalJWT) GetContent() any {
	return m.Content
}

// CreateToken creates token
func (m *GlobalJWT) CreateToken() (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS512, m)
	return token.SignedString([]byte(configs.GetConfig().Jwt.Secret))
}

// Option is the option for jwt
type Option func(j *JwtMiddleware)

// KeyFunc is the key func for jwt
type KeyFunc = func() (any, error)

// JwtMiddleware is the jwt middleware
type JwtMiddleware struct {
	keyFunc KeyFunc
	claims  Claims
	cfg     *JwtConfig
}

// NewJWTMiddleware creates a new JWT with defined signingMethod, secret & claims
func NewJWTMiddleware(options ...Option) *JwtMiddleware {
	j := &JwtMiddleware{
		cfg: defaultJwtConfig,
	}

	for _, option := range options {
		option(j)
	}

	return j
}

// WithClaims sets claims
func WithClaims(claims Claims) Option {
	return func(j *JwtMiddleware) {
		j.claims = claims
	}
}

// WithKeyFunc sets keyFunc
func WithKeyFunc(keyFunc KeyFunc) Option {
	return func(j *JwtMiddleware) {
		j.keyFunc = keyFunc
	}
}

// WithSigningMethod sets signingMethod
func WithSigningMethod(signingMethod jwt.SigningMethod) Option {
	return func(j *JwtMiddleware) {
		j.cfg.SigningMethod = signingMethod
	}
}

func WithJWTHeaderKey(key string) Option {
	return func(j *JwtMiddleware) {
		j.cfg.JWTHeaderKey = key
	}
}

// ParseToken parses token string with defined callback
func (j *JwtMiddleware) ParseToken(token string) (Claims, error) {
	key, err := j.keyFunc()
	if err != nil {
		return nil, err
	}

	secret, err := jwt.ParseWithClaims(token, j.claims, func(token *jwt.Token) (interface{}, error) {
		return key, nil
	})
	if err != nil {
		return nil, err
	}

	return secret.Claims.(Claims), nil
}

// Builder builds jwt middleware
func (j *JwtMiddleware) Builder() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// get token from header
		tokenHeader := ctx.GetHeader(j.cfg.JWTHeaderKey)
		if tokenHeader == "" {
			_ = ctx.AbortWithError(http.StatusUnauthorized, errors.New("token is empty"))
			return
		}

		segs := strings.Split(tokenHeader, " ")
		if len(segs) < 2 {
			_ = ctx.AbortWithError(http.StatusUnauthorized, errors.New("token is invalid"))
			return
		}

		token := segs[1]
		// parse token
		claims, err := j.ParseToken(token)
		if err != nil {
			_ = ctx.AbortWithError(http.StatusUnauthorized, err)
			return
		}

		// set claims to context
		ctx.Set("claims", claims)
		ctx.Next()
	}
}
