package api

import (
	"github.com/gin-gonic/gin"
	"github.com/ogreks/meeseeks-box/internal/repository/orm"
)

func InitApiServer(db orm.Repo) *gin.Engine {
	return gin.Default()
}
