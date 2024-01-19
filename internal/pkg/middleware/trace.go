package middleware

import (
	"crypto/rand"
	"fmt"
	"io"

	"github.com/gin-gonic/gin"
)

type TraceConfig struct {
	SkippedPathPrefixes []string
	RequestHeaderKey    string
	ResponseTraceKey    string
}

var DefaultTraceConfig = TraceConfig{
	SkippedPathPrefixes: []string{"/health"},
	RequestHeaderKey:    "X-Request-Id",
	ResponseTraceKey:    "X-Trace-Id",
}

func Trace() gin.HandlerFunc {
	return TraceWithConfig(DefaultTraceConfig)
}

func TraceWithConfig(config TraceConfig) gin.HandlerFunc {
	return func(c *gin.Context) {
		if SkippedPathPrefixes(c, config.SkippedPathPrefixes...) {
			c.Next()
			return
		}

		traceID := c.GetHeader(config.RequestHeaderKey)
		if traceID == "" {
			buf := make([]byte, 16)
			_, err := io.ReadFull(rand.Reader, buf)
			if err != nil {
				return
			}
			traceID = fmt.Sprintf("%x", buf)
		}

		c.Writer.Header().Set(config.ResponseTraceKey, traceID)
		c.Next()
	}
}
