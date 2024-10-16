package controllers

import (
	"encoding/json"
	"fmt"
	"github.com/ZPI-2024-25/KubernetesAccessManager/cluster"
	"github.com/ZPI-2024-25/KubernetesAccessManager/models"
	"github.com/gorilla/mux"
	"net/http"
)

func GetResourceController(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	params := mux.Vars(r)
	resourceType := params["resourceType"]
	resourceName := params["resourceName"]

	queryParams := r.URL.Query()
	namespace := queryParams.Get("namespace")

	resource, err := cluster.GetResource(resourceType, namespace, resourceName)
	if err != nil {
		w.WriteHeader(int(err.Code))
		json.NewEncoder(w).Encode(err)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(resource)
}

func ListResourcesController(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	params := mux.Vars(r)
	resourceType := params["resourceType"]

	queryParams := r.URL.Query()
	namespace := queryParams.Get("namespace")

	resources, err := cluster.ListResources(resourceType, namespace)
	if err != nil {
		w.WriteHeader(int(err.Code))
		json.NewEncoder(w).Encode(err)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(resources)
}

func CreateResourceController(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	params := mux.Vars(r)
	resourceType := params["resourceType"]

	queryParams := r.URL.Query()
	namespace := queryParams.Get("namespace")

	var resource models.ResourceDetails
	jsonErr := json.NewDecoder(r.Body).Decode(&resource.ResourceDetails)
	if jsonErr != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(models.ModelError{Code: 404, Message: "Invalid request body"})
		return
	}

	resource, err := cluster.CreateResource(resourceType, namespace, resource)
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(int(err.Code))
		json.NewEncoder(w).Encode(err)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(resource)
}

func DeleteClusterResourceController(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	params := mux.Vars(r)
	resourceType := params["resourceType"]
	resourceName := params["resourceName"]

	queryParams := r.URL.Query()
	namespace := queryParams.Get("namespace")

	err := cluster.DeleteResource(resourceType, namespace, resourceName)
	if err != nil {
		w.WriteHeader(int(err.Code))
		json.NewEncoder(w).Encode(err)
		return
	}

	w.WriteHeader(http.StatusOK)

	status := models.Status{
		Status:  "Success",
		Code:    http.StatusOK,
		Message: fmt.Sprintf("Resource %s deleted successfully", resourceName),
	}
	json.NewEncoder(w).Encode(status)
}

func UpdateResourceController(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	params := mux.Vars(r)

	resourceType := params["resourceType"]
	resourceName := params["resourceName"]

	queryParams := r.URL.Query()
	namespace := queryParams.Get("namespace")

	var resource models.ResourceDetails
	jsonErr := json.NewDecoder(r.Body).Decode(&resource.ResourceDetails)
	if jsonErr != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(models.ModelError{Code: 404, Message: "Invalid request body"})
		return
	}

	resource, err := cluster.UpdateResource(resourceType, namespace, resourceName, resource)
	if err != nil {
		w.WriteHeader(int(err.Code))
		json.NewEncoder(w).Encode(err)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(resource)
}
