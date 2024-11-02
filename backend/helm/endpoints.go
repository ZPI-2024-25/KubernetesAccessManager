package helm

import (
	"errors"
	"github.com/ZPI-2024-25/KubernetesAccessManager/models"
	"helm.sh/helm/v3/pkg/storage/driver"
)

func GetHelmRelease(releaseName string, namespace string) (*models.HelmRelease, *models.ModelError) {
	actionConfig, cErr := PrepareActionConfig(namespace, false)
	if cErr != nil {
		return nil, cErr
	}

	release, err := getRelease(actionConfig, releaseName)
	if err != nil {
		return nil, &models.ModelError{Code: 404, Message: "Release not found: " + err.Error()}
	}

	return GetReleaseData(release), nil
}

func ListHelmReleases(namespace string) ([]models.HelmRelease, *models.ModelError) {
	actionConfig, cErr := PrepareActionConfig(namespace, false)
	if cErr != nil {
		return nil, cErr
	}

	releases, err := listReleases(actionConfig, namespace == "")
	if err != nil {
		return nil, &models.ModelError{Code: 500, Message: "Failed to list releases: " + err.Error()}
	}

	var helmReleases []models.HelmRelease
	for _, release := range releases {
		helmReleases = append(helmReleases, *GetReleaseData(release))
	}

	return helmReleases, nil
}

func UninstallHelmRelease(releaseName string, namespace string) *models.ModelError {
	actionConfig, cErr := PrepareActionConfig(namespace, true)
	if cErr != nil {
		return cErr
	}

	_, err := uninstallRelease(actionConfig, releaseName)
	if err != nil {
		if errors.Is(err, driver.ErrReleaseNotFound) {
			return &models.ModelError{Code: 404, Message: "Release not found: " + err.Error()}
		}
		return &models.ModelError{Code: 500, Message: "Failed to uninstall release: " + err.Error()}
	}

	return nil
}

func GetHelmReleaseHistory(releaseName string, namespace string) ([]models.HelmReleaseHistory, *models.ModelError) {
	actionConfig, cErr := PrepareActionConfig(namespace, true)
	if cErr != nil {
		return nil, cErr
	}

	releases, err := getReleaseHistory(actionConfig, releaseName, 0)
	if err != nil {
		return nil, &models.ModelError{Code: 404, Message: "Failed to get release history"}
	}

	var helmReleases []models.HelmReleaseHistory
	for _, release := range releases {
		helmReleases = append(helmReleases, *GetReleaseHistoryData(release))
	}

	return helmReleases, nil
}

func RollbackHelmRelease(releaseName string, namespace string, version int) (*models.HelmRelease, *models.ModelError) {
	actionConfig, cErr := PrepareActionConfig(namespace, true)
	if cErr != nil {
		return nil, cErr
	}

	err := rollbackRelease(actionConfig, releaseName, version)
	if err != nil {
		return nil, &models.ModelError{Code: 500, Message: "Failed to rollback release: " + err.Error()}
	}

	release, err := getRelease(actionConfig, releaseName)
	if err != nil {
		return nil, &models.ModelError{Code: 404, Message: "Failed to get release: " + err.Error()}
	}

	return GetReleaseData(release), nil
}
