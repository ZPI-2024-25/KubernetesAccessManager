package helm

import (
	"fmt"
	"github.com/ZPI-2024-25/KubernetesAccessManager/cluster"
	"github.com/ZPI-2024-25/KubernetesAccessManager/models"
	"helm.sh/helm/v3/pkg/action"
	"helm.sh/helm/v3/pkg/release"
	"time"
)

func getReleaseData(release *release.Release) *models.HelmRelease {
	if release == nil {
		return nil
	}

	name := release.Name
	namespace := release.Namespace
	var chartName, chartVersion, appVersion string
	if release.Chart != nil {
		if release.Chart.Metadata != nil {
			chartName = release.Chart.Name()
			chartVersion = release.Chart.Metadata.Version
			appVersion = release.Chart.AppVersion()
		}
	}
	chart := fmt.Sprintf("%s-%s", chartName, chartVersion)
	if chartName == "" && chartVersion == "" {
		chart = "-"
	}

	var status string
	var updated time.Time
	if release.Info != nil {
		status = release.Info.Status.String()
		updated = release.Info.LastDeployed.Time
	}

	revision := fmt.Sprintf("%d", release.Version)

	return &models.HelmRelease{
		Name:       name,
		Namespace:  namespace,
		Chart:      chart,
		Status:     status,
		Updated:    updated,
		Revision:   revision,
		AppVersion: appVersion,
	}
}

func getReleaseHistoryData(releaseHistory *release.Release) *models.HelmReleaseHistory {
	if releaseHistory == nil {
		return nil
	}

	var chartName, chartVersion, appVersion string
	if releaseHistory.Chart != nil {
		if releaseHistory.Chart.Metadata != nil {
			chartName = releaseHistory.Chart.Name()
			chartVersion = releaseHistory.Chart.Metadata.Version
			appVersion = releaseHistory.Chart.AppVersion()
		}
	}
	chart := fmt.Sprintf("%s-%s", chartName, chartVersion)
	if chartName == "" && chartVersion == "" {
		chart = "-"
	}

	var description, status string
	var updated time.Time
	if releaseHistory.Info != nil {
		description = releaseHistory.Info.Description
		status = releaseHistory.Info.Status.String()
		updated = releaseHistory.Info.LastDeployed.Time
	}

	revision := releaseHistory.Version

	return &models.HelmReleaseHistory{
		AppVersion:  appVersion,
		Description: description,
		Updated:     updated,
		Chart:       chart,
		Revision:    int32(revision),
		Status:      status,
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
