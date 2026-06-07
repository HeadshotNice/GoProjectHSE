package config

import (
	"os"
	"strconv"
	"time"
)

type Config struct {
	HTTP     HTTPConfig
	DB       DBConfig
	Auth     AuthConfig
	RabbitMQ RabbitMQConfig
}

type HTTPConfig struct {
	Addr string
}

type DBConfig struct {
	Host     string
	Port     int
	User     string
	Password string
	Name     string
	SSLMode  string
}

type AuthConfig struct {
	JWTSecret string
	JWTIssuer string
	JWTTTL    time.Duration
}

type RabbitMQConfig struct {
	Enabled                bool
	URL                    string
	Exchange               string
	DocumentSubmittedKey   string
	DocumentStatusKey      string
	DocumentSubmittedQueue string
	DocumentStatusQueue    string
}

func FromEnv() Config {
	return Config{
		HTTP: HTTPConfig{
			Addr: env("HTTP_ADDR", ":8080"),
		},
		DB: DBConfig{
			Host:     env("DB_HOST", "localhost"),
			Port:     envInt("DB_PORT", 5432),
			User:     env("DB_USER", "postgres"),
			Password: env("DB_PASSWORD", "postgres"),
			Name:     env("DB_NAME", "postgres"),
			SSLMode:  env("DB_SSLMODE", "disable"),
		},
		Auth: AuthConfig{
			JWTSecret: env("JWT_SECRET", "dev-secret-change-me"),
			JWTIssuer: env("JWT_ISSUER", "hse-server"),
			JWTTTL:    envDuration("JWT_TTL", 24*time.Hour),
		},
		RabbitMQ: RabbitMQConfig{
			Enabled:                envBool("RABBITMQ_ENABLED", true),
			URL:                    env("RABBITMQ_URL", "amqp://admin:admin@localhost:5672/"),
			Exchange:               env("RABBITMQ_EXCHANGE", "main_exchange"),
			DocumentSubmittedKey:   env("RABBITMQ_DOCUMENT_SUBMITTED_KEY", "document-submitted"),
			DocumentStatusKey:      env("RABBITMQ_DOCUMENT_STATUS_KEY", "document-status-change"),
			DocumentSubmittedQueue: env("RABBITMQ_DOCUMENT_SUBMITTED_QUEUE", "document_submitted"),
			DocumentStatusQueue:    env("RABBITMQ_DOCUMENT_STATUS_QUEUE", "document_status"),
		},
	}
}

func env(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}

func envInt(key string, def int) int {
	if v := os.Getenv(key); v != "" {
		if n, err := strconv.Atoi(v); err == nil {
			return n
		}
	}
	return def
}

func envBool(key string, def bool) bool {
	if v := os.Getenv(key); v != "" {
		if b, err := strconv.ParseBool(v); err == nil {
			return b
		}
	}
	return def
}

func envDuration(key string, def time.Duration) time.Duration {
	if v := os.Getenv(key); v != "" {
		if d, err := time.ParseDuration(v); err == nil {
			return d
		}
	}
	return def
}
