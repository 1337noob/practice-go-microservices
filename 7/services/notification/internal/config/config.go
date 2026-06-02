package config

import (
	"os"
)

type Config struct {
	KafkaBroker    string
	ConsumerGroup  string
	SMTPHost       string
	SMTPPort       string
	SMTPFrom       string
}

func Load() Config {
	return Config{
		KafkaBroker:   getEnv("KAFKA_BROKER", "localhost:9092"),
		ConsumerGroup: getEnv("KAFKA_CONSUMER_GROUP", "notification-group"),
		SMTPHost:      getEnv("SMTP_HOST", "localhost"),
		SMTPPort:      getEnv("SMTP_PORT", "1025"),
		SMTPFrom:      getEnv("SMTP_FROM", "noreply@example.com"),
	}
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
