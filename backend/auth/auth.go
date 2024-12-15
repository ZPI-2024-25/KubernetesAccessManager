package auth

import (
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/MicahParks/keyfunc"
	"github.com/ZPI-2024-25/KubernetesAccessManager/common"
	"github.com/ZPI-2024-25/KubernetesAccessManager/models"
	"github.com/golang-jwt/jwt/v4"
)

var jwks *keyfunc.JWKS

func InitializeAuth() {
	var err error
	log.Println("Connecting to auth provider on URL:", common.KeycloakJwksUrl)
	jwks, err = keyfunc.Get(common.KeycloakJwksUrl, keyfunc.Options{
		RefreshInterval: time.Hour,
	})
	if err != nil {
		log.Printf("Failed to connect to auth provider: %s", err)
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
	if jwks == nil {
		log.Printf("Authentication service not initialized\n")
		return false, nil
	}
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
	authService, err := GetRoleMapInstance()
	if err != nil {
		log.Printf("Error while getting Rolemap instance: %v\n", err)
		return false, err
	}

	if authService.HasPermission(roles, &operation) {
		return true, nil
	}

	return false, nil
}

func ExtractRoles(claims *jwt.MapClaims) ([]string, *models.ModelError) {
	var roles []string
	if realmAccess, ok := (*claims)["realm_access"].(map[string]interface{}); ok {
		extractRolesFromMapInterface(realmAccess, "roles", &roles)
	}
	if resourceAccess, ok := (*claims)["resource_access"].(map[string]interface{}); ok {
		if resource, ok := resourceAccess[common.KeycloakClient]; ok {
			if resourceMap, ok := resource.(map[string]interface{}); ok {
				extractRolesFromMapInterface(resourceMap, "roles", &roles)
			}
		}
	}
	return roles, nil
}

func extractRolesFromMapInterface(claims map[string]interface{}, rolekey string, roles *[]string) {
	if realm, ok := claims[rolekey].([]interface{}); ok {
		for _, role := range realm {
			if roleStr, ok := role.(string); ok {
				*roles = append((*roles), roleStr)
			}
		}
	}
}

func ExtractUserStatus(claims *jwt.MapClaims) (int32, string, string) {
	var exp int32
	var preferredUsername string
	var email string
	if expFloat, ok := (*claims)["exp"].(float64); ok {
		exp = int32(expFloat)
	}
	if preferredUsernameStr, ok := (*claims)["preferred_username"].(string); ok {
		preferredUsername = preferredUsernameStr
	}
	if emailStr, ok := (*claims)["email"].(string); ok {
		email = emailStr
	}
	return exp, preferredUsername, email
}

func FilterRestrictedResources(resources *models.ResourceList, claims *jwt.MapClaims, resourceType string) (*models.ResourceList, *models.ModelError) {
	namespaces := make(map[string]struct{})
	for _, resource := range resources.ResourceList {
		namespaces[resource.Namespace] = struct{}{}
	}
	allowed, err := getAllowedNamespaces(claims, resourceType, models.List, namespaces)
	if err != nil {
		return nil, err
	}
	filteredResources := make([]models.ResourceListResourceList, 0)
	for _, resource := range resources.ResourceList {
		if _, ok := allowed[resource.Namespace]; ok {
			filteredResources = append(filteredResources, resource)
		}
	}
	resources.ResourceList = filteredResources

	return resources, nil
}

func FilterRestrictedReleases(releases []models.HelmRelease, claims *jwt.MapClaims) ([]models.HelmRelease, *models.ModelError) {
	namespaces := make(map[string]struct{})
	for _, release := range releases {
		namespaces[release.Namespace] = struct{}{}
	}
	allowed, err := getAllowedNamespaces(claims, "Helm", models.List, namespaces)
	if err != nil {
		return nil, err
	}
	filteredReleases := make([]models.HelmRelease, 0)
	for _, release := range releases {
		if _, ok := allowed[release.Namespace]; ok {
			filteredReleases = append(filteredReleases, release)
		}
	}

	return filteredReleases, nil
}

func getAllowedNamespaces(claims *jwt.MapClaims, resourceType string, opType models.OperationType, namespaces map[string]struct{}) (map[string]struct{}, *models.ModelError) {
	roles, errM := ExtractRoles(claims)
	if errM != nil {
		return nil, errM
	}
	if len(roles) == 0 {
		return nil, &models.ModelError{
			Code:    http.StatusForbidden,
			Message: "No roles found in token",
		}
	}
	roleMap, err := GetRoleMapInstance()
	if err != nil {
		return nil, &models.ModelError{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
		}
	}
	if !roleMap.HasPermissionInAnyNamespace(roles, resourceType, opType) {
		return nil, &models.ModelError{
			Code:    http.StatusForbidden,
			Message: fmt.Sprintf("User does not have permission to %v resources", opType),
		}
	}

	allowed := make(map[string]struct{})
	for ns := range namespaces {
		op := models.Operation{
			Resource:  resourceType,
			Namespace: ns,
			Type:      opType,
		}
		hasPermission, err := IsUserAuthorized(op, roles)
		if err != nil {
			return nil, &models.ModelError{
				Code:    http.StatusInternalServerError,
				Message: fmt.Sprintf("Error when checking permissions: %v", err),
			}
		}
		if hasPermission {
			allowed[ns] = struct{}{}
		}
	}
	return allowed, nil
}
