package datatype

import (
	"os"
	"strconv"
)

type Config struct {
	UserServiceURL string
	CatServiceURL  string
	FoodServiceURL string
	EmailConfig    EmailConfig
	RedisConfig    RedisConfig
	GoogleConfig   GoogleConfig
}

var config *Config

func NewConfig() *Config {
	var smtpPort int
	portStr := os.Getenv("SMTP_PORT")
	if portStr != "" {
		smtpPort, _ = strconv.Atoi(portStr)
	}

	if config == nil {
		config = &Config{
			UserServiceURL: os.Getenv("USER_SERVICE_URL"),
			CatServiceURL:  os.Getenv("CAT_SERVICE_URL"),
			FoodServiceURL: os.Getenv("FOOD_SERVICE_URL"),
			EmailConfig: EmailConfig{
				SMTPHost:     os.Getenv("SMTP_HOST"),
				SMTPPort:     smtpPort,
				SMTPUsername: os.Getenv("SMTP_USERNAME"),
				SMTPPassword: os.Getenv("SMTP_PASSWORD"),
			},
			RedisConfig: RedisConfig{
				Host:     os.Getenv("REDIS_ADDR"),
				Password: os.Getenv("REDIS_PASSWORD"),
			},
			GoogleConfig: GoogleConfig{
				ClientId:     os.Getenv("GG_CLIENT_ID"),
				ClientSecret: os.Getenv("GG_CLIENT_SECRET"),
				RedirectUrl:  os.Getenv("GG_REDIRECT_URL"),
			},
		}
	}
	return config
}

func GetConfig() *Config {
	return NewConfig()
}

type EmailConfig struct {
	SMTPHost     string
	SMTPPort     int
	SMTPUsername string
	SMTPPassword string
}

type RedisConfig struct {
	Host     string
	Password string
}

type GoogleConfig struct {
	ClientId     string
	ClientSecret string
	RedirectUrl  string
}
