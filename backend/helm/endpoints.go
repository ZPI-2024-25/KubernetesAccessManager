package helm

import (
	"github.com/ZPI-2024-25/KubernetesAccessManager/models"
	"helm.sh/helm/v3/pkg/action"
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
