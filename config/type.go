package config

type Log struct {
	// 日志级别
	Level string `yaml:"level"`
	// 日志格式
	Format string `yaml:"format"`
	// 日志文件
	File string `yaml:"file"`
	// 日志文件最大大小
	MaxSize int `yaml:"max_size"`
	// 日志文件最大保留天数
}

type Database struct {
	// 数据库驱动
	Driver string `yaml:"driver"`
	// 数据库连接
	Source string `yaml:"source"`
	// 数据库最大空闲连接数
	MaxIdleConns int `yaml:"max_idle_conns"`
	// 数据库最大打开连接数
	MaxOpenConns int `yaml:"max_open_conns"`
}

type Server struct {
	// 服务器地址
	Addr string `yaml:"addr"`
	// 服务器端口
	Port int `yaml:"port"`
	// 服务器最大并发连接数
	MaxConn int `yaml:"max_conn"`
	// 服务器读超时
	ReadTimeout int `yaml:"read_timeout"`
}
