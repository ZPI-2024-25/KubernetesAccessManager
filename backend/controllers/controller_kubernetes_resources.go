package controllers

import (
	"encoding/json"
	"fmt"
	"github.com/ZPI-2024-25/KubernetesAccessManager/cluster"
	"github.com/ZPI-2024-25/KubernetesAccessManager/models"
	"github.com/gorilla/mux"
	"net/http"
)

//func GetClusterResourceController(w http.ResponseWriter, r *http.Request) {
//	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
//	w.WriteHeader(http.StatusOK)
//
//	params := mux.Vars(r)
//
//	resource, err := cluster.GetResource(params["resourceType"], "", params["resourceName"])
//	if err != nil {
//		http.Error(w, err.Message, int(err.Code))
//	}
//	json.NewEncoder(w).Encode(resource)
//}

//func GetNamespacedResourceController(w http.ResponseWriter, r *http.Request) {
//	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
//	w.WriteHeader(http.StatusOK)
//
//	params := mux.Vars(r)
//
//	resource, err := cluster.GetResource(params["resourceType"], params["namespace"], params["resourceName"])
//	if err != nil {
//		http.Error(w, err.Message, int(err.Code))
//	}
//	json.NewEncoder(w).Encode(resource)
//}

func ListResourcesController(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	params := mux.Vars(r)
	resourceType := params["resourceType"]

	queryParams := r.URL.Query()
	namespace := queryParams.Get("namespace")

	resources, err := cluster.ListResources(resourceType, namespace)
	if err != nil {
		http.Error(w, err.Message, int(err.Code))
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
	jsonErr := json.NewDecoder(r.Body).Decode(&resource)
	if jsonErr != nil {
		fmt.Println(jsonErr)
		http.Error(w, `{"error": "Invalid request body"}`, http.StatusBadRequest)
		return
	}

	resource, err := cluster.CreateResource(resourceType, namespace, resource)
	if err != nil {
		fmt.Println(err)
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

	queryParams := r.URL.Query()
	namespace := queryParams.Get("namespace")

	err := cluster.DeleteResource(resourceType, namespace, resourceName)
	if err != nil {
		http.Error(w, err.Message, int(err.Code))
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"status": "Resource deleted successfully"}`))
}

//func UpdateResourceController(w http.ResponseWriter, r *http.Request) {
//	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
//	params := mux.Vars(r)
//
//	resourceType := params["resourceType"]
//	namespace := params["namespace"]
//	resourceName := params["resourceName"]
//
//	var resource models.Resource
//	jsonErr := json.NewDecoder(r.Body).Decode(&resource)
//	if jsonErr != nil {
//		http.Error(w, `{"error": "Invalid request body"}`, http.StatusBadRequest)
//		return
//	}
//
//	resource, err := cluster.UpdateResource(resourceType, namespace, resourceName, resource)
//	if err != nil {
//		fmt.Println(err)
//		http.Error(w, err.Message, int(err.Code))
//	}
//
//	w.WriteHeader(http.StatusOK)
//	json.NewEncoder(w).Encode(resource)
//}
