package ioc

import "github.com/ogreks/meeseeks-box/configs"

func InitConfig() configs.Config {
	return *configs.GetConfig()
}
