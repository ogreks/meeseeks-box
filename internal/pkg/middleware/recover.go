package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// RecoveryConfig defines the configs for Recovery middleware
type RecoveryConfig struct {
	StackAll  bool
	StackSize int
	StackSkip int
}

// DefaultRecoveryConfig is the default recovery configs
var DefaultRecoveryConfig = RecoveryConfig{
	StackAll:  false,
	StackSize: 1024 * 8,
	StackSkip: 3,
}

// Recovery returns a middleware for recovering from any panics and writes a 500 if there was one.
func Recovery(logger *zap.Logger) gin.HandlerFunc {
	return RecoveryWithConfig(logger, DefaultRecoveryConfig)
}

// RecoveryWithConfig returns a middleware for recovering from any panics and writes a 500 if there was one.
func RecoveryWithConfig(logger *zap.Logger, config RecoveryConfig) gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				var fields []zap.Field = []zap.Field{
					zap.StackSkip("stack", config.StackSkip),
					zap.String("request", c.Request.RequestURI),
					zap.String("method", c.Request.Method),
					zap.String("client_ip", c.ClientIP()),
					zap.String("user_agent", c.Request.UserAgent()),
					zap.Error(err.(error)),
				}

				logger.Error("[Recovery] panic recovered", fields...)

				c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
					"code":    http.StatusInternalServerError,
					"message": "Internal server error, please try again later.",
					"data":    map[string]interface{}{},
				})
			}
		}()

		c.Next()
	}
}
