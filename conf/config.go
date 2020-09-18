package conf

import "sync"

var (
	config *Config
	_confOnce sync.Once
)
type Config struct {
	DBConnectionString string `json:"db_connection_string"`
	LogPath string `json:"log_path"`
	UploadPath string `json:"upload_path"`
	RedisConnectionString string `json:"redis_connection_string"`
	AllowOrigin string `json:"allow_origin"`
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