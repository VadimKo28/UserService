package config

import (
	"log"
	"os"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
	"github.com/joho/godotenv"
)

type Config struct {
	Env        string      `yaml:"env" env-required:"true"`
	DatabasePath string    `yaml:"database_path" env-required:"true"`
	HTTPServer HTTPServer  `yaml:"http_server"`
}

type HTTPServer struct {
	Host        string `yaml:"host" env-default:"localhost"`
	Port        string `yaml:"port" env-default:"8081"`
	Timeout     time.Duration    `yaml:"timeout" env-default:"4"`
	IdleTimeout time.Duration    `yaml:"idle_timeout" env-default:"60"`
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