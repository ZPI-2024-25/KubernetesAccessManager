package helm

import (
	"github.com/ZPI-2024-25/KubernetesAccessManager/models"
	helmclient "github.com/mittwald/go-helm-client"
	"helm.sh/helm/v3/pkg/action"
	"time"
)

func GetHelmRelease(releaseName string, namespace string) (*models.HelmRelease, *models.ModelError) {
	helmClient, err := GetHelmClient(namespace)
	if err != nil {
		return nil, &models.ModelError{Code: 500, Message: "Failed to get helm client"}
	}

	release, err := helmClient.GetRelease(releaseName)
	if err != nil {
		return nil, &models.ModelError{Code: 404, Message: "Release not found"}
	}

	return GetReleaseData(release), nil
}

func ListHelmReleases(namespace string) (*[]models.HelmRelease, *models.ModelError) {
	helmClient, err := GetHelmClient(namespace)
	if err != nil {
		return nil, &models.ModelError{Code: 500, Message: "Failed to get helm client"}
	}

	releases, err := helmClient.ListReleasesByStateMask(action.ListAll)
	if err != nil {
		return nil, &models.ModelError{Code: 500, Message: "Failed to list releases"}
	}

	var helmReleases []models.HelmRelease
	for _, release := range releases {
		helmReleases = append(helmReleases, *GetReleaseData(release))
	}

	return &helmReleases, nil
}

func UninstallHelmRelease(release string, namespace string) *models.ModelError {
	helmClient, err := GetHelmClient(namespace)
	if err != nil {
		return &models.ModelError{Code: 500, Message: "Failed to get helm client"}
	}

	err = helmClient.UninstallReleaseByName(release)
	if err != nil {
		return &models.ModelError{Code: 500, Message: "Failed to uninstall release"}
	}

	return nil
}

func GetHelmReleaseHistory(name string, namespace string) (*[]models.HelmReleaseHistory, *models.ModelError) {
	helmClient, err := GetHelmClient(namespace)
	if err != nil {
		return nil, &models.ModelError{Code: 500, Message: "Failed to get helm client"}
	}

	releases, err := helmClient.ListReleaseHistory(name, 0)
	if err != nil {
		return nil, &models.ModelError{Code: 404, Message: "Failed to get release history"}
	}

	var helmReleases []models.HelmReleaseHistory
	for _, release := range releases {
		helmReleases = append(helmReleases, *GetReleaseHistoryData(release))
	}

	return &helmReleases, nil
}

func RollbackHelmRelease(name string, namespace string, version int) (*models.HelmRelease, *models.ModelError) {
	helmClient, err := GetHelmClient(namespace)
	if err != nil {
		return nil, &models.ModelError{Code: 500, Message: "Failed to get helm client"}
	}

	release, err := helmClient.GetRelease(name)
	if err != nil {
		return nil, &models.ModelError{Code: 404, Message: "Release not found"}
	}

	chartSpec := helmclient.ChartSpec{
		ReleaseName: release.Name,
		ChartName:   release.Chart.Name(),
		Namespace:   release.Namespace,
		UpgradeCRDs: true,
		Wait:        true,
		Timeout:     time.Duration(30) * time.Second,
	}

	err = helmClient.RollbackRelease(&chartSpec)
	if err != nil {
		return nil, &models.ModelError{Code: 500, Message: "Failed to rollback release: " + err.Error()}
	}

	return GetHelmRelease(name, namespace)

}
