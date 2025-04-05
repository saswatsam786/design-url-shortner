package config

import (
	"fmt"
	"os"
)

type Config struct {
	Port           string
	PostgresConfig PostgresConfig
	RedisConfig    RedisConfig
	BaseURL        string
}

type PostgresConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
}

type RedisConfig struct {
	Host string
	Port string
}

func LoadConfig() (*Config, error) {
	port := getEnvOrDefault("PORT", "8080")
	baseURL := getEnvOrDefault("BASE_URL", "http://localhost:8080")

	postgresConfig := PostgresConfig{
		Host:     getEnvOrDefault("POSTGRES_HOST", "localhost"),
		Port:     getEnvOrDefault("POSTGRES_PORT", "5432"),
		User:     getEnvOrDefault("POSTGRES_USER", "postgres"),
		Password: getEnvOrDefault("POSTGRES_PASSWORD", "postgres"),
		DBName:   getEnvOrDefault("POSTGRES_DB", "urlshortner"),
	}

	redisConfig := RedisConfig{
		Host: getEnvOrDefault("REDIS_HOST", "localhost"),
		Port: getEnvOrDefault("REDIS_PORT", "6379"),
	}

	return &Config{
		Port:           port,
		PostgresConfig: postgresConfig,
		RedisConfig:    redisConfig,
		BaseURL:        baseURL,
	}, nil
}

func (c *Config) GetPostgresDSN() string {
	return fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		c.PostgresConfig.Host,
		c.PostgresConfig.Port,
		c.PostgresConfig.User,
		c.PostgresConfig.Password,
		c.PostgresConfig.DBName,
	)
}

func (c *Config) GetRedisAddr() string {
	return fmt.Sprintf("%s:%s", c.RedisConfig.Host, c.RedisConfig.Port)
}

func getEnvOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}

	return defaultValue
}
