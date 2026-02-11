package config

import (
	"log"
	"os"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
	"github.com/joho/godotenv"
)

type Config struct {
	Env          string     `yaml:"env" env-required:"true"`
	DatabasePath string     `yaml:"database_path" env:"DATABASE_URL" env-required:"true"`
	HTTPServer   HTTPServer `yaml:"http_server"`
	Auth         Auth       `yaml:"auth"`
	Kafka        Kafka      `yaml:"kafka"`
}

type Auth struct {
	TokenTTL        time.Duration `yaml:"token_ttl" env-default:"15"`
	RefreshTokenTTL time.Duration `yaml:"refresh_token_ttl" env-default:"720h"`
	JWTSecret       string        `yaml:"jwt_secret" env:"JWT_SECRET" env-required:"true"`
	PasswordSalt    string        `yaml:"password_salt" env:"SALT_HASH" env-required:"true"`
}

type HTTPServer struct {
	Host        string        `yaml:"host" env-default:"localhost"`
	Port        string        `yaml:"port" env-default:"8081"`
	Timeout     time.Duration `yaml:"timeout" env-default:"4"`
	IdleTimeout time.Duration `yaml:"idle_timeout" env-default:"60"`
}

type Kafka struct {
	Brokers           []string `yaml:"brokers" env:"KAFKA_BROKERS" env-separator:","`
	SubscriptionTopic string   `yaml:"subscription_topic" env:"KAFKA_SUBSCRIPTION_TOPIC" env-default:"subscription-events"`
}

func MustLoad() Config {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	config_path := os.Getenv("CONFIG_PATH")

	if config_path == "" {
		log.Fatal("CONFIG_PATH environment variable is not set")
	}

	if _, err := os.Stat(config_path); os.IsNotExist(err) {
		log.Fatalf("Config file does not exist at path: %s", config_path)
	}

	var cfg Config

	if err := cleanenv.ReadConfig(config_path, &cfg); err != nil {
		log.Fatalf("Failed to read config: %v", err)
	}

	return cfg
}
