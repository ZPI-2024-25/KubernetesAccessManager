package helm

import (
	"fmt"
	"github.com/ZPI-2024-25/KubernetesAccessManager/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"helm.sh/helm/v3/pkg/release"
	"testing"
)

type MockActionConfigGetter struct {
	mock.Mock
}

func (m *MockActionConfigGetter) Get(namespace string, useDefaultNamespace bool) (ActionConfigInterface, *models.ModelError) {
	args := m.Called(namespace, useDefaultNamespace)
	var actionConfig ActionConfigInterface
	if res := args.Get(0); res != nil {
		actionConfig = res.(ActionConfigInterface)
	}
	var modelErr *models.ModelError
	if err := args.Get(1); err != nil {
		modelErr = err.(*models.ModelError)
	}
	return actionConfig, modelErr
}

type MockActionConfig struct {
	mock.Mock
}

func (m *MockActionConfig) getRelease(name string) (*release.Release, error) {
	args := m.Called(name)
	var rel *release.Release
	if res := args.Get(0); res != nil {
		rel = res.(*release.Release)
	}
	return rel, args.Error(1)
}

func (m *MockActionConfig) listReleases(allNamespaces bool) ([]*release.Release, error) {
	args := m.Called(allNamespaces)
	var rels []*release.Release
	if res := args.Get(0); res != nil {
		rels = res.([]*release.Release)
	}
	return rels, args.Error(1)
}

func (m *MockActionConfig) rollbackRelease(name string, version int) error {
	args := m.Called(name, version)
	return args.Error(0)
}

func (m *MockActionConfig) uninstallRelease(name string) (*release.UninstallReleaseResponse, error) {
	args := m.Called(name)
	var resp *release.UninstallReleaseResponse
	if res := args.Get(0); res != nil {
		resp = res.(*release.UninstallReleaseResponse)
	}
	return resp, args.Error(1)
}

func (m *MockActionConfig) getReleaseHistory(name string, max int) ([]*release.Release, error) {
	args := m.Called(name, max)
	var history []*release.Release
	if res := args.Get(0); res != nil {
		history = res.([]*release.Release)
	}
	return history, args.Error(1)
}

func TestGetHelmRelease(t *testing.T) {
	tests := []struct {
		name            string
		releaseName     string
		namespace       string
		mockRelease     *release.Release
		mockConfigError *models.ModelError
		mockError       error
		expectedError   bool
		expectedCode    int
		expectedMsg     string
	}{
		{
			name:          "Release Not Found",
			releaseName:   "non-existent-release",
			namespace:     "test-namespace",
			mockRelease:   nil,
			mockError:     fmt.Errorf("release not found"),
			expectedError: true,
			expectedCode:  404,
			expectedMsg:   "Release not found",
		},
		{
			name:          "Success",
			releaseName:   "test-release",
			namespace:     "test-namespace",
			mockRelease:   &release.Release{Name: "test-release"},
			mockError:     nil,
			expectedError: false,
			expectedCode:  0,
			expectedMsg:   "",
		},
		{
			name:            "Config Error",
			releaseName:     "test-release",
			namespace:       "test-namespace",
			mockRelease:     nil,
			mockConfigError: &models.ModelError{Code: 500, Message: "Failed to get cluster config"},
			mockError:       nil,
			expectedError:   true,
			expectedCode:    500,
			expectedMsg:     "Failed to get cluster config",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockActionConfigGetter := new(MockActionConfigGetter)
			mockActionConfig := new(MockActionConfig)

			mockActionConfigGetter.On("Get", tt.namespace, false).Return(mockActionConfig, tt.mockConfigError)
			if tt.mockConfigError == nil {
				mockActionConfig.On("getRelease", tt.releaseName).Return(tt.mockRelease, tt.mockError)
			}

			result, err := GetHelmRelease(tt.releaseName, tt.namespace, mockActionConfigGetter.Get)

			if tt.expectedError {
				assert.NotNil(t, err)
				assert.Contains(t, err.Message, tt.expectedMsg)
				assert.Equal(t, tt.expectedCode, int(err.Code))
			} else {
				assert.Nil(t, err)
				assert.NotNil(t, result)
				assert.Equal(t, tt.releaseName, result.Name)
			}

			mockActionConfigGetter.AssertExpectations(t)
			mockActionConfig.AssertExpectations(t)
		})
	}
}

func TestListHelmReleases(t *testing.T) {
	tests := []struct {
		name            string
		namespace       string
		allNamespaces   bool
		mockReleases    []*release.Release
		mockConfigError *models.ModelError
		mockError       error
		expectedError   bool
		expectedCode    int
		expectedLen     int
		expectedMsg     string
	}{
		{
			name:            "Success With Namespace",
			namespace:       "test-namespace",
			allNamespaces:   false,
			mockReleases:    []*release.Release{{Name: "test-release-1"}, {Name: "test-release-2"}},
			mockConfigError: nil,
			mockError:       nil,
			expectedError:   false,
			expectedCode:    0,
			expectedLen:     2,
			expectedMsg:     "",
		},
		{
			name:            "Success All Namespaces",
			namespace:       "",
			allNamespaces:   true,
			mockReleases:    []*release.Release{{Name: "test-release-1"}, {Name: "test-release-2"}},
			mockConfigError: nil,
			mockError:       nil,
			expectedError:   false,
			expectedCode:    0,
			expectedLen:     2,
			expectedMsg:     "",
		},
		{
			name:            "Error Listing Releases",
			namespace:       "test-namespace",
			allNamespaces:   false,
			mockReleases:    nil,
			mockConfigError: nil,
			mockError:       fmt.Errorf("failed to list releases"),
			expectedError:   true,
			expectedCode:    500,
			expectedLen:     0,
			expectedMsg:     "Failed to list releases",
		},
		{
			name:            "Config Error",
			namespace:       "test-namespace",
			allNamespaces:   false,
			mockReleases:    nil,
			mockConfigError: &models.ModelError{Code: 500, Message: "Failed to get cluster config"},
			mockError:       nil,
			expectedError:   true,
			expectedCode:    500,
			expectedLen:     0,
			expectedMsg:     "Failed to get cluster config",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockActionConfigGetter := new(MockActionConfigGetter)
			mockActionConfig := new(MockActionConfig)

			mockActionConfigGetter.On("Get", tt.namespace, false).Return(mockActionConfig, tt.mockConfigError)
			if tt.mockConfigError == nil {
				mockActionConfig.On("listReleases", tt.allNamespaces).Return(tt.mockReleases, tt.mockError)
			}

			result, err := ListHelmReleases(tt.namespace, mockActionConfigGetter.Get)

			if tt.expectedError {
				assert.NotNil(t, err)
				assert.Contains(t, err.Message, tt.expectedMsg)
				assert.Equal(t, tt.expectedCode, int(err.Code))
			} else {
				assert.Nil(t, err)
				assert.Len(t, result, tt.expectedLen)
			}

			mockActionConfigGetter.AssertExpectations(t)
			mockActionConfig.AssertExpectations(t)
		})
	}
}
