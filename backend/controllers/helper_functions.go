package controllers

import (
	"encoding/json"
	"github.com/ZPI-2024-25/KubernetesAccessManager/models"
	"github.com/gorilla/mux"
	"net/http"
)

func setJSONContentType(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
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

func getReleaseName(r *http.Request) string {
	return mux.Vars(r)["releaseName"]
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
