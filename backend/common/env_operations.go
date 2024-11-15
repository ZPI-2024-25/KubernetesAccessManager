package common

import (
	"log"
	"os"
)

func GetOrDefaultEnv(key, defaultValue string) string {
	val := os.Getenv(key)
	if val == "" {
		log.Printf("Environment variable %s not set, using default value %s", key, defaultValue)
		return defaultValue
	}
	return val
}