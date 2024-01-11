package configs

import (
	"bytes"
	"os"
	"path/filepath"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
)

var cfg = NewConfig()

type Config struct {
	// TODO watch config file change broadcast notification
	watch    chan Config
	Server   Server   `mapstructure:"server"`
	RCache   RCache   `mapstructure:"rcache"`
	Database Database `mapstructure:"database"`
	Log      Log      `mapstructure:"log"`
	Jwt      Jwt      `mapstructure:"jwt"`
	WebHook  struct {
		Feishu Feishu `mapstructure:"feishu"`
	} `mapstructure:"webhook"`
}

// init config default value
func NewConfig() *Config {
	return &Config{
		watch: make(chan Config, 1),
		Server: Server{
			Debug:       true,
			Addr:        "0.0.0.0",
			Port:        80,
			MaxConn:     1000,
			ReadTimeout: 60,
			LogPath:     "./log/server.log",
		},
	}
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
