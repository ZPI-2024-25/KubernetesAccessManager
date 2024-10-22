package cluster

import (
	"context"
	"errors"
	"testing"

	"github.com/ZPI-2024-25/KubernetesAccessManager/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/client-go/dynamic"
)

type MockResourceInterface struct {
	dynamic.ResourceInterface
	mock.Mock
	ReturnedGetValue *unstructured.Unstructured
	ReturnedGetError error
}

func (m *MockResourceInterface) Get(ctx context.Context, name string, options metav1.GetOptions, subresources ...string) (*unstructured.Unstructured, error) {
	return m.ReturnedGetValue, m.ReturnedGetError
}

func TestGetResourceError(t *testing.T) {

	mockGetResourceI := new(mock.Mock)
	mockGetResourceI.On("func1", "Pod", "default", "default").Return(&MockResourceInterface{}, &models.ModelError{Code: 404, Message: "Not found"})
	getResourceI = func(resourceType string, namespace string, emptyNamespace string) (dynamic.ResourceInterface, *models.ModelError) {
		args := mockGetResourceI.Called(resourceType, namespace, emptyNamespace)
		return args.Get(0).(dynamic.ResourceInterface), args.Get(1).(*models.ModelError)
	}

	t.Run("Test GetResource", func(t *testing.T) {
		_, err := GetResource("Pod", "default", "default")
		assert.NotNil(t, err)
		assert.Equal(t, &models.ModelError{Code: 404, Message: "Not found"}, err)
	})
}

func TestGetResourceSuccess(t *testing.T) {

	mockGetResourceI := new(mock.Mock)
	mockGetResourceI.On("func1", mock.Anything, mock.Anything, mock.Anything).
		Return(&MockResourceInterface{ReturnedGetValue: &unstructured.Unstructured{
				Object: map[string]interface{}{
					"key": "value",
				},
			}}, nil)
	getResourceI = func(resourceType string, namespace string, emptyNamespace string) (dynamic.ResourceInterface, *models.ModelError) {
		args := mockGetResourceI.Called(resourceType, namespace, emptyNamespace)
		return args.Get(0).(dynamic.ResourceInterface), nil
	}

	t.Run("Test GetResource", func(t *testing.T) {
		result, err := GetResource("Pod", "default", "default")
		assert.Nil(t, err)
		expectedResourceDetails := &models.ResourceDetails{
			ResourceDetails: func() *interface{} {
				obj := &unstructured.Unstructured{
					Object: map[string]interface{}{
						"key": "value",
					},
				}
				var i interface{} = obj
				return &i
			}(),
		}

		expectedObj := (*expectedResourceDetails.ResourceDetails).(*unstructured.Unstructured)
		resultObj := (*result.ResourceDetails).(*unstructured.Unstructured)
		for key, value := range expectedObj.Object {
			assert.Equal(t, value, resultObj.Object[key])
		}
	})
}

func TestGetResourceErrorFromGet(t *testing.T) {
	mockGetResourceI := new(mock.Mock)
	mockGetResourceI.On("func1", mock.Anything, mock.Anything, mock.Anything).
		Return(&MockResourceInterface{ReturnedGetValue: &unstructured.Unstructured{
				Object: map[string]interface{}{
					"key": "value",
				},
			}, ReturnedGetError: errors.New("error")}, nil)
	getResourceI = func(resourceType string, namespace string, emptyNamespace string) (dynamic.ResourceInterface, *models.ModelError) {
		args := mockGetResourceI.Called(resourceType, namespace, emptyNamespace)
		return args.Get(0).(dynamic.ResourceInterface), nil
	}
	t.Run("Test GetResourceErrorFromGet", func(t *testing.T) {
		res, err := GetResource("Pod", "default", "default")

		assert.NotNil(t, err)
		assert.EqualValues(t, err.Code, 500)
		assert.EqualValues(t, err.Message, "Internal server error: error")
		assert.EqualValues(t, res, models.ResourceDetails{})
	})
}
