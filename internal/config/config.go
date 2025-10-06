package config

import (
	"os"
	"strconv"
)

const (
	workers_count         = "WORKERS"
	default_workers_count = 4
	queue_size            = "QUEUE_SIZE"
	default_queue_size    = 64
)

type Config struct {
	Workers   int
	QueueSize int
}

func Load() Config {
	return Config{
		Workers:   getValFromEnv(workers_count, default_workers_count),
		QueueSize: getValFromEnv(queue_size, default_queue_size),
	}
}

func getValFromEnv(envKey string, def int) int {
	sVal := os.Getenv(envKey)
	if sVal == "" {
		return def
	}

	val, err := strconv.Atoi(sVal)
	if err != nil || val <= 0 {
		return def
	}
	return val
}
