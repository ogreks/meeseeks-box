package config

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
	})
}

func GetConfig() *Config {
	return cfg
}
