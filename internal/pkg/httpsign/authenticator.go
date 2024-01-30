package httpsign

import (
	"bytes"
	"context"
	"encoding/base64"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/ogreks/meeseeks-box/internal/pkg/httpsign/validator"
	"net/http"
	"strings"
)

const (
	requestTarget = "(request-target)"
	date          = "date"
	digest        = "digest"
	host          = "host"
)

var defaultRequiredHeaders = []string{requestTarget, date, digest}

type Authenticator struct {
	secrets    *Secrets
	validators []validator.Validator
	headers    []string
	encryptKey EncryptionKey // TODO 需要实时变更
}

type Option func(*Authenticator)

func WithValidator(validators ...validator.Validator) Option {
	return func(a *Authenticator) {
		a.validators = validators
	}
}

func WithRequiredHeaders(headers []string) Option {
	return func(a *Authenticator) {
		a.headers = headers
	}
}

func WithEncryptKey(e EncryptionKey) Option {
	return func(a *Authenticator) {
		a.encryptKey = e
	}
}

func NewAuthenticator(options ...Option) *Authenticator {
	a := &Authenticator{secrets: NewSecrets()}

	for _, option := range options {
		option(a)
	}

	if a.validators == nil {
		a.validators = []validator.Validator{
			validator.NewDateValidator(),
		}
	}

	if len(a.headers) == 0 {
		a.headers = defaultRequiredHeaders
	}

	return a
}

func (a *Authenticator) isValidHeader(headers []string) bool {
	m := len(headers)
	for _, h := range a.headers {
		i := 0
		for i := 0; i < m; i++ {
			if h == headers[i] {
				break
			}
		}
		if i == m {
			return false
		}
	}

	return true
}

func (a *Authenticator) getSecret(ctx context.Context, keyID KeyID, algorithm string) (*Secret, error) {
	secret, ok := a.secrets.container[keyID]
	if !ok {
		var err error
		secret, err = a.encryptKey.Secret(ctx, keyID)
		if err != nil {
			return nil, ErrInvalidKeyID
		}
	}

	if secret.Algorithm.Name() != algorithm {
		if algorithm != "" {
			return nil, ErrIncorrectAlgorithm
		}
	}

	return secret, nil
}

func constructSignMessage(r *http.Request, headers []string) string {
	var signBuffer bytes.Buffer
	for i, field := range headers {
		var fieldValue string
		switch field {
		case host:
			fieldValue = r.Host
		case requestTarget:
			fieldValue = fmt.Sprintf("%s %s", strings.ToLower(r.Method), r.URL.RequestURI())
		default:
			fieldValue = r.Header.Get(field)
		}
		signString := fmt.Sprintf("%s: %s", field, fieldValue)
		signBuffer.WriteString(signString)
		if i < len(headers)-1 {
			signBuffer.WriteString("\n")
		}
	}
	return signBuffer.String()
}

func (a *Authenticator) Authenticated() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		sigHeader, err := NewSignatureHeader(ctx.Request)
		if err != nil {
			_ = ctx.AbortWithError(http.StatusUnauthorized, err)
			return
		}

		for _, v := range a.validators {
			if err := v.Validator(ctx.Request); err != nil {
				_ = ctx.AbortWithError(http.StatusBadRequest, err)
				return
			}
		}

		if !a.isValidHeader(sigHeader.headers) {
			_ = ctx.AbortWithError(http.StatusBadRequest, ErrHeaderNotEnough)
			return
		}

		secret, err := a.getSecret(ctx.Request.Context(), sigHeader.keyID, sigHeader.algorithm)
		if err != nil {
			if errors.Is(err, ErrInvalidKeyID) {
				_ = ctx.AbortWithError(http.StatusUnauthorized, err)
				return
			}
			_ = ctx.AbortWithError(http.StatusBadRequest, err)
			return
		}

		signString := constructSignMessage(ctx.Request, sigHeader.headers)
		signature, err := secret.Algorithm.Sign(signString, secret.Key)
		if err != nil {
			_ = ctx.AbortWithError(http.StatusInternalServerError, err)
			return
		}

		signatureBase64 := base64.StdEncoding.EncodeToString(signature)
		if signatureBase64 != sigHeader.signature {
			_ = ctx.AbortWithError(http.StatusUnauthorized, ErrInvalidSign)
			return
		}

		ctx.Next()
	}
}
