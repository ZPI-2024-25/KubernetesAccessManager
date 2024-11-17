package controllers

import (
	"net/http"

	"github.com/ZPI-2024-25/KubernetesAccessManager/auth"
	"github.com/ZPI-2024-25/KubernetesAccessManager/models"
	"github.com/golang-jwt/jwt/v4"
)

func CheckLoginStatus(w http.ResponseWriter, r *http.Request)  {
	token, err := auth.GetJWTTokenFromHeader(r)
	isValid, claims := auth.IsTokenValid(token)

	if err != nil || !isValid {
		writeJSONResponse(w, http.StatusUnauthorized, models.ModelError{
			Message: "Unauthorized",
			Code: http.StatusUnauthorized,
		})
		return
	}
	rolemap, err := auth.GetRoleMapInstance()
	if err != nil {
		writeJSONResponse(w, http.StatusInternalServerError, &models.ModelError{
			Message: "Internal Server Error",
			Code: http.StatusInternalServerError,
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
		return nil, &models.ModelError{
			Message: "Roles not found in bearer token",
			Code: http.StatusInternalServerError,
		}
	}
	access := rolemap.GetAllPermissions(roles)
	auth.PrunePermissions(access)
	permissions := shortenPermissions(access)
	exp, preferredUsername, email := auth.ExtractUserStatus(claims)

	return &models.UserStatus{
		Permissions: permissions,
		User: &models.UserStatusUser{
			PreferredUsername: preferredUsername,
			Email: email,
			Exp: exp,
		},
	}, nil
}

func shortenPermissions(pmatrix map[string]map[string]map[models.OperationType]struct{}) map[string]map[string][]string {
	res := make(map[string]map[string][]string)

	for k, v := range pmatrix {
		res[k] = make(map[string][]string)
		for k2, v2 := range v {
			res[k][k2] = make([]string, len(v2))
			i := 0
			for k3 := range v2 {
				res[k][k2][i] = shorterP(k3)
				i++
			}
		}
	}
	return res
}

func shorterP(p models.OperationType) string {
	switch p {
		case models.Create:
			return "c"
		case models.Read:
			return "r"
		case models.Update:
			return "u"
		case models.Delete:
			return "d"
		case models.List:
			return "l"
		default:
			return "x"
		}
}