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

func checkVersion(version int32) *models.ModelError {
	if version < 0 {
		return &models.ModelError{Code: 400, Message: "Invalid version"}
	}
	return nil
}
