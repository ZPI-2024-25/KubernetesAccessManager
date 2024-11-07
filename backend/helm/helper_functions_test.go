package helm

import (
	"github.com/ZPI-2024-25/KubernetesAccessManager/models"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"helm.sh/helm/v3/pkg/chart"
	"helm.sh/helm/v3/pkg/release"
	helmtime "helm.sh/helm/v3/pkg/time"
)

func TestGetReleaseData(t *testing.T) {
	inputRelease := &release.Release{
		Name:      "my-release",
		Namespace: "test-namespace",
		Chart: &chart.Chart{
			Metadata: &chart.Metadata{
				Name:       "my-chart",
				Version:    "1.2.3",
				AppVersion: "4.5.6",
			},
		},
		Info: &release.Info{
			Status:       release.StatusDeployed,
			LastDeployed: helmtime.Time{Time: time.Date(2023, time.April, 1, 12, 0, 0, 0, time.UTC)},
		},
		Version: 7,
	}

	expectedHelmRelease := &models.HelmRelease{
		Name:       "my-release",
		Namespace:  "test-namespace",
		Chart:      "my-chart-1.2.3",
		Status:     "deployed",
		Updated:    time.Date(2023, time.April, 1, 12, 0, 0, 0, time.UTC),
		Revision:   "7",
		AppVersion: "4.5.6",
	}

	result := getReleaseData(inputRelease)

	assert.Equal(t, expectedHelmRelease, result)
}

func TestGetReleaseHistoryData(t *testing.T) {
	inputReleaseHistory := &release.Release{
		Name:      "my-release",
		Namespace: "test-namespace",
		Chart: &chart.Chart{
			Metadata: &chart.Metadata{
				Name:       "my-chart",
				Version:    "1.2.3",
				AppVersion: "4.5.6",
			},
		},
		Info: &release.Info{
			Status:       release.StatusSuperseded,
			Description:  "Upgrade complete",
			LastDeployed: helmtime.Time{Time: time.Date(2023, time.May, 10, 15, 30, 0, 0, time.UTC)},
		},
		Version: 8,
	}

	expectedHistory := &models.HelmReleaseHistory{
		AppVersion:  "4.5.6",
		Description: "Upgrade complete",
		Updated:     time.Date(2023, time.May, 10, 15, 30, 0, 0, time.UTC),
		Chart:       "my-chart-1.2.3",
		Revision:    8,
		Status:      "superseded",
	}

	result := getReleaseHistoryData(inputReleaseHistory)

	assert.Equal(t, expectedHistory, result)
}
