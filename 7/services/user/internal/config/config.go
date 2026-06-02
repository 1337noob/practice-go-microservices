package config

import (
	"os"
	"strconv"
)

type Config struct {
	GRPCPort    int
	KafkaBroker string
}

func Load() Config {
	return Config{
		GRPCPort:    getEnvInt("GRPC_PORT", 50051),
		KafkaBroker: getEnv("KAFKA_BROKER", "localhost:9092"),
	}
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

func getEnvInt(key string, fallback int) int {
	if v := os.Getenv(key); v != "" {
		if i, err := strconv.Atoi(v); err == nil {
			return i
		}
	}
	return fallback
}
