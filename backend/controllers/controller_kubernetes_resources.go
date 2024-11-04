package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/ZPI-2024-25/KubernetesAccessManager/auth"
	"github.com/ZPI-2024-25/KubernetesAccessManager/cluster"
	"github.com/ZPI-2024-25/KubernetesAccessManager/models"
	"github.com/gorilla/mux"
	"k8s.io/utils/env"
)

func GetResourceController(w http.ResponseWriter, r *http.Request) {
	handleResourceOperation(w, r, models.Read, func(resourceType, namespace, resourceName string) (interface{}, *models.ModelError) {
		return cluster.GetResource(resourceType, namespace, resourceName, cluster.GetResourceInterface)
	})
}

func ListResourcesController(w http.ResponseWriter, r *http.Request) {
	handleResourceOperation(w, r, models.List, func(resourceType, namespace, _ string) (interface{}, *models.ModelError) {
		return cluster.ListResources(resourceType, namespace, cluster.GetResourceInterface)
	})
}

func CreateResourceController(w http.ResponseWriter, r *http.Request) {
	handleResourceOperation(w, r, models.Create, func(resourceType, namespace, _ string) (interface{}, *models.ModelError) {
		var resource models.ResourceDetails
		if !decodeJSONBody(r, &resource.ResourceDetails) {
			return nil, &models.ModelError{Code: http.StatusBadRequest, Message: "Invalid request body"}
		}
		return cluster.CreateResource(resourceType, namespace, resource, cluster.GetResourceInterface)
	})
}

func DeleteResourceController(w http.ResponseWriter, r *http.Request) {
	handleResourceOperation(w, r, models.Delete, func(resourceType, namespace, resourceName string) (interface{}, *models.ModelError) {
		if err := cluster.DeleteResource(resourceType, namespace, resourceName, cluster.GetResourceInterface); err != nil {
			return nil, err
		}
		return models.Status{
			Status:  "Success",
			Code:    http.StatusOK,
			Message: fmt.Sprintf("Resource %s deleted successfully", resourceName),
		}, nil
	})
}

func UpdateResourceController(w http.ResponseWriter, r *http.Request) {
	handleResourceOperation(w, r, models.Update, func(resourceType, namespace, resourceName string) (interface{}, *models.ModelError) {
		var resource models.ResourceDetails
		if !decodeJSONBody(r, &resource.ResourceDetails) {
			return nil, &models.ModelError{Code: http.StatusBadRequest, Message: "Invalid request body"}
		}
		return cluster.UpdateResource(resourceType, namespace, resourceName, resource, cluster.GetResourceInterface)
	})
}

func handleResourceOperation(w http.ResponseWriter, r *http.Request, opType models.OperationType, operationFunc func(string, string, string) (interface{}, *models.ModelError)) {
	resourceType := getResourceType(r)
	resourceName := getResourceName(r)
	namespace := getNamespace(r)

	operation := models.Operation{
		Resource:  resourceType,
		Namespace: namespace,
		Type:      opType,
	}

	if err := authenticateAndAuthorize(w, r, operation); err != nil {
		writeJSONResponse(w, int(err.Code), err)
		return
	}

	result, err := operationFunc(resourceType, namespace, resourceName)
	if err != nil {
		writeJSONResponse(w, int(err.Code), err)
		return
	}

	statusCode := http.StatusOK
	if opType == models.Create {
		statusCode = http.StatusCreated
	}

	writeJSONResponse(w, statusCode, result)
}

func authenticateAndAuthorize(w http.ResponseWriter, r *http.Request, operation models.Operation) *models.ModelError {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	
	// temporary solution to disable auth if we don't have keycloak running
	if env.GetString("KEYCLOAK_URL", "") == "" {
		return nil
	}
	token, err := auth.GetJWTTokenFromHeader(r)
	isValid, claims := auth.IsTokenValid(token)

	if err != nil || !isValid {
		return &models.ModelError{
			Code:    http.StatusUnauthorized,
			Message: "Authentication failed",
		}
	}

	roles, err := auth.ExtractRoles(claims)
	if err != nil {
		return &models.ModelError{
			Code:    http.StatusBadRequest,
			Message: "Roles not found in bearer token",
		}
	}

	authorized, err := auth.IsUserAuthorized(operation, roles)
	if err != nil {
		return &models.ModelError{
			Code:    http.StatusInternalServerError,
			Message: "Internal Server Error",
		}
	}

	if !authorized {
		return &models.ModelError{
			Code:    http.StatusForbidden,
			Message: "Insufficient permissions",
		}
	}
	return nil
}

func getResourceType(r *http.Request) string {
	return mux.Vars(r)["resourceType"]
}

func getResourceName(r *http.Request) string {
	return mux.Vars(r)["resourceName"]
}

func getNamespace(r *http.Request) string {
	return r.URL.Query().Get("namespace")
}

func writeJSONResponse(w http.ResponseWriter, statusCode int, data interface{}) {
	if w.Header().Get("Content-Type") == "" {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	}
	if statusCode != http.StatusOK {
		w.WriteHeader(statusCode)
	}
	if err := json.NewEncoder(w).Encode(data); err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}

func decodeJSONBody(r *http.Request, dst interface{}) bool {
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(dst)
	if err != nil {
		return false
	}
	return true
}
