package configs

import (
	"bytes"
	"os"
	"path/filepath"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
)

var cfg = new(Config)

type Config struct {
	Server   Server   `mapstructure:"server"`
	Database Database `mapstructure:"database"`
	Log      Log      `mapstructure:"log"`
	Jwt      Jwt      `mapstructure:"jwt"`
}

func (c *Config) GetServer() Server {
	return c.Server
}

func (c *Config) GetDatabase() Database {
	return c.Database
}

func (c *Config) GetLog() Log {
	return c.Log
}

// initDefault 初始化默认配置
func initDefaultConfig() {
	if cfg.Database.LogPath == "" {
		cfg.Database.LogPath = "./log/sql.log"
	}

	if cfg.Jwt.Issuer == "" {
		cfg.Jwt.Issuer = "meeseeks-box"
	}

	if cfg.Jwt.HeaderKey == "" {
		cfg.Jwt.HeaderKey = RequestHeaderJWTKey
	}

	if cfg.Jwt.Expire == 0 {
		cfg.Jwt.Expire = RequestHeaderJWTExpireTime
	}
}

func InitConfig(configFile string) *Config {
	f, err := os.ReadFile(configFile)
	if err != nil {
		panic(err)
	}

	viper.SetConfigType("yaml")
	err = viper.ReadConfig(bytes.NewReader(f))
	if err != nil {
		panic(err)
	}

	if err := viper.Unmarshal(cfg); err != nil {
		panic(err)
	}

	initDefaultConfig()

	cfg.Watch(configFile)

	return cfg
}

func (c *Config) Watch(configurePath string) {
	viper.SetConfigFile(configurePath)
	if _, err := os.Stat(configurePath); err != nil {
		if err := os.MkdirAll(filepath.Dir(configurePath), 0655); err != nil {
			panic(err)
		}

		f, err := os.Create(configurePath)
		if err != nil {
			panic(err)
		}

		defer func() {
			_ = f.Close()
		}()

		if err := viper.WriteConfig(); err != nil {
			panic(err)
		}
	}

	viper.WatchConfig()
	viper.OnConfigChange(func(e fsnotify.Event) {
		if err := viper.Unmarshal(cfg); err != nil {
			panic(err)
		}
		initDefaultConfig()
	})
}

func GetConfig() *Config {
	return cfg
}
