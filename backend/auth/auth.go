package auth

import (
	"fmt"
	"github.com/MicahParks/keyfunc"
	"github.com/golang-jwt/jwt/v4"
	"log"
	"net/http"
	"strings"
	"time"
)

// Globalna zmienna przechowująca klucze publiczne z Keycloak
var jwks *keyfunc.JWKS

func init() {
	// URL JWKS z Keycloak - Zaktualizuj URL zgodnie z Twoją konfiguracją Keycloak
	jwksURL := "http://localhost:4000/realms/ZPI-realm/protocol/openid-connect/certs"

	// Inicjalizuj JWKS
	var err error
	jwks, err = keyfunc.Get(jwksURL, keyfunc.Options{
		RefreshInterval: time.Hour, // Automatyczne odświeżanie kluczy co godzinę
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

	// Sprawdź, czy nagłówek zaczyna się od "Bearer"
	parts := strings.Split(authHeader, " ")
	if len(parts) != 2 || parts[0] != "Bearer" {
		return "", fmt.Errorf("invalid authorization header format")
	}

	return parts[1], nil
}

func IsTokenValid(tokenStr string) bool {
	// Zweryfikuj token JWT przy użyciu kluczy JWKS
	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		return jwks.Keyfunc(token)
	})
	if err != nil {
		fmt.Printf("Token parsing error: %v\n", err)
		return false
	}

	// Sprawdź, czy token jest ważny i podpis jest poprawny
	if !token.Valid {
		return false
	}

	// Sprawdź datę wygaśnięcia tokena (exp)
	if claims, ok := token.Claims.(jwt.MapClaims); ok {
		if exp, ok := claims["exp"].(float64); ok {
			expirationTime := time.Unix(int64(exp), 0)
			log.Printf("%s", expirationTime.String())
			if time.Now().After(expirationTime) {
				fmt.Println("Token has expired")
				return false
			}
		} else {
			fmt.Println("Expiration claim missing or invalid")
			return false
		}
	} else {
		fmt.Println("Invalid token claims")
		return false
	}

	return true
}
