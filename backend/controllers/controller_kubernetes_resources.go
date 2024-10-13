package controllers

import (
	"encoding/json"
	"fmt"
	"github.com/ZPI-2024-25/KubernetesUserManager/cluster"
	"github.com/ZPI-2024-25/KubernetesUserManager/models"
	"github.com/gorilla/mux"
	"net/http"
)

func GetClusterResourceController(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)

	params := mux.Vars(r)

	resource, err := cluster.GetResource(params["resourceType"], "", params["resourceName"])
	if err != nil {
		http.Error(w, err.Message, int(err.Code))
	}
	json.NewEncoder(w).Encode(resource)
}

func GetNamespacedResourceController(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)

	params := mux.Vars(r)

	resource, err := cluster.GetResource(params["resourceType"], params["namespace"], params["resourceName"])
	if err != nil {
		http.Error(w, err.Message, int(err.Code))
	}
	json.NewEncoder(w).Encode(resource)
}

func ListClusterResourcesController(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
}

func ListNamespacedResourcesController(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
}

func CreateClusterResourceController(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	params := mux.Vars(r)

	resourceType := params["resourceType"]

	var resource models.Resource
	jsonErr := json.NewDecoder(r.Body).Decode(&resource)
	if jsonErr != nil {
		http.Error(w, `{"error": "Invalid request body"}`, http.StatusBadRequest)
		return
	}

	resource, err := cluster.CreateResource(resourceType, "", resource)
	if err != nil {
		http.Error(w, err.Message, int(err.Code))
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(resource)
}

func CreateNamespacedResourceController(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	params := mux.Vars(r)

	resourceType := params["resourceType"]
	namespace := params["namespace"]

	var resource models.Resource
	jsonErr := json.NewDecoder(r.Body).Decode(&resource)
	if jsonErr != nil {
		http.Error(w, `{"error": "Invalid request body"}`, http.StatusBadRequest)
		return
	}

	resource, err := cluster.CreateResource(resourceType, namespace, resource)
	if err != nil {
		http.Error(w, err.Message, int(err.Code))
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(resource)
}

func DeleteClusterResourceController(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)

	params := mux.Vars(r)

	resourceType := params["resourceType"]
	resourceName := params["resourceName"]

	err := cluster.DeleteResource(resourceType, "", resourceName)
	if err != nil {
		http.Error(w, err.Message, int(err.Code))
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"status": "Resource deleted successfully"}`))
}

func DeleteNamespacedResourceController(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	params := mux.Vars(r)

	resourceType := params["resourceType"]
	namespace := params["namespace"]
	resourceName := params["resourceName"]

	err := cluster.DeleteResource(resourceType, namespace, resourceName)
	if err != nil {
		http.Error(w, err.Message, int(err.Code))
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"status": "Resource deleted successfully"}`))
}

func UpdateResourceController(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	params := mux.Vars(r)

	resourceType := params["resourceType"]
	namespace := params["namespace"]
	resourceName := params["resourceName"]

	var resource models.Resource
	jsonErr := json.NewDecoder(r.Body).Decode(&resource)
	if jsonErr != nil {
		http.Error(w, `{"error": "Invalid request body"}`, http.StatusBadRequest)
		return
	}

	resource, err := cluster.UpdateResource(resourceType, namespace, resourceName, resource)
	if err != nil {
		fmt.Println(err)
		http.Error(w, err.Message, int(err.Code))
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(resource)
}
