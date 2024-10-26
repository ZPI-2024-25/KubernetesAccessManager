package auth

import (
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"
	"github.com/MicahParks/keyfunc"
	"github.com/ZPI-2024-25/KubernetesAccessManager/models"
	"github.com/golang-jwt/jwt/v4"
	"github.com/joho/godotenv"
)

var jwks *keyfunc.JWKS

func init() {
	err := godotenv.Load()
	jwksURL := os.Getenv("KEYCLOAK_URL")
	jwks, err = keyfunc.Get(jwksURL, keyfunc.Options{
		RefreshInterval: time.Hour,
	})
	if err != nil {
		panic(fmt.Sprintf("Failed to create JWKS: %s", err))
	}
}

func GetJWTTokenFromHeader(r *http.Request) (string, error) {
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		return "", fmt.Errorf("authorization header missing")
	}

	parts := strings.Split(authHeader, " ")
	if len(parts) != 2 || parts[0] != "Bearer" {
		return "", fmt.Errorf("invalid authorization header format")
	}
	return parts[1], nil
}

func IsTokenValid(tokenStr string) bool {
	claims := jwt.MapClaims{}
	token, err := jwt.ParseWithClaims(tokenStr, &claims, jwks.Keyfunc)
	if err != nil {
		fmt.Printf("Token parsing error: %v\n", err)
		return false
	}
	if !token.Valid {
		return false
	}

	if exp, ok := claims["exp"].(float64); ok {
		expirationTime := time.Unix(int64(exp), 0)
		if time.Now().After(expirationTime) {
			fmt.Println("Token has expired")
			return false
		}
	} else {
		fmt.Println("Expiration claim missing or invalid")
		return false
	}

	return true
}

func IsUserAuthorized(operation models.Operation, roles []string) (bool, error) {
	authService, err := GetInstance()
	if err != nil {
		fmt.Printf("Error when loading config: %v\n", err)
		return false, err
	}

	for _, role := range roles {
		if authService.HasPermission(role, &operation) {
			return true, nil
		}
	}
	return false, nil
}
