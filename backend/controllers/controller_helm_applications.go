package controllers

import (
	"fmt"
	"net/http"
	"time"

	"github.com/ZPI-2024-25/KubernetesAccessManager/auth"
	"github.com/ZPI-2024-25/KubernetesAccessManager/common"
	"github.com/ZPI-2024-25/KubernetesAccessManager/helm"
	"github.com/ZPI-2024-25/KubernetesAccessManager/models"
	"k8s.io/utils/env"
)

const (
	DefaultOperationTimeout = 5 * time.Second
)

func GetHelmReleaseController(w http.ResponseWriter, r *http.Request) {
	handleHelmOperation(w, r, models.Read, func(releaseName, namespace string) (interface{}, *models.ModelError) {
		return helm.GetHelmRelease(releaseName, namespace, helm.PrepareActionConfig)
	})
}

func GetHelmReleaseHistoryController(w http.ResponseWriter, r *http.Request) {
	handleHelmOperation(w, r, models.Read, func(releaseName, namespace string) (interface{}, *models.ModelError) {
		return helm.GetHelmReleaseHistory(releaseName, namespace, helm.PrepareActionConfig)
	})
}

func ListHelmReleasesController(w http.ResponseWriter, r *http.Request) {
	handleHelmOperation(w, r, models.List, func(releaseName, namespace string) (interface{}, *models.ModelError) {
		releases, err := helm.ListHelmReleases(namespace, helm.PrepareActionConfig)
		if err != nil {
			return nil, err
		}
		if namespace != "" {
			return releases, nil
		}

		// temporary solution to disable auth if we don't have keycloak running
		if env.GetString("KEYCLOAK_URL", "") == "" {
			return releases, nil
		}
		token, err2 := auth.GetJWTTokenFromHeader(r)
		isValid, claims := auth.IsTokenValid(token)

		if err2 != nil || !isValid {
			return nil, &models.ModelError{
				Message: "Unauthorized",
				Code:    http.StatusUnauthorized,
			}
		}
		filtered, errM := auth.FilterRestrictedReleases(releases, claims)
		return filtered, errM
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
		release, completed, err := helm.RollbackHelmRelease(releaseName, namespace, int(version.Version), timeout, helm.PrepareActionConfig)
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
		completed, err := helm.UninstallHelmRelease(releaseName, namespace, timeout, helm.PrepareActionConfig)
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

	if namespace == "" && opType != models.List {
		namespace = common.DEFAULT_NAMESPACE
	}

	// The only operation that can be done for all namespaces - list without namespace mentioned
	if !(opType == models.List && namespace == "") {
		operation := models.Operation{
			Resource:  "Helm",
			Namespace: namespace,
			Type:      opType,
		}

		if err := authenticateAndAuthorize(r, operation); err != nil {
			writeJSONResponse(w, int(err.Code), err)
			return
		}
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
		if mystery, ok := result.(models.Status); ok {
			statusCode = int(mystery.Code)
		}
	}

	writeJSONResponse(w, statusCode, result)
}
