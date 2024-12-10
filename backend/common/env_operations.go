package common

import (
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

var (
	HealthPort             int
	AppPort                int
	KeycloakURL            string
	KeycloakClient         string
	KeycloakRealm          string
	KeycloakJwksUrl        string
	RoleMapNamespace       string
	RoleMapName            string
	USE_JWT_TOKEN_PATHS    bool
	TOKEN_ROLE_PATHS       string
	TOKEN_PATHS_SEP        string
	TOKEN_PATH_SEGMENT_SEP string
)

func InitEnv() {
	err := godotenv.Load()
	if err != nil {
		log.Println("Can't load .env file:", err)
	}
	var exists bool
	KeycloakJwksUrl, exists = os.LookupEnv("KEYCLOAK_JWKS_URL")
	if exists {
		log.Println("Using KEYCLOAK_JWKS_URL environment variable.")
	} else {
		log.Println("KEYCLOAK_JWKS_URL environment variable not set, setting default JWKS URL.")
		KeycloakURL = getEnvOrPanic("VITE_KEYCLOAK_URL")
		KeycloakRealm = getEnvOrPanic("VITE_KEYCLOAK_REALMNAME")
		KeycloakJwksUrl = fmt.Sprintf("%s/realms/%s/protocol/openid-connect/certs", KeycloakURL, KeycloakRealm)
		log.Printf("Using JWKS URL: %s\n", KeycloakJwksUrl)
	}
	usePaths := getEnvOrDefault("USE_JWT_TOKEN_PATHS", "false")
	if usePaths == "true" {
		USE_JWT_TOKEN_PATHS = true
		TOKEN_ROLE_PATHS = getEnvOrDefault("TOKEN_ROLE_PATHS", DEFAULT_TOKEN_ROLE_PATHS)
		TOKEN_PATHS_SEP = getEnvOrDefault("TOKEN_PATHS_SEP", DEFAULT_PATHS_SEP)
		TOKEN_PATH_SEGMENT_SEP = getEnvOrDefault("TOKEN_PATH_SEGMENT_SEP", DEFAULT_PATH_SEGMENT_SEP)
	}
	KeycloakClient = getEnvOrPanic("VITE_KEYCLOAK_CLIENTNAME")
	HealthPort = getEnvAsInt("HEALTH_PORT", 8082)
	AppPort = getEnvAsInt("BACKEND_PORT", 8080)
	RoleMapNamespace = getEnvOrDefault("ROLEMAP_NAMESPACE", DEFAULT_ROLEMAP_NAMESPACE)
	RoleMapName = getEnvOrDefault("ROLEMAP_NAME", DEFAULT_ROLEMAP_NAME)
}

func getEnvOrDefault(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists && value != "" {
		return value
	}
	log.Printf("Environment variable %s not set, using default value %s", key, defaultValue)
	return defaultValue
}

func getEnvOrPanic(key string) string {
	value, exists := os.LookupEnv(key)
	if !exists {
		log.Fatalf("Environment variable %s is not set. Exiting...", key)
	}
	return value
}

func getEnvAsInt(key string, defaultValue int) int {
	valueStr := getEnvOrDefault(key, strconv.Itoa(defaultValue))
	value, err := strconv.Atoi(valueStr)
	if err != nil {
		log.Fatalf("Invalid value for %s: %s. Must be an integer. Exiting...", key, valueStr)
	}
	return value
}
