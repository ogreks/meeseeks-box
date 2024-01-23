package middleware

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func Cors() gin.HandlerFunc {
	return CorsWithConfig(cors.DefaultConfig())
}

func CorsWithConfig(cc cors.Config) gin.HandlerFunc {
	return cors.New(cc)
}
