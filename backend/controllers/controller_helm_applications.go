package controllers

import (
	"github.com/ZPI-2024-25/KubernetesAccessManager/helm"
	"net/http"
)

func GetHelmReleaseController(w http.ResponseWriter, r *http.Request) {
	setJSONContentType(w)

	releaseName := getReleaseName(r)
	namespace := getNamespace(r)

	release, err := helm.GetHelmRelease(releaseName, namespace)
	if err != nil {
		writeJSONResponse(w, int(err.Code), err)
		return
	}

	writeJSONResponse(w, http.StatusOK, release)
}

func GetHelmReleaseHistoryController(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
}

func ListHelmReleasesController(w http.ResponseWriter, r *http.Request) {
	setJSONContentType(w)

	namespace := getNamespace(r)

	releases, err := helm.ListHelmReleases(namespace)
	if err != nil {
		writeJSONResponse(w, int(err.Code), err)
		return
	}

	writeJSONResponse(w, http.StatusOK, releases)
}

func RollbackHelmReleaseController(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
}

func UninstallHelmReleaseController(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
}
