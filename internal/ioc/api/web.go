package api

import "github.com/gin-gonic/gin"

func InitApiServer() *gin.Engine {
	return gin.Default()
}
