package config

type Config struct {
	server   Server   `yaml:"server"`
	database Database `yaml:"database"`
	log      Log      `yaml:"log"`
}

func (c *Config) GetServer() Server {
	return c.server
}

func (c *Config) GetDatabase() Database {
	return c.database
}

func (c *Config) GetLog() Log {
	return c.log
}
