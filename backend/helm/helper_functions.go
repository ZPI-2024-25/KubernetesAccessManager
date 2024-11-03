package helm

import (
	"fmt"
	"github.com/ZPI-2024-25/KubernetesAccessManager/cluster"
	"github.com/ZPI-2024-25/KubernetesAccessManager/models"
	"helm.sh/helm/v3/pkg/action"
	"helm.sh/helm/v3/pkg/release"
)

func getReleaseData(release *release.Release) *models.HelmRelease {
	return &models.HelmRelease{
		Name:       release.Name,
		Namespace:  release.Namespace,
		Chart:      fmt.Sprintf("%s-%s", release.Chart.Name(), release.Chart.Metadata.Version),
		Status:     release.Info.Status.String(),
		Updated:    release.Info.LastDeployed.Time,
		Revision:   fmt.Sprintf("%d", release.Version),
		AppVersion: release.Chart.AppVersion(),
	}
}

func getReleaseHistoryData(releaseHistory *release.Release) *models.HelmReleaseHistory {
	return &models.HelmReleaseHistory{
		AppVersion:  releaseHistory.Chart.AppVersion(),
		Description: releaseHistory.Info.Description,
		Updated:     releaseHistory.Info.LastDeployed.Time,
		Chart:       fmt.Sprintf("%s-%s", releaseHistory.Chart.Name(), releaseHistory.Chart.Metadata.Version),
		Revision:    int32(releaseHistory.Version),
		Status:      releaseHistory.Info.Status.String(),
	}
}

func prepareActionConfig(namespace string, useDefaultNamespace bool) (*action.Configuration, *models.ModelError) {
	config, err := cluster.GetConfig()
	if err != nil {
		return nil, &models.ModelError{Code: 500, Message: "Failed to get cluster config"}
	}

	if namespace == "" && useDefaultNamespace {
		namespace = "default"
	}

	actionConfig, err := getActionConfig(config, namespace)
	if err != nil {
		return nil, &models.ModelError{Code: 500, Message: "Failed to create Helm action configuration: " + err.Error()}
	}

	return actionConfig, nil
}
