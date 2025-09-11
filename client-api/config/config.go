package config

import (
	"case-itau/utils/logger"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	APIPort      string
	RateLimitMax int64
	DBPath       string
}

func Load() *Config {
	godotenv.Load()

	checkEnvs()

	port := os.Getenv("API_PORT")
	if port == "" {
		port = "3000"
	}

	rateLimit := os.Getenv("RATE_LIMIT_MAX")
	if rateLimit == "" {
		rateLimit = "100"
	}
	rateLimitMax, err := strconv.Atoi(rateLimit)
	if err != nil {
		rateLimitMax = 100
	}

	dbPath := os.Getenv("DB_PATH")
	if dbPath == "" {
		dbPath = "database.db"
	}

	logger.NewLogger()

	return &Config{
		APIPort:      port,
		RateLimitMax: int64(rateLimitMax),
		DBPath:       dbPath,
	}
}

func checkEnvs() {
	requiredEnvs := []string{
		"API_PORT",
		"RATE_LIMIT_MAX",
		"DB_PATH",
	}

	for _, env := range requiredEnvs {
		if os.Getenv(env) == "" {
			panic("Environment variable " + env + " is not set")
		}
	}
}
