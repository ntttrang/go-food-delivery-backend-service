package datatype

import "os"

type Config struct {
	UserServiceURL string
	CatServiceURL  string
}

var config *Config

func NewConfig() *Config {
	if config == nil {
		config = &Config{
			UserServiceURL: os.Getenv("USER_SERVICE_URL"),
			CatServiceURL:  os.Getenv("CAT_SERVICE_URL"),
		}
	}
	return config
}

func GetConfig() *Config {
	return NewConfig()
}
