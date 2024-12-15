package common

import (
	"fmt"
	"github.com/joho/godotenv"
	"log"
	"os"
	"strconv"
)

var (
	HealthPort       int
	AppPort          int
	KeycloakURL      string
	KeycloakClient   string
	KeycloakRealm    string
	KeycloakJwksUrl  string
	RoleMapNamespace string
	RoleMapName      string
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
		KeycloakRealm = getEnvOrPanic("VITE_KEYCLOAK_REALM_NAME")
		KeycloakJwksUrl = fmt.Sprintf("%s/realms/%s/protocol/openid-connect/certs", KeycloakURL, KeycloakRealm)
		log.Printf("Using JWKS URL: %s\n", KeycloakJwksUrl)
	}
	KeycloakClient = getEnvOrPanic("VITE_KEYCLOAK_CLIENT_NAME")
	log.Printf("Using Keycloak client: %s\n", KeycloakClient)
	HealthPort = getEnvAsInt("HEALTH_PORT", 8082)
	log.Printf("Using health port: %d\n", HealthPort)
	AppPort = getEnvAsInt("BACKEND_PORT", 8080)
	log.Printf("Using application port: %d\n", AppPort)
	RoleMapNamespace = getEnvOrDefault("ROLEMAP_NAMESPACE", DEFAULT_ROLEMAP_NAMESPACE)
	log.Printf("Using role map namespace: %s\n", RoleMapNamespace)
	RoleMapName = getEnvOrDefault("ROLEMAP_NAME", DEFAULT_ROLEMAP_NAME)
	log.Printf("Using role map name: %s\n", RoleMapName)
}

func getEnvOrDefault(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
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
