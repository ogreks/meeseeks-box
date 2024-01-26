package auth

import (
	"github.com/gin-gonic/gin"
	"github.com/ogreks/meeseeks-box/internal/pkg/httpsign"
	"github.com/ogreks/meeseeks-box/internal/service/platform"
)

type OpenMiddleware struct {
	s *httpsign.Authenticator
}

func NewOpenMiddleware(s platform.Service) *OpenMiddleware {
	return &OpenMiddleware{
		s: httpsign.NewAuthenticator(
			httpsign.WithEncryptKey(s),
		),
	}
}

func (o *OpenMiddleware) Builder() gin.HandlerFunc {
	return o.s.Authenticated()
}
