package helm

import (
	"fmt"
	"github.com/ZPI-2024-25/KubernetesAccessManager/models"
)

func GetHelmRelease(releaseName string, namespace string) (*models.HelmRelease, *models.ModelError) {
	helmClient, err := GetHelmClient()
	if err != nil {
		return nil, &models.ModelError{Code: 500, Message: "Failed to get helm client"}
	}

	release, err := helmClient.GetRelease(releaseName)
	if err != nil {
		return nil, &models.ModelError{Code: 404, Message: "Release not found"}
	}

	var helmRelease models.HelmRelease
	helmRelease.Name = release.Name
	helmRelease.Namespace = release.Namespace
	helmRelease.Chart = fmt.Sprintf("%s-%s", release.Chart.Name(), release.Chart.Metadata.Version)
	helmRelease.Status = release.Info.Status.String()
	helmRelease.Updated = release.Info.LastDeployed.Time
	helmRelease.Revision = fmt.Sprintf("%d", release.Version)
	helmRelease.AppVersion = release.Chart.AppVersion()

	return &helmRelease, nil
}
