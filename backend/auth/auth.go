package auth

import (
	"errors"
	"fmt"
	"log"
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
	if jwksURL == "" {
		log.Println("KEYCLOAK_URL environment variable not set")
		return
	}
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

func IsTokenValid(tokenStr string) (bool, *jwt.MapClaims) {
	claims := jwt.MapClaims{}
	token, err := jwt.ParseWithClaims(tokenStr, &claims, jwks.Keyfunc)
	if err != nil {
		log.Printf("Token parsing error: %v\n", err)
		return false, nil
	}
	if !token.Valid {
		log.Println("Token is invalid")
		return false, nil
	}

	if exp, ok := claims["exp"].(float64); ok {
		expirationTime := time.Unix(int64(exp), 0)
		if time.Now().After(expirationTime) {
			log.Println("Token has expired")
			return false, nil
		}
	} else {
		log.Println("Expiration claim missing or invalid")
		return false, nil
	}

	return true, &claims
}

func IsUserAuthorized(operation models.Operation, roles []string) (bool, error) {
	authService, err := GetInstance()
	if err != nil {
		log.Printf("Error when loading auth service: %v\n", err)
		return false, err
	}

	if authService.HasPermission(roles, &operation) {
		return true, nil
	}
	
	return false, nil
}

func ExtractRoles(claims *jwt.MapClaims) ([]string, error) {
	var roles []string
	if resourceAccess, ok := (*claims)["resource_access"].(map[string]interface{}); ok {
		for _, resource := range resourceAccess {
			if resourceMap, ok := resource.(map[string]interface{}); ok {
				if resourceRoles, ok := resourceMap["roles"].([]interface{}); ok {
					for _, role := range resourceRoles {
						if roleStr, ok := role.(string); ok {
							roles = append(roles, roleStr)
						}
					}
				}
			}
		}
	} else {
		return nil, errors.New("resource_access claim missing or invalid")
	}

	return roles, nil
}
