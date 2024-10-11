package controllers

import (
	"encoding/json"
	"github.com/ZPI-2024-25/KubernetesUserManager/go/cluster"
	"github.com/gorilla/mux"
	"net/http"
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
	w.WriteHeader(http.StatusOK)
}

func CreateNamespacedResourceController(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
}

func DeleteClusterResourceController(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)

	params := mux.Vars(r)

	cluster.DeleteClusterResource(params["resourceType"], params["resourceName"])
}

func DeleteNamespacedResourceController(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)

	params := mux.Vars(r)

	cluster.DeleteNamespacedResource(params["resourceType"], params["namespace"], params["resourceName"])
}
