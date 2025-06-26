package datatype

import (
	"os"
	"strconv"
)

type Config struct {
	EmailConfig  EmailConfig
	RedisConfig  RedisConfig
	GoogleConfig GoogleConfig
	Minio        MinIoConfig // Same as Amazon S3

	// URL for RPC
	UserServiceURL string
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
			Minio: MinIoConfig{
				AccessKey:  os.Getenv("MINIO_ACCESS_KEY"),
				BucketName: os.Getenv("MINIO_BUCKET_NAME"),
				Domain:     os.Getenv("MINIO_DOMAIN"),
				Region:     os.Getenv("MINIO_REGION"),
				SecretKey:  os.Getenv("MINIO_SECRET_KEY"),
				UseSSL:     false,
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
type MinIoConfig struct {
	AccessKey  string
	BucketName string
	Domain     string
	Region     string
	SecretKey  string
	UseSSL     bool
}
