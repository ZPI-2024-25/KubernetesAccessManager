package helm

import (
	"errors"
	"github.com/ZPI-2024-25/KubernetesAccessManager/models"
	"helm.sh/helm/v3/pkg/storage/driver"
	"time"
)

type ActionConfigGetter func(namespace string, useDefaultNamespace bool) (ActionConfigInterface, *models.ModelError)

func GetHelmRelease(releaseName string, namespace string, getActionConfig ActionConfigGetter) (*models.HelmRelease, *models.ModelError) {
	actionConfig, cErr := getActionConfig(namespace, false)
	if cErr != nil {
		return nil, cErr
	}

	release, err := actionConfig.getRelease(releaseName)
	if err != nil {
		return nil, &models.ModelError{Code: 404, Message: "Release not found: " + err.Error()}
	}

	return getReleaseData(release), nil
}

func ListHelmReleases(namespace string, getActionConfig ActionConfigGetter) ([]models.HelmRelease, *models.ModelError) {
	actionConfig, cErr := getActionConfig(namespace, false)
	if cErr != nil {
		return nil, cErr
	}

	releases, err := actionConfig.listReleases(namespace == "")
	if err != nil {
		return nil, &models.ModelError{Code: 500, Message: "Failed to list releases: " + err.Error()}
	}

	var helmReleases []models.HelmRelease
	for _, release := range releases {
		helmReleases = append(helmReleases, *getReleaseData(release))
	}

	return helmReleases, nil
}

func UninstallHelmRelease(releaseName string, namespace string, timeout time.Duration, getActionConfig ActionConfigGetter) (bool, *models.ModelError) {
	actionConfig, cErr := getActionConfig(namespace, true)
	if cErr != nil {
		return false, cErr
	}

	errCh := make(chan error, 1)

	go func() {
		_, err := actionConfig.uninstallRelease(releaseName)
		errCh <- err
	}()

	select {
	case err := <-errCh:
		if err != nil {
			if errors.Is(err, driver.ErrReleaseNotFound) {
				return false, &models.ModelError{Code: 404, Message: "Release not found: " + err.Error()}
			}
			return false, &models.ModelError{Code: 500, Message: "Internal server error: " + err.Error()}
		}
		return true, nil
	case <-time.After(timeout):
		return false, nil
	}
}

func GetHelmReleaseHistory(releaseName string, namespace string, getActionConfig ActionConfigGetter) ([]models.HelmReleaseHistory, *models.ModelError) {
	actionConfig, cErr := getActionConfig(namespace, true)
	if cErr != nil {
		return nil, cErr
	}

	releases, err := actionConfig.getReleaseHistory(releaseName, 0)
	if err != nil {
		return nil, &models.ModelError{Code: 404, Message: "Failed to get release history"}
	}

	var helmReleases []models.HelmReleaseHistory
	for _, release := range releases {
		helmReleases = append(helmReleases, *getReleaseHistoryData(release))
	}

	return helmReleases, nil
}

func RollbackHelmRelease(releaseName string, namespace string, version int, timeout time.Duration, getActionConfig ActionConfigGetter) (*models.HelmRelease, bool, *models.ModelError) {
	actionConfig, cErr := getActionConfig(namespace, true)
	if cErr != nil {
		return nil, false, cErr
	}

	type rollbackResult struct {
		release *models.HelmRelease
		err     error
	}

	resultCh := make(chan rollbackResult, 1)

	go func() {
		err := actionConfig.rollbackRelease(releaseName, version)
		if err != nil {
			resultCh <- rollbackResult{nil, err}
			return
		}
		release, err := actionConfig.getRelease(releaseName)
		if err != nil {
			resultCh <- rollbackResult{nil, err}
			return
		}
		resultCh <- rollbackResult{getReleaseData(release), nil}
	}()

	select {
	case result := <-resultCh:
		if result.err != nil {
			if errors.Is(result.err, driver.ErrReleaseNotFound) {
				return nil, false, &models.ModelError{Code: 404, Message: "Release not found: " + result.err.Error()}
			}
			return nil, false, &models.ModelError{Code: 500, Message: "Internal server error: " + result.err.Error()}
		}
		return result.release, true, nil
	case <-time.After(timeout):
		return nil, false, nil
	}
}
