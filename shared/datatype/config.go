package datatype

import (
	"os"
	"strconv"
)

type Config struct {
	UserServiceURL       string
	CatServiceURL        string
	FoodServiceURL       string
	RestaurantServiceURL string
	CartServiceURL       string
	PaymentServiceURL    string
	EmailConfig          EmailConfig
	RedisConfig          RedisConfig
	GoogleConfig         GoogleConfig
	MinioS3              MinioS3Config
	ElasticSearch        ElasticSearchConfig
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
			UserServiceURL:       os.Getenv("USER_SERVICE_URL"),
			CatServiceURL:        os.Getenv("CAT_SERVICE_URL"),
			FoodServiceURL:       os.Getenv("FOOD_SERVICE_URL"),
			RestaurantServiceURL: os.Getenv("RESTAURANT_SERVICE_URL"),
			CartServiceURL:       os.Getenv("CART_SERVICE_URL"),
			PaymentServiceURL:    os.Getenv("PAYMENT_SERVICE_URL"),
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
			MinioS3: MinioS3Config{
				AccessKey:  os.Getenv("MINIO_ACCESS_KEY"),
				BucketName: os.Getenv("MINIO_BUCKET_NAME"),
				Domain:     os.Getenv("MINIO_DOMAIN"),
				Region:     os.Getenv("MINIO_REGION"),
				SecretKey:  os.Getenv("MINIO_SECRET_KEY"),
				UseSSL:     false,
			},
			ElasticSearch: ElasticSearchConfig{
				Addresses: []string{os.Getenv("ES_ADDRESS")},
				Username:  os.Getenv("ES_USERNAME"),
				Password:  os.Getenv("ES_PASSWORD"),
				CloudID:   os.Getenv("ES_CLOUD_ID"),
				APIKey:    os.Getenv("ES_API_KEY"),
				IndexName: os.Getenv("ES_INDEX_NAME"),
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

type MinioS3Config struct {
	AccessKey  string
	BucketName string
	Domain     string
	Region     string
	SecretKey  string
	UseSSL     bool
}
