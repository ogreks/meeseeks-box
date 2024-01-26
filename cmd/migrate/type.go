package migrate

import (
	"fmt"

	"github.com/ogreks/meeseeks-box/internal/model"
	"github.com/ogreks/meeseeks-box/pkg/command"
)

var (
	defaultConfigFile = fmt.Sprintf("%s/config.yml", command.HelpGetWorkDir())
	daoPath           = "internal/dao"
	modelConfigs      = map[string][]any{
		//"system":       systemModels(),
		"user":         userModels(),
		"account":      accountModels(),
		"session_keys": sessionKeysModels(),
	}
)

//func systemModels() []any {
//	return []any{
//		&model.Config{},
//		&model.VerifyCode{},
//	}
//}

func userModels() []any {
	return []any{
		&model.User{},
	}
}

func accountModels() []any {
	return []any{
		&model.Account{},
		&model.AccountConnect{},
	}
}

func sessionKeysModels() []any {
	return []any{
		&model.SessionKey{},
	}
}
