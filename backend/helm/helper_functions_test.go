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

func TestGetReleaseData_NilRelease(t *testing.T) {
	result := getReleaseData(nil)
	assert.Nil(t, result)
}

func TestGetReleaseData_NilChart(t *testing.T) {
	inputRelease := &release.Release{
		Name:      "my-release",
		Namespace: "test-namespace",
		Chart:     nil,
		Info: &release.Info{
			Status:       release.StatusDeployed,
			LastDeployed: helmtime.Time{Time: time.Now()},
		},
		Version: 1,
	}

	result := getReleaseData(inputRelease)

	expected := &models.HelmRelease{
		Name:       "my-release",
		Namespace:  "test-namespace",
		Chart:      "-",
		Status:     "deployed",
		Updated:    inputRelease.Info.LastDeployed.Time,
		Revision:   "1",
		AppVersion: "",
	}

	assert.Equal(t, expected, result)
}

func TestGetReleaseData_NilInfo(t *testing.T) {
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
		Info:    nil,
		Version: 1,
	}

	result := getReleaseData(inputRelease)

	expected := &models.HelmRelease{
		Name:       "my-release",
		Namespace:  "test-namespace",
		Chart:      "my-chart-1.2.3",
		Status:     "",
		Updated:    time.Time{},
		Revision:   "1",
		AppVersion: "4.5.6",
	}

	assert.Equal(t, expected, result)
}

func TestGetReleaseData_NilChartMetadata(t *testing.T) {
	inputRelease := &release.Release{
		Name:      "my-release",
		Namespace: "test-namespace",
		Chart: &chart.Chart{
			Metadata: nil,
		},
		Info: &release.Info{
			Status:       release.StatusDeployed,
			LastDeployed: helmtime.Time{Time: time.Now()},
		},
		Version: 1,
	}

	result := getReleaseData(inputRelease)

	expected := &models.HelmRelease{
		Name:       "my-release",
		Namespace:  "test-namespace",
		Chart:      "-",
		Status:     "deployed",
		Updated:    inputRelease.Info.LastDeployed.Time,
		Revision:   "1",
		AppVersion: "",
	}

	assert.Equal(t, expected, result)
}

func TestGetReleaseData_EmptyChartNameAndVersion(t *testing.T) {
	inputRelease := &release.Release{
		Name:      "my-release",
		Namespace: "test-namespace",
		Chart: &chart.Chart{
			Metadata: &chart.Metadata{
				Name:       "",
				Version:    "",
				AppVersion: "4.5.6",
			},
		},
		Info: &release.Info{
			Status:       release.StatusDeployed,
			LastDeployed: helmtime.Time{Time: time.Now()},
		},
		Version: 1,
	}

	result := getReleaseData(inputRelease)

	expected := &models.HelmRelease{
		Name:       "my-release",
		Namespace:  "test-namespace",
		Chart:      "-",
		Status:     "deployed",
		Updated:    inputRelease.Info.LastDeployed.Time,
		Revision:   "1",
		AppVersion: "4.5.6",
	}

	assert.Equal(t, expected, result)
}

func TestGetReleaseData_ZeroVersion(t *testing.T) {
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
			LastDeployed: helmtime.Time{Time: time.Now()},
		},
		Version: 0,
	}

	result := getReleaseData(inputRelease)

	expected := &models.HelmRelease{
		Name:       "my-release",
		Namespace:  "test-namespace",
		Chart:      "my-chart-1.2.3",
		Status:     "deployed",
		Updated:    inputRelease.Info.LastDeployed.Time,
		Revision:   "0",
		AppVersion: "4.5.6",
	}

	assert.Equal(t, expected, result)
}

func TestGetReleaseData_ZeroUpdatedTime(t *testing.T) {
	zeroTime := time.Time{}
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
			LastDeployed: helmtime.Time{Time: zeroTime},
		},
		Version: 1,
	}

	result := getReleaseData(inputRelease)

	expected := &models.HelmRelease{
		Name:       "my-release",
		Namespace:  "test-namespace",
		Chart:      "my-chart-1.2.3",
		Status:     "deployed",
		Updated:    zeroTime,
		Revision:   "1",
		AppVersion: "4.5.6",
	}

	assert.Equal(t, expected, result)
}

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

func TestGetReleaseHistoryData_NilReleaseHistory(t *testing.T) {
	result := getReleaseHistoryData(nil)
	assert.Nil(t, result)
}

