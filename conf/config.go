package conf

import "sync"

var (
	config *Config
	_confOnce sync.Once
)
type Config struct {
	DBConnectionString string `json:"db_connection_string"`
}

func Set(c *Config) {
	config = c
}

func Get() *Config{
	_confOnce.Do(func() {
		if config == nil {
			config = new(Config)
		}
	})
	return config
}