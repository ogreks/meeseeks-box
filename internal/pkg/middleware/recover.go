package middleware

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.uber.org/multierr"
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
				if err, ok := err.(error); ok {
					var fields []zap.Field = []zap.Field{
						zap.StackSkip("stack", config.StackSkip),
						zap.String("trace_id", c.Writer.Header().Get(DefaultTraceConfig.ResponseTraceKey)),
						zap.String("request", c.Request.RequestURI),
						zap.String("method", c.Request.Method),
						zap.String("client_ip", c.ClientIP()),
						zap.String("user_agent", c.Request.UserAgent()),
						zap.Error(err),
					}

					logger.Error("[Recovery] panic recovered", fields...)
				}

				fmt.Printf("%v\n", err)

				c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
					"code":    http.StatusInternalServerError,
					"message": "Internal server error, please try again later.",
					"data":    map[string]interface{}{},
				})
				return
			}

			if c.IsAborted() {
				var abortErr error
				for i := range c.Errors {
					multierr.AppendInto(&abortErr, c.Errors[i])
					var fields []zap.Field = []zap.Field{
						zap.StackSkip("stack", config.StackSkip),
						zap.String("trace_id", c.Writer.Header().Get(DefaultTraceConfig.ResponseTraceKey)),
						zap.String("request", c.Request.RequestURI),
						zap.String("method", c.Request.Method),
						zap.String("client_ip", c.ClientIP()),
						zap.String("user_agent", c.Request.UserAgent()),
						zap.Error(c.Errors[i]),
					}

					logger.Error("[gin aborted] errors:", fields...)
				}

				c.JSONP(c.Writer.Status(), gin.H{
					"code":    c.Writer.Status(),
					"message": c.Errors.Errors()[0],
					"data":    map[string]interface{}{},
				})
			}
		}()

		c.Next()
	}
}
