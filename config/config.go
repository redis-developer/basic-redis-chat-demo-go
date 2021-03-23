package config

import (
	_ "github.com/joho/godotenv/autoload"
	"os"
)

const (
	envNameServerAddress  = "SERVER_ADDRESS"
	envNameClientLocation = "CLIENT_LOCATION"
	envNameRedisAddress   = "REDIS_ADDRESS"
	envNameRedisPassword  = "REDIS_PASSWORD"

	defaultRedisAddress   = "localhost:6379"
	defaultServerAddress  = ":40080"
	defaultClientLocation = "/usr/local/share/dinamicka/public"
)

type Config struct {
	ServerAddress  string
	ClientLocation string
	RedisAddress   string
	RedisPassword  string
}

func NewConfig() *Config {
	config := &Config{
		ServerAddress:  os.Getenv(envNameServerAddress),
		ClientLocation: os.Getenv(envNameClientLocation),
		RedisAddress:   os.Getenv(envNameRedisAddress),
		RedisPassword:  os.Getenv(envNameRedisPassword),
	}
	if config.ServerAddress == "" {
		config.ServerAddress = defaultServerAddress
	}
	if config.ClientLocation == "" {
		config.ClientLocation = defaultClientLocation
	}
	if config.RedisAddress == "" {
		config.RedisAddress = defaultRedisAddress
	}

	return config
}
