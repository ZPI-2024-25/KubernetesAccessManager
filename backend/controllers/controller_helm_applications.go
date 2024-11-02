package controllers

import (
	"fmt"
	"github.com/ZPI-2024-25/KubernetesAccessManager/helm"
	"github.com/ZPI-2024-25/KubernetesAccessManager/models"
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
	setJSONContentType(w)

	releaseName := getReleaseName(r)
	namespace := getNamespace(r)

	release, err := helm.GetHelmReleaseHistory(releaseName, namespace)
	if err != nil {
		writeJSONResponse(w, int(err.Code), err)
		return
	}

	writeJSONResponse(w, http.StatusOK, release)
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
	setJSONContentType(w)

	releaseName := getReleaseName(r)
	namespace := getNamespace(r)

	var version models.ReleaseNameRollbackBody
	if !decodeJSONBody(w, r, &version) {
		return
	}
	err := checkVersion(version.Version)
	if err != nil {
		writeJSONResponse(w, int(err.Code), err)
		return
	}

	release, err := helm.RollbackHelmRelease(releaseName, namespace, int(version.Version))
	if err != nil {
		writeJSONResponse(w, int(err.Code), err)
		return
	}

	writeJSONResponse(w, http.StatusOK, release)
}

func UninstallHelmReleaseController(w http.ResponseWriter, r *http.Request) {
	setJSONContentType(w)

	releaseName := getReleaseName(r)
	namespace := getNamespace(r)

	err := helm.UninstallHelmRelease(releaseName, namespace)
	if err != nil {
		writeJSONResponse(w, int(err.Code), err)
		return
	}

	status := models.Status{
		Status:  "Success",
		Code:    200,
		Message: fmt.Sprintf("Release %s uninstalled successfully", releaseName),
	}
	writeJSONResponse(w, http.StatusOK, status)
}
