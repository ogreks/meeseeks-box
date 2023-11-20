package configs

type Log struct {
	// 日志级别
	Level string `mapstructure:"level"`
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

type Server struct {
	Debug bool `mapstructure:"debug"`
	// 服务器地址
	Addr string `mapstructure:"addr"`
	// 服务器端口
	Port int `mapstructure:"port"`
	// 服务器最大并发连接数
	MaxConn int `mapstructure:"max_conn"`
	// 服务器读超时
	ReadTimeout int `mapstructure:"read_timeout"`
}

type Jwt struct {
	Secret    string `mapstructure:"secret"`
	Expire    int    `mapstructure:"expire"`
	Issuer    string `mapstructure:"-"`
	HeaderKey string `mapstructure:"-"`
}
