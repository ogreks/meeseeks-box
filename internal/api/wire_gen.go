// Code generated by Wire. DO NOT EDIT.

//go:generate go run github.com/google/wire/cmd/wire
//go:build !wireinject
// +build !wireinject

package api

import (
	"github.com/gin-gonic/gin"
	api2 "github.com/ogreks/meeseeks-box/internal/ioc/api"
)

// Injectors from wire.go:

func InitApiServer() *gin.Engine {
	engine := api2.InitApiServer()
	return engine
}
