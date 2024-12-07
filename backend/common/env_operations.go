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

func GetEnvOrPanic(envVar string) string {
	value, exists := os.LookupEnv(envVar)
	if !exists {
		log.Fatalf("Environment variable %s is not set. Exiting...", envVar)
	}
	return value
}
