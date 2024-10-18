package controllers

import (
	"encoding/json"
	"github.com/ZPI-2024-25/KubernetesAccessManager/kubernetes"
	"github.com/ZPI-2024-25/KubernetesAccessManager/models"
	"github.com/gorilla/mux"
	"io"
	"net/http"
)

func CreateResourceController(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	bodyVal, err := io.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode("Invalid request body")
		return
	}

	resourceManifest := string(bodyVal)

	resource, creationErr := kubernetes.ResourceCreation(resourceManifest)

	if creationErr != nil {
		w.WriteHeader(int(creationErr.Code))
		json.NewEncoder(w).Encode(creationErr)
	} else {
		w.WriteHeader(201)
		json.NewEncoder(w).Encode(resource)
	}
}

func DeleteResourceController(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	params := mux.Vars(r)
	_, err := kubernetes.ResourceGet(params["resourceType"], params["resourceName"])
	if err != nil {
		w.WriteHeader(int(err.Code))
		json.NewEncoder(w).Encode(err)
	} else {
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(models.Status{Status: "Usuniete", Message: "No usuniete, przeciez pisze", Code: 200})
	}
}

func GetResourceController(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	params := mux.Vars(r)
	resource, err := kubernetes.ResourceGet(params["resourceType"], params["resourceName"])
	if err != nil {
		w.WriteHeader(int(err.Code))
		json.NewEncoder(w).Encode(err)
	} else {
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(resource)
	}
}

func ListResourcesController(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	params := mux.Vars(r)
	resource, err := kubernetes.ResourceListing(params["resourceType"])
	if err != nil {
		w.WriteHeader(int(err.Code))
		json.NewEncoder(w).Encode(err)
	} else {
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(resource)
	}
}

func UpdateResourceController(w http.ResponseWriter, r *http.Request) {
	GetResourceController(w, r)
}
