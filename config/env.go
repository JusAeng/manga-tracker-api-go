package config

import (
	"os"
	"strings"
	"sync"

	"github.com/joho/godotenv"
)

var loadDotEnvOnce sync.Once

// loadDotEnv loads a .env file for local development if one is present.
// In deployed environments (Cloud Run, Docker with injected env vars) there
// is no .env file, and that is expected, not an error, so we don't fail here.
func loadDotEnv() {
	loadDotEnvOnce.Do(func() {
		_ = godotenv.Load()
	})
}

func GetEnv(key string) (string, error) {
	loadDotEnv()

	return os.Getenv(strings.ToUpper(key)), nil
}