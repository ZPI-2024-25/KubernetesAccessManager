package controllers

import (
	"encoding/json"
	"fmt"
	"github.com/ZPI-2024-25/KubernetesAccessManager/auth"
	"github.com/ZPI-2024-25/KubernetesAccessManager/cluster"
	"github.com/ZPI-2024-25/KubernetesAccessManager/models"
	"github.com/gorilla/mux"
	"net/http"
)

func GetResourceController(w http.ResponseWriter, r *http.Request) {
	if err := setJSONCTAndAuth(w, r); err != nil {
		writeJSONResponse(w, int(err.Code), err)
		return
	}
	resourceType := getResourceType(r)
	resourceName := getResourceName(r)
	namespace := getNamespace(r)

	resource, err := cluster.GetResource(resourceType, namespace, resourceName)
	if err != nil {
		writeJSONResponse(w, int(err.Code), err)
		return
	}

	writeJSONResponse(w, http.StatusOK, resource)
}

func ListResourcesController(w http.ResponseWriter, r *http.Request) {
	if err := setJSONCTAndAuth(w, r); err != nil {
		writeJSONResponse(w, int(err.Code), err)
		return
	}
	resourceType := getResourceType(r)
	namespace := getNamespace(r)

	resources, err := cluster.ListResources(resourceType, namespace)
	if err != nil {
		writeJSONResponse(w, int(err.Code), err)
		return
	}

	writeJSONResponse(w, http.StatusOK, resources)
}

func CreateResourceController(w http.ResponseWriter, r *http.Request) {
	if err := setJSONCTAndAuth(w, r); err != nil {
		writeJSONResponse(w, int(err.Code), err)
		return
	}
	resourceType := getResourceType(r)
	namespace := getNamespace(r)

	var resource models.ResourceDetails
	if !decodeJSONBody(w, r, &resource.ResourceDetails) {
		return
	}

	resource, err := cluster.CreateResource(resourceType, namespace, resource)
	if err != nil {
		fmt.Println(err)
		writeJSONResponse(w, int(err.Code), err)
		return
	}

	writeJSONResponse(w, http.StatusCreated, resource)
}

func DeleteClusterResourceController(w http.ResponseWriter, r *http.Request) {
	if err := setJSONCTAndAuth(w, r); err != nil {
		writeJSONResponse(w, int(err.Code), err)
		return
	}
	resourceType := getResourceType(r)
	resourceName := getResourceName(r)
	namespace := getNamespace(r)

	err := cluster.DeleteResource(resourceType, namespace, resourceName)
	if err != nil {
		writeJSONResponse(w, int(err.Code), err)
		return
	}

	w.WriteHeader(http.StatusOK)

	status := models.Status{
		Status:  "Success",
		Code:    http.StatusOK,
		Message: fmt.Sprintf("Resource %s deleted successfully", resourceName),
	}
	writeJSONResponse(w, http.StatusOK, status)
}

func UpdateResourceController(w http.ResponseWriter, r *http.Request) {
	if err := setJSONCTAndAuth(w, r); err != nil {
		writeJSONResponse(w, int(err.Code), err)
		return
	}
	resourceType := getResourceType(r)
	resourceName := getResourceName(r)
	namespace := getNamespace(r)

	var resource models.ResourceDetails
	if !decodeJSONBody(w, r, &resource.ResourceDetails) {
		return
	}

	resource, err := cluster.UpdateResource(resourceType, namespace, resourceName, resource)
	if err != nil {
		writeJSONResponse(w, int(err.Code), err)
		return
	}

	writeJSONResponse(w, http.StatusOK, resource)
}

func setJSONCTAndAuth(w http.ResponseWriter, r *http.Request) *models.ModelError {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	unauthorizedError := models.ModelError{
		Code:    http.StatusUnauthorized,
		Message: "Unauthorized",
	}
	token, err := auth.GetJWTTokenFromHeader(r)
	if err != nil {
		return &unauthorizedError
	}
	if !auth.IsTokenValid(token) {
		return &unauthorizedError
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
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(data)
}

func decodeJSONBody(w http.ResponseWriter, r *http.Request, dst interface{}) bool {
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(dst)
	if err != nil {
		writeJSONResponse(w, 400, &models.ModelError{Code: 400, Message: "Invalid request body"})
		return false
	}
	return true
}
