package controllers

import (
	"encoding/json"
	"github.com/ZPI-2024-25/KubernetesAccessManager/kubernetes"
	"github.com/gorilla/mux"
	"net/http"
)

func CreateResourceController(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
}

func DeleteResourceController(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
}

func GetResourceController(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
}

func ListResourcesController(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	params := mux.Vars(r)

	resource, err := kubernetes.ResourceListing(params["resourceType"])
	if err != nil {
		http.Error(w, err.Message, int(err.Code))
	}
	json.NewEncoder(w).Encode(resource)
}

func UpdateResourceController(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
}
