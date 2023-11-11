package ioc

import "github.com/ogreks/meeseeks-box/config"

func InitConfig() config.Config {
	return *config.GetConfig()
}
