package controllers

import (
	"fmt"
	"github.com/ZPI-2024-25/KubernetesAccessManager/helm"
	"github.com/ZPI-2024-25/KubernetesAccessManager/models"
	"net/http"
	"time"
)

const (
	DefaultOperationTimeout = 5 * time.Second
)

func GetHelmReleaseController(w http.ResponseWriter, r *http.Request) {
	handleHelmOperation(w, r, models.Read, func(releaseName, namespace string) (interface{}, *models.ModelError) {
		return helm.GetHelmRelease(releaseName, namespace)
	})
}

func GetHelmReleaseHistoryController(w http.ResponseWriter, r *http.Request) {
	handleHelmOperation(w, r, models.Read, func(releaseName, namespace string) (interface{}, *models.ModelError) {
		return helm.GetHelmReleaseHistory(releaseName, namespace)
	})
}

func ListHelmReleasesController(w http.ResponseWriter, r *http.Request) {
	handleHelmOperation(w, r, models.List, func(releaseName, namespace string) (interface{}, *models.ModelError) {
		return helm.ListHelmReleases(namespace)
	})
}

func RollbackHelmReleaseController(w http.ResponseWriter, r *http.Request) {
	handleHelmOperation(w, r, models.Update, func(releaseName, namespace string) (interface{}, *models.ModelError) {
		var version models.ReleaseNameRollbackBody
		if !decodeJSONBody(r, &version) {
			return nil, &models.ModelError{Code: http.StatusBadRequest, Message: "Invalid request body"}
		}
		if err := checkVersion(version.Version); err != nil {
			return nil, err
		}

		timeout := DefaultOperationTimeout
		release, completed, err := helm.RollbackHelmRelease(releaseName, namespace, int(version.Version), timeout)
		if err != nil {
			return nil, err
		}

		if completed {
			return release, nil
		} else {
			return models.Status{
				Status:  "Accepted",
				Code:    202,
				Message: fmt.Sprintf("Rolling back release %s to version %d in progress", releaseName, version.Version),
			}, nil
		}
	})
}

func UninstallHelmReleaseController(w http.ResponseWriter, r *http.Request) {
	handleHelmOperation(w, r, models.Delete, func(releaseName, namespace string) (interface{}, *models.ModelError) {
		timeout := DefaultOperationTimeout
		completed, err := helm.UninstallHelmRelease(releaseName, namespace, timeout)
		if err != nil {
			return nil, err
		}

		if completed {
			return models.Status{
				Status:  "Success",
				Code:    200,
				Message: fmt.Sprintf("Release %s uninstalled successfully", releaseName),
			}, nil
		} else {
			return models.Status{
				Status:  "Accepted",
				Code:    202,
				Message: fmt.Sprintf("Uninstalling release %s in progress", releaseName),
			}, nil
		}
	})
}

func handleHelmOperation(w http.ResponseWriter, r *http.Request, opType models.OperationType, operationFunc func(string, string) (interface{}, *models.ModelError)) {
	releaseName := getReleaseName(r)
	namespace := getNamespace(r)

	operation := models.Operation{
		Resource:  "Helm",
		Namespace: namespace,
		Type:      opType,
	}

	if err := authenticateAndAuthorize(r, operation); err != nil {
		writeJSONResponse(w, int(err.Code), err)
		return
	}

	result, err := operationFunc(releaseName, namespace)
	if err != nil {
		writeJSONResponse(w, int(err.Code), err)
		return
	}

	statusCode := http.StatusOK
	if opType == models.Create {
		statusCode = http.StatusCreated
	} else if opType == models.Delete {
		status := result.(models.Status)
		statusCode = int(status.Code)
	} else if opType == models.Update {
		status := result.(models.Status)
		statusCode = int(status.Code)
	}

	writeJSONResponse(w, statusCode, result)
}
