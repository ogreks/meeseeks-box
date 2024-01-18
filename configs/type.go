package configs

const (
	Debug int8 = iota - 1
	Info
	Warn
)

type Log struct {
	// 日志级别
	Level int `mapstructure:"level"`
	// 日志格式
	Format string `mapstructure:"format"`
	// 日志文件
	File string `mapstructure:"file"`
	// 日志文件最大大小
	MaxSize int `mapstructure:"max_size"`
	// 日志文件最大保留天数
}

type Database struct {
	// 数据库驱动
	Driver string `mapstructure:"driver"`
	// 数据库连接
	Source string `mapstructure:"source"`
	// 数据库最大空闲连接数
	// 数据库最大打开连接数
	MaxOpenConn int `mapstructure:"max_open_conn"`

	MaxIdleConn int `mapstructure:"max_idle_conn"`

	MaxLifetime int `mapstructure:"max_life_time"`

	Mode int `mapstructure:"mode"`

	Charset string `mapstructure:"charset"`

	LogPath string `mapstructure:"log_path"`
}

type RCache struct {
	// driver localhost:6379
	Driver string `mapstructure:"driver"`
	// password auth redis password
	Password string `mapstructure:"password"`
}

type Server struct {
	// server debug mode (default: true)
	Debug bool `mapstructure:"debug"`
	// server address to listen (default: 0.0.0.0)
	Addr string `mapstructure:"addr"`
	// server port (default: 80)
	Port int `mapstructure:"port"`
	// server max connection (default: 1000)
	MaxConn int `mapstructure:"max_conn"`
	// server read timeout (default: 60)
	ReadTimeout int `mapstructure:"read_timeout"`
	// server log path (default: ./log/server.log)
	LogPath string `mapstructure:"log_path"`
}

type Jwt struct {
	Secret    string `mapstructure:"secret"`
	Expire    int    `mapstructure:"expire"`
	Issuer    string `mapstructure:"issuer"`
	HeaderKey string `mapstructure:"header_key"`
	RefersKey string `mapstructure:"refers_key"`
}

type Feishu struct {
	AppId             string `mapstructure:"app_id"`
	AppSecret         string `mapstructure:"secret"`
	EncryptKey        string `mapstructure:"encrypt_key"`
	VerificationToken string `mapstructure:"verification_token"`
}
