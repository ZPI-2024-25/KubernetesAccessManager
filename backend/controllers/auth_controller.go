package controllers

import (
	"fmt"
	"net/http"

	"github.com/ZPI-2024-25/KubernetesAccessManager/auth"
	"github.com/ZPI-2024-25/KubernetesAccessManager/models"
	"github.com/golang-jwt/jwt/v4"
)

func CheckLoginStatus(w http.ResponseWriter, r *http.Request) {
	token, err := auth.GetJWTTokenFromHeader(r)
	isValid, claims := auth.IsTokenValid(token)

	if err != nil || !isValid {
		writeJSONResponse(w, http.StatusUnauthorized, models.ModelError{
			Message: "Unauthorized",
			Code:    http.StatusUnauthorized,
		})
		return
	}
	rolemap, err := auth.GetRoleMapInstance()
	if err != nil {
		writeJSONResponse(w, http.StatusInternalServerError, &models.ModelError{
			Message: fmt.Sprintf("Failed to get client: %s", err),
			Code:    http.StatusInternalServerError,
		})
	}
	status, errM := getLoginStatus(claims, rolemap)
	if errM != nil {
		writeJSONResponse(w, int(errM.Code), errM)
		return
	}
	writeJSONResponse(w, http.StatusOK, status)
}

func getLoginStatus(claims *jwt.MapClaims, rolemap *auth.RoleMapRepository) (*models.UserStatus, *models.ModelError) {
	roles, err := auth.ExtractRoles(claims)
	if err != nil {
		return nil, err
	}
	access := rolemap.GetAllPermissions(roles)
	auth.PrunePermissions(access)
	permissions := toPermissionModel(access)
	exp, preferredUsername, email := auth.ExtractUserStatus(claims)

	return &models.UserStatus{
		Permissions: permissions,
		User: &models.UserStatusUser{
			PreferredUsername: preferredUsername,
			Email:             email,
			Exp:               exp,
		},
	}, nil
}

func toPermissionModel(pmatrix auth.PermissionMatrix) map[string]map[string][]string {
	result := make(map[string]map[string][]string)

	for namespace, resources := range pmatrix {
		result[namespace] = make(map[string][]string)
		for resource, operations := range resources {
			result[namespace][resource] = make([]string, len(operations))
			i := 0
			for op := range operations {
				result[namespace][resource][i] = op.ShortString()
				i++
			}
		}
	}
	return result
}
