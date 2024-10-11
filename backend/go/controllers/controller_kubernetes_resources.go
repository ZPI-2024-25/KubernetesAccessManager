package controllers

import (
	"encoding/json"
	"github.com/ZPI-2024-25/KubernetesUserManager/go/cluster"
	"github.com/ZPI-2024-25/KubernetesUserManager/go/models"
	"github.com/gorilla/mux"
	"net/http"
	"strings"
)

func GetClusterResourceController(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)

	params := mux.Vars(r)

	resource, _ := cluster.GetClusterResource(params["resourceType"], params["resourceName"])
	json.NewEncoder(w).Encode(resource)
}

func GetNamespacedResourceController(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)

	params := mux.Vars(r)

	resource, _ := cluster.GetNamespacedResource(params["resourceType"], params["namespace"], params["resourceName"])
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
	err := json.NewDecoder(r.Body).Decode(&resource)
	if err != nil {
		http.Error(w, `{"error": "Invalid request body"}`, http.StatusBadRequest)
		return
	}

	resource, err = cluster.CreateClusterResource(resourceType, resource)
	if err != nil {
		if strings.Contains(err.Error(), "unauthorized") {
			http.Error(w, `{"error": "Unauthorized"}`, http.StatusUnauthorized)
			return
		} else if strings.Contains(err.Error(), "forbidden") {
			http.Error(w, `{"error": "Forbidden"}`, http.StatusForbidden)
			return
		} else if strings.Contains(err.Error(), "invalid") {
			http.Error(w, `{"error": "Invalid input"}`, http.StatusBadRequest)
			return
		} else {
			http.Error(w, `{"error": "Internal server error"}`, http.StatusInternalServerError)
			return
		}
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
	err := json.NewDecoder(r.Body).Decode(&resource)
	if err != nil {
		http.Error(w, `{"error": "Invalid request body"}`, http.StatusBadRequest)
		return
	}

	resource, err = cluster.CreateNamespacedResource(resourceType, namespace, resource)
	if err != nil {
		if strings.Contains(err.Error(), "unauthorized") {
			http.Error(w, `{"error": "Unauthorized"}`, http.StatusUnauthorized)
			return
		} else if strings.Contains(err.Error(), "forbidden") {
			http.Error(w, `{"error": "Forbidden"}`, http.StatusForbidden)
			return
		} else if strings.Contains(err.Error(), "invalid") {
			http.Error(w, `{"error": "Invalid input"}`, http.StatusBadRequest)
			return
		} else {
			http.Error(w, `{"error": "Internal server error"}`, http.StatusInternalServerError)
			return
		}
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

	err := cluster.DeleteClusterResource(resourceType, resourceName)
	if err != nil {
		if strings.Contains(err.Error(), "invalid resource type") {
			http.Error(w, `{"error": "Invalid resource type"}`, http.StatusBadRequest)
			return
		} else if strings.Contains(err.Error(), "not found") {
			http.Error(w, `{"error": "Resource not found"}`, http.StatusNotFound)
			return
		} else if strings.Contains(err.Error(), "forbidden") {
			http.Error(w, `{"error": "Forbidden"}`, http.StatusForbidden)
			return
		} else if strings.Contains(err.Error(), "unauthorized") {
			http.Error(w, `{"error": "Unauthorized"}`, http.StatusUnauthorized)
			return
		} else {
			http.Error(w, `{"error": "Internal server error"}`, http.StatusInternalServerError)
			return
		}
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

	err := cluster.DeleteNamespacedResource(resourceType, namespace, resourceName)
	if err != nil {
		if strings.Contains(err.Error(), "invalid resource type") {
			http.Error(w, `{"error": "Invalid resource type"}`, http.StatusBadRequest)
			return
		} else if strings.Contains(err.Error(), "not found") {
			http.Error(w, `{"error": "Resource not found"}`, http.StatusNotFound)
			return
		} else if strings.Contains(err.Error(), "forbidden") {
			http.Error(w, `{"error": "Forbidden"}`, http.StatusForbidden)
			return
		} else if strings.Contains(err.Error(), "unauthorized") {
			http.Error(w, `{"error": "Unauthorized"}`, http.StatusUnauthorized)
			return
		} else {
			http.Error(w, `{"error": "Internal server error"}`, http.StatusInternalServerError)
			return
		}
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"status": "Resource deleted successfully"}`))
}
