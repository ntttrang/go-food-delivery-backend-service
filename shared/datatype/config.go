package datatype

import (
	"os"
)

type Config struct {
	UserServiceURL string
	CatServiceURL  string
	EmailConfig    EmailConfig
	RedisConfig    RedisConfig
}

var config *Config

func NewConfig() *Config {
	const defaultPort = 587
	// portStr := os.Getenv("SMTP_PORT")
	// if portStr == "" {
	// 	fmt.Println("SMTP_PORT not set, using default")
	// 	portStr = fmt.Sprintf("%d", defaultPort)
	// }

	// port, err := strconv.Atoi(portStr)
	// if err != nil {
	// 	fmt.Printf("Invalid SMTP_PORT value '%s', using default: %d\n", portStr, defaultPort)
	// 	port = defaultPort
	// }

	port := defaultPort

	if config == nil {
		config = &Config{
			UserServiceURL: os.Getenv("USER_SERVICE_URL"),
			CatServiceURL:  os.Getenv("CAT_SERVICE_URL"),
			EmailConfig: EmailConfig{
				SMTPHost:     os.Getenv("SMTP_HOST"),
				SMTPPort:     port,
				SMTPUsername: os.Getenv("SMTP_USERNAME"),
				SMTPPassword: os.Getenv("SMTP_PASSWORD"),
			},
			RedisConfig: RedisConfig{
				Host:     os.Getenv("REDIS_ADDR"),
				Password: os.Getenv("REDIS_PASSWORD"),
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
