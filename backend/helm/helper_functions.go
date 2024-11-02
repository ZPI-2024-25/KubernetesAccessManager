package helm

import (
	"fmt"
	"github.com/ZPI-2024-25/KubernetesAccessManager/models"
	"helm.sh/helm/v3/pkg/release"
)

func GetReleaseData(release *release.Release) *models.HelmRelease {
	var helmRelease models.HelmRelease
	helmRelease.Name = release.Name
	helmRelease.Namespace = release.Namespace
	helmRelease.Chart = fmt.Sprintf("%s-%s", release.Chart.Name(), release.Chart.Metadata.Version)
	helmRelease.Status = release.Info.Status.String()
	helmRelease.Updated = release.Info.LastDeployed.Time
	helmRelease.Revision = fmt.Sprintf("%d", release.Version)
	helmRelease.AppVersion = release.Chart.AppVersion()
	return &helmRelease
}