func TestGetReleaseHistoryData_NilChart(t *testing.T) {
	inputReleaseHistory := &release.Release{
		Name:      "my-release",
		Namespace: "test-namespace",
		Chart:     nil,
		Info: &release.Info{
			Status:       release.StatusSuperseded,
			Description:  "Upgrade complete",
			LastDeployed: helmtime.Time{Time: time.Now()},
		},
		Version: 8,
	}

	result := getReleaseHistoryData(inputReleaseHistory)

	expected := &models.HelmReleaseHistory{
		AppVersion:  "",
		Description: "Upgrade complete",
		Updated:     inputReleaseHistory.Info.LastDeployed.Time,
		Chart:       "-",
		Revision:    8,
		Status:      "superseded",
	}

	assert.Equal(t, expected, result)
}

func TestGetReleaseHistoryData_NilInfo(t *testing.T) {
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
		Info:    nil,
		Version: 8,
	}

	result := getReleaseHistoryData(inputReleaseHistory)

	expected := &models.HelmReleaseHistory{
		AppVersion:  "4.5.6",
		Description: "",
		Updated:     time.Time{},
		Chart:       "my-chart-1.2.3",
		Revision:    8,
		Status:      "",
	}

	assert.Equal(t, expected, result)
}

func TestGetReleaseHistoryData_NilChartMetadata(t *testing.T) {
	inputReleaseHistory := &release.Release{
		Name:      "my-release",
		Namespace: "test-namespace",
		Chart: &chart.Chart{
			Metadata: nil,
		},
		Info: &release.Info{
			Status:       release.StatusSuperseded,
			Description:  "Upgrade complete",
			LastDeployed: helmtime.Time{Time: time.Now()},
		},
		Version: 8,
	}

	result := getReleaseHistoryData(inputReleaseHistory)

	expected := &models.HelmReleaseHistory{
		AppVersion:  "",
		Description: "Upgrade complete",
		Updated:     inputReleaseHistory.Info.LastDeployed.Time,
		Chart:       "-",
		Revision:    8,
		Status:      "superseded",
	}

	assert.Equal(t, expected, result)
}

func TestGetReleaseHistoryData_EmptyChartNameAndVersion(t *testing.T) {
	inputReleaseHistory := &release.Release{
		Name:      "my-release",
		Namespace: "test-namespace",
		Chart: &chart.Chart{
			Metadata: &chart.Metadata{
				Name:       "",
				Version:    "",
				AppVersion: "",
			},
		},
		Info: &release.Info{
			Status:       release.StatusSuperseded,
			Description:  "Upgrade complete",
			LastDeployed: helmtime.Time{Time: time.Now()},
		},
		Version: 8,
	}

	result := getReleaseHistoryData(inputReleaseHistory)

	expected := &models.HelmReleaseHistory{
		AppVersion:  "",
		Description: "Upgrade complete",
		Updated:     inputReleaseHistory.Info.LastDeployed.Time,
		Chart:       "-",
		Revision:    8,
		Status:      "superseded",
	}

	assert.Equal(t, expected, result)
}

func TestGetReleaseHistoryData_ZeroVersion(t *testing.T) {
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
			LastDeployed: helmtime.Time{Time: time.Now()},
		},
		Version: 0,
	}

	result := getReleaseHistoryData(inputReleaseHistory)

	expected := &models.HelmReleaseHistory{
		AppVersion:  "4.5.6",
		Description: "Upgrade complete",
		Updated:     inputReleaseHistory.Info.LastDeployed.Time,
		Chart:       "my-chart-1.2.3",
		Revision:    0,
		Status:      "superseded",
	}

	assert.Equal(t, expected, result)
}

func TestGetReleaseHistoryData_ZeroUpdatedTime(t *testing.T) {
	zeroTime := time.Time{}
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
			LastDeployed: helmtime.Time{Time: zeroTime},
		},
		Version: 8,
	}

	result := getReleaseHistoryData(inputReleaseHistory)

	expected := &models.HelmReleaseHistory{
		AppVersion:  "4.5.6",
		Description: "Upgrade complete",
		Updated:     zeroTime,
		Chart:       "my-chart-1.2.3",
		Revision:    8,
		Status:      "superseded",
	}

	assert.Equal(t, expected, result)
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
