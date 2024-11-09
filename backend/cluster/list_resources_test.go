package cluster

import (
	"context"
	"errors"
	"github.com/ZPI-2024-25/KubernetesAccessManager/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/client-go/dynamic"
	"testing"
)

type MockListResourceInterface struct {
	dynamic.ResourceInterface
	mock.Mock
	ReturnedList  *unstructured.UnstructuredList
	ReturnedError error
}

func (m *MockListResourceInterface) List(ctx context.Context, opts metav1.ListOptions) (*unstructured.UnstructuredList, error) {
	return m.ReturnedList, m.ReturnedError
}

func MockUnstructuredList() *unstructured.UnstructuredList {
	return &unstructured.UnstructuredList{
		Items: []unstructured.Unstructured{
			{
				Object: map[string]interface{}{
					"metadata": map[string]interface{}{
						"name":      "resource1",
						"namespace": "validNamespace",
					},
				},
			},
			{
				Object: map[string]interface{}{
					"metadata": map[string]interface{}{
						"name":      "resource2",
						"namespace": "validNamespace",
					},
				},
			},
		},
	}
}

func TestListResourcesError(t *testing.T) {
	mockGetResourceI := new(mock.Mock)
	expectedModelError := &models.ModelError{Code: 404, Message: "Not found"}
	mockGetResourceI.On("func1", "Pod", "validNamespace", "").Return(&MockListResourceInterface{}, expectedModelError)

	getResourceI := func(resourceType string, namespace string, emptyNamespace string) (dynamic.ResourceInterface, *models.ModelError) {
		args := mockGetResourceI.Called(resourceType, namespace, emptyNamespace)
		return args.Get(0).(dynamic.ResourceInterface), args.Get(1).(*models.ModelError)
	}

	t.Run("Test ListResources Error", func(t *testing.T) {
		result, err := ListResources("Pod", "validNamespace", getResourceI)
		assert.NotNil(t, err)
		assert.Equal(t, expectedModelError, err)
		assert.Equal(t, models.ResourceList{}, result)
	})
}

func TestListResourcesSuccess(t *testing.T) {
	mockGetResourceI := new(mock.Mock)
	mockResourceInterface := &MockListResourceInterface{
		ReturnedList: MockUnstructuredList(),
	}
	mockGetResourceI.On("func1", "Pod", "validNamespace", "").
		Return(mockResourceInterface, nil)

	getResourceI := func(resourceType string, namespace string, emptyNamespace string) (dynamic.ResourceInterface, *models.ModelError) {
		args := mockGetResourceI.Called(resourceType, namespace, emptyNamespace)
		return args.Get(0).(dynamic.ResourceInterface), nil
	}

	t.Run("Test ListResources Success", func(t *testing.T) {
		result, err := ListResources("Pod", "validNamespace", getResourceI)
		assert.Nil(t, err)
		assert.Equal(t, 2, len(result.ResourceList))
		for _, resource := range result.ResourceList {
			assert.Equal(t, "validNamespace", resource.Namespace)
		}
	})
}

func TestListResourcesErrorFromList(t *testing.T) {
	mockGetResourceI := new(mock.Mock)
	expectedError := errors.New("failed to list resources")
	mockResourceInterface := &MockListResourceInterface{
		ReturnedError: expectedError,
	}
	mockGetResourceI.On("func1", "Pod", "validNamespace", "").
		Return(mockResourceInterface, nil)

	getResourceI := func(resourceType string, namespace string, emptyNamespace string) (dynamic.ResourceInterface, *models.ModelError) {
		args := mockGetResourceI.Called(resourceType, namespace, emptyNamespace)
		return args.Get(0).(dynamic.ResourceInterface), nil
	}

	t.Run("Test ListResources Error from List", func(t *testing.T) {
		result, err := ListResources("Pod", "validNamespace", getResourceI)
		expectedModelError := &models.ModelError{Code: 500, Message: "Internal server error: " + expectedError.Error()}
		assert.NotNil(t, err)
		assert.Equal(t, expectedModelError, err)
		assert.Equal(t, models.ResourceList{}, result)
	})
}

func TestExtractActive(t *testing.T) {
	tests := []struct {
		name                      string
		resource                  unstructured.Unstructured
		resourceType              string
		expectedActive            string
		transposedColumnsOverride map[string][]string
	}{
		{
			name:           "ResourceType not in transposedResourceListColumns['active']",
			resource:       unstructured.Unstructured{},
			resourceType:   "ConfigMap",
			expectedActive: "",
		},
		{
			name: "Status does not exist",
			resource: unstructured.Unstructured{
				Object: map[string]interface{}{},
			},
			resourceType:   "CronJob",
			expectedActive: "",
		},
		{
			name: "Active field missing in status",
			resource: unstructured.Unstructured{
				Object: map[string]interface{}{
					"status": map[string]interface{}{},
				},
			},
			resourceType:   "CronJob",
			expectedActive: "0",
		},
		{
			name: "Active field is empty slice",
			resource: unstructured.Unstructured{
				Object: map[string]interface{}{
					"status": map[string]interface{}{
						"active": []interface{}{},
					},
				},
			},
			resourceType:   "CronJob",
			expectedActive: "0",
		},
		{
			name: "Active field has elements",
			resource: unstructured.Unstructured{
				Object: map[string]interface{}{
					"status": map[string]interface{}{
						"active": []interface{}{1, 2, 3},
					},
				},
			},
			resourceType:   "CronJob",
			expectedActive: "3",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resourceDetailsTruncated := &models.ResourceListResourceList{}

			extractActive(tt.resource, tt.resourceType, resourceDetailsTruncated)

			if resourceDetailsTruncated.Active != tt.expectedActive {
				t.Errorf("Expected Active to be '%s', got '%s'", tt.expectedActive, resourceDetailsTruncated.Active)
			}
		})
	}
}

func TestExtractAge(t *testing.T) {
	tests := []struct {
		name                      string
		resource                  unstructured.Unstructured
		resourceType              string
		expectedAge               string
		transposedColumnsOverride map[string][]string
	}{
		{
			name:         "ResourceType not in transposedResourceListColumns['age']",
			resource:     unstructured.Unstructured{},
			resourceType: "NotExistingResourceType",
			expectedAge:  "",
		},
		{
			name: "CreationTimestamp missing in metadata",
			resource: unstructured.Unstructured{
				Object: map[string]interface{}{
					"metadata": map[string]interface{}{},
				},
			},
			resourceType: "Pod",
			expectedAge:  "",
		},
		{
			name: "CreationTimestamp present",
			resource: unstructured.Unstructured{
				Object: map[string]interface{}{
					"metadata": map[string]interface{}{
						"creationTimestamp": "2023-05-01T10:00:00Z",
					},
				},
			},
			resourceType: "Pod",
			expectedAge:  "2023-05-01T10:00:00Z",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resourceDetailsTruncated := &models.ResourceListResourceList{}

			extractAge(tt.resource, tt.resourceType, resourceDetailsTruncated)

			if resourceDetailsTruncated.Age != tt.expectedAge {
				t.Errorf("Expected Age to be '%s', got '%s'", tt.expectedAge, resourceDetailsTruncated.Age)
			}
		})
	}
}

func TestExtractBindings(t *testing.T) {
	tests := []struct {
		name           string
		resource       unstructured.Unstructured
		resourceType   string
		expectedResult string
	}{
		{
			name:           "ResourceType not in transposedResourceListColumns['bindings']",
			resource:       unstructured.Unstructured{},
			resourceType:   "ConfigMap",
			expectedResult: "",
		},
		{
			name: "Subjects field does not exist",
			resource: unstructured.Unstructured{
				Object: map[string]interface{}{},
			},
			resourceType:   "ClusterRoleBinding",
			expectedResult: "",
		},
		{
			name: "Subjects field is empty",
			resource: unstructured.Unstructured{
				Object: map[string]interface{}{
					"subjects": []interface{}{},
				},
			},
			resourceType:   "ClusterRoleBinding",
			expectedResult: "",
		},
		{
			name: "Subjects field has entries",
			resource: unstructured.Unstructured{
				Object: map[string]interface{}{
					"subjects": []interface{}{
						map[string]interface{}{"name": "user1"},
						map[string]interface{}{"name": "user2"},
					},
				},
			},
			resourceType:   "ClusterRoleBinding",
			expectedResult: "user1, user2",
		},
		{
			name: "Subjects field has entries with missing names",
			resource: unstructured.Unstructured{
				Object: map[string]interface{}{
					"subjects": []interface{}{
						map[string]interface{}{"name": "user1"},
						map[string]interface{}{"kind": "User"},
					},
				},
			},
			resourceType:   "ClusterRoleBinding",
			expectedResult: "user1",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resourceDetailsTruncated := &models.ResourceListResourceList{}

			extractBindings(tt.resource, tt.resourceType, resourceDetailsTruncated)

			if resourceDetailsTruncated.Bindings != tt.expectedResult {
				t.Errorf("Expected Bindings to be '%s', got '%s'", tt.expectedResult, resourceDetailsTruncated.Bindings)
			}
		})
	}
}

func TestExtractCapacity(t *testing.T) {
	tests := []struct {
		name           string
		resource       unstructured.Unstructured
		resourceType   string
		expectedResult string
	}{
		{
			name:           "ResourceType not in transposedResourceListColumns['capacity']",
			resource:       unstructured.Unstructured{},
			resourceType:   "ConfigMap",
			expectedResult: "",
		},
		{
			name: "Spec field does not exist",
			resource: unstructured.Unstructured{
				Object: map[string]interface{}{},
			},
			resourceType:   "PersistentVolume",
			expectedResult: "",
		},
		{
			name: "Capacity field does not exist",
			resource: unstructured.Unstructured{
				Object: map[string]interface{}{
					"spec": map[string]interface{}{},
				},
			},
			resourceType:   "PersistentVolume",
			expectedResult: "",
		},
		{
			name: "Capacity field exists",
			resource: unstructured.Unstructured{
				Object: map[string]interface{}{
					"spec": map[string]interface{}{
						"capacity": map[string]interface{}{
							"storage": "10Gi",
						},
					},
				},
			},
			resourceType:   "PersistentVolume",
			expectedResult: "10Gi",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resourceDetailsTruncated := &models.ResourceListResourceList{}

			extractCapacity(tt.resource, tt.resourceType, resourceDetailsTruncated)

			if resourceDetailsTruncated.Capacity != tt.expectedResult {
				t.Errorf("Expected Capacity to be '%s', got '%s'", tt.expectedResult, resourceDetailsTruncated.Capacity)
			}
		})
	}
}

func TestExtractClaim(t *testing.T) {
	tests := []struct {
		name           string
		resource       unstructured.Unstructured
		resourceType   string
		expectedResult string
	}{
		{
			name:           "ResourceType not in transposedResourceListColumns['claim']",
			resource:       unstructured.Unstructured{},
			resourceType:   "ConfigMap",
			expectedResult: "",
		},
		{
			name: "Spec field does not exist",
			resource: unstructured.Unstructured{
				Object: map[string]interface{}{},
			},
			resourceType:   "PersistentVolume",
			expectedResult: "",
		},
		{
			name: "ClaimRef field does not exist",
			resource: unstructured.Unstructured{
				Object: map[string]interface{}{
					"spec": map[string]interface{}{},
				},
			},
			resourceType:   "PersistentVolume",
			expectedResult: "",
		},
		{
			name: "ClaimRef field exists with name",
			resource: unstructured.Unstructured{
				Object: map[string]interface{}{
					"spec": map[string]interface{}{
						"claimRef": map[string]interface{}{
							"name": "my-pvc",
						},
					},
				},
			},
			resourceType:   "PersistentVolume",
			expectedResult: "my-pvc",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resourceDetailsTruncated := &models.ResourceListResourceList{}

			extractClaim(tt.resource, tt.resourceType, resourceDetailsTruncated)

			if resourceDetailsTruncated.Claim != tt.expectedResult {
				t.Errorf("Expected Claim to be '%s', got '%s'", tt.expectedResult, resourceDetailsTruncated.Claim)
			}
		})
	}
}

func TestExtractClusterIp(t *testing.T) {
	tests := []struct {
		name           string
		resource       unstructured.Unstructured
		resourceType   string
		expectedResult string
	}{
		{
			name:           "ResourceType not in transposedResourceListColumns['cluster_ip']",
			resource:       unstructured.Unstructured{},
			resourceType:   "ConfigMap",
			expectedResult: "",
		},
		{
			name: "Spec field does not exist",
			resource: unstructured.Unstructured{
				Object: map[string]interface{}{},
			},
			resourceType:   "Service",
			expectedResult: "",
		},
		{
			name: "ClusterIP field does not exist",
			resource: unstructured.Unstructured{
				Object: map[string]interface{}{
					"spec": map[string]interface{}{},
				},
			},
			resourceType:   "Service",
			expectedResult: "",
		},
		{
			name: "ClusterIP field exists",
			resource: unstructured.Unstructured{
				Object: map[string]interface{}{
					"spec": map[string]interface{}{
						"clusterIP": "10.0.0.1",
					},
				},
			},
			resourceType:   "Service",
			expectedResult: "10.0.0.1",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resourceDetailsTruncated := &models.ResourceListResourceList{}

			extractClusterIp(tt.resource, tt.resourceType, resourceDetailsTruncated)

			if resourceDetailsTruncated.ClusterIp != tt.expectedResult {
				t.Errorf("Expected ClusterIp to be '%s', got '%s'", tt.expectedResult, resourceDetailsTruncated.ClusterIp)
			}
		})
	}
}

func TestExtractCompletions(t *testing.T) {
	tests := []struct {
		name           string
		resource       unstructured.Unstructured
		resourceType   string
		expectedResult string
	}{
		{
			name:           "ResourceType not in transposedResourceListColumns['completions']",
			resource:       unstructured.Unstructured{},
			resourceType:   "ConfigMap",
			expectedResult: "",
		},
		{
			name: "Spec and Status do not exist",
			resource: unstructured.Unstructured{
				Object: map[string]interface{}{},
			},
			resourceType:   "Job",
			expectedResult: "0/1",
		},
		{
			name: "Spec exists with completions",
			resource: unstructured.Unstructured{
				Object: map[string]interface{}{
					"spec": map[string]interface{}{
						"completions": int64(5),
					},
				},
			},
			resourceType:   "Job",
			expectedResult: "0/5",
		},
		{
			name: "Status exists with succeeded",
			resource: unstructured.Unstructured{
				Object: map[string]interface{}{
					"status": map[string]interface{}{
						"succeeded": int64(3),
					},
				},
			},
			resourceType:   "Job",
			expectedResult: "3/1",
		},
		{
			name: "Spec and Status exist with completions and succeeded",
			resource: unstructured.Unstructured{
				Object: map[string]interface{}{
					"spec": map[string]interface{}{
						"completions": int64(5),
					},
					"status": map[string]interface{}{
						"succeeded": int64(3),
					},
				},
			},
			resourceType:   "Job",
			expectedResult: "3/5",
		},
		{
			name: "Spec completions not int64",
			resource: unstructured.Unstructured{
				Object: map[string]interface{}{
					"spec": map[string]interface{}{
						"completions": "5",
					},
				},
			},
			resourceType:   "Job",
			expectedResult: "0/1",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resourceDetailsTruncated := &models.ResourceListResourceList{}

			extractCompletions(tt.resource, tt.resourceType, resourceDetailsTruncated)

			if resourceDetailsTruncated.Completions != tt.expectedResult {
				t.Errorf("Expected Completions to be '%s', got '%s'", tt.expectedResult, resourceDetailsTruncated.Completions)
			}
		})
	}
}

func TestExtractConditions(t *testing.T) {
	tests := []struct {
		name           string
		resource       unstructured.Unstructured
		resourceType   string
		expectedResult string
	}{
		{
			name:           "ResourceType not in transposedResourceListColumns['conditions']",
			resource:       unstructured.Unstructured{},
			resourceType:   "ConfigMap",
			expectedResult: "",
		},
		{
			name: "Status does not exist",
			resource: unstructured.Unstructured{
				Object: map[string]interface{}{},
			},
			resourceType:   "Deployment",
			expectedResult: "",
		},
		{
			name: "Deployment with conditions",
			resource: unstructured.Unstructured{
				Object: map[string]interface{}{
					"status": map[string]interface{}{
						"conditions": []interface{}{
							map[string]interface{}{
								"type":   "Available",
								"status": "True",
							},
							map[string]interface{}{
								"type":   "Progressing",
								"status": "False",
							},
						},
					},
				},
			},
			resourceType:   "Deployment",
			expectedResult: "Available",
		},
		{
			name: "Deployment with more conditions",
			resource: unstructured.Unstructured{
				Object: map[string]interface{}{
					"status": map[string]interface{}{
						"conditions": []interface{}{
							map[string]interface{}{
								"type":   "Available",
								"status": "True",
							},
							map[string]interface{}{
								"type":   "Progressing",
								"status": "True",
							},
						},
					},
				},
			},
			resourceType:   "Deployment",
			expectedResult: "Available, Progressing",
		},
		{
			name: "Node with conditions",
			resource: unstructured.Unstructured{
				Object: map[string]interface{}{
					"status": map[string]interface{}{
						"conditions": []interface{}{
							map[string]interface{}{
								"type":   "Ready",
								"status": "True",
							},
							map[string]interface{}{
								"type":   "DiskPressure",
								"status": "False",
							},
						},
					},
				},
			},
			resourceType:   "Node",
			expectedResult: "Ready",
		},
		{
			name: "Job with conditions",
			resource: unstructured.Unstructured{
				Object: map[string]interface{}{
					"status": map[string]interface{}{
						"conditions": []interface{}{
							map[string]interface{}{
								"type": "Complete",
							},
						},
					},
				},
			},
			resourceType:   "Job",
			expectedResult: "Complete",
		},
		{
			name: "Job with more conditions",
			resource: unstructured.Unstructured{
				Object: map[string]interface{}{
					"status": map[string]interface{}{
						"conditions": []interface{}{
							map[string]interface{}{
								"type": "SuccessCriteriaMet",
							},
							map[string]interface{}{
								"type": "Complete",
							},
						},
					},
				},
			},
			resourceType:   "Job",
			expectedResult: "SuccessCriteriaMet",
		},
		{
			name: "Job with unknown condition",
			resource: unstructured.Unstructured{
				Object: map[string]interface{}{
					"status": map[string]interface{}{
						"conditions": []interface{}{
							"invalid",
						},
					},
				},
			},
			resourceType:   "Job",
			expectedResult: "Unknown",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resourceDetailsTruncated := &models.ResourceListResourceList{}

			extractConditions(tt.resource, tt.resourceType, resourceDetailsTruncated)

			if resourceDetailsTruncated.Conditions != tt.expectedResult {
				t.Errorf("Expected Conditions to be '%s', got '%s'", tt.expectedResult, resourceDetailsTruncated.Conditions)
			}
		})
	}
}

func TestExtractContainers(t *testing.T) {
	tests := []struct {
		name           string
		resource       unstructured.Unstructured
		resourceType   string
		expectedResult string
	}{
		{
			name:           "ResourceType not in transposedResourceListColumns['containers']",
			resource:       unstructured.Unstructured{},
			resourceType:   "ConfigMap",
			expectedResult: "",
		},
		{
			name: "Status does not exist",
			resource: unstructured.Unstructured{
				Object: map[string]interface{}{},
			},
			resourceType:   "Pod",
			expectedResult: "",
		},
		{
			name: "ContainerStatuses does not exist",
			resource: unstructured.Unstructured{
				Object: map[string]interface{}{
					"status": map[string]interface{}{},
				},
			},
			resourceType:   "Pod",
			expectedResult: "",
		},
		{
			name: "No containers",
			resource: unstructured.Unstructured{
				Object: map[string]interface{}{
					"status": map[string]interface{}{
						"containerStatuses": []interface{}{},
					},
				},
			},
			resourceType:   "Pod",
			expectedResult: "0/0",
		},
		{
			name: "Containers with readiness",
			resource: unstructured.Unstructured{
				Object: map[string]interface{}{
					"status": map[string]interface{}{
						"containerStatuses": []interface{}{
							map[string]interface{}{
								"ready": true,
							},
							map[string]interface{}{
								"ready": false,
							},
							map[string]interface{}{
								"ready": true,
							},
						},
					},
				},
			},
			resourceType:   "Pod",
			expectedResult: "2/3",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resourceDetailsTruncated := &models.ResourceListResourceList{}

			extractContainers(tt.resource, tt.resourceType, resourceDetailsTruncated)

			if resourceDetailsTruncated.Containers != tt.expectedResult {
				t.Errorf("Expected Containers to be '%s', got '%s'", tt.expectedResult, resourceDetailsTruncated.Containers)
			}
		})
	}
}

func TestExtractControlledBy(t *testing.T) {
	tests := []struct {
		name           string
		resource       unstructured.Unstructured
		resourceType   string
		expectedResult string
	}{
		{
			name:           "ResourceType not in transposedResourceListColumns['controlled_by']",
			resource:       unstructured.Unstructured{},
			resourceType:   "ConfigMap",
			expectedResult: "",
		},
		{
			name: "OwnerReferences does not exist",
			resource: unstructured.Unstructured{
				Object: map[string]interface{}{
					"metadata": map[string]interface{}{},
				},
			},
			resourceType:   "Pod",
			expectedResult: "",
		},
		{
			name: "OwnerReferences exist",
			resource: unstructured.Unstructured{
				Object: map[string]interface{}{
					"metadata": map[string]interface{}{
						"ownerReferences": []interface{}{
							map[string]interface{}{
								"kind": "ReplicaSet",
								"name": "my-replicaset",
							},
							map[string]interface{}{
								"kind": "Deployment",
								"name": "my-deployment",
							},
						},
					},
				},
			},
			resourceType:   "Pod",
			expectedResult: "ReplicaSet:my-replicaset, Deployment:my-deployment",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resourceDetailsTruncated := &models.ResourceListResourceList{}

			extractControlledBy(tt.resource, tt.resourceType, resourceDetailsTruncated)

			if resourceDetailsTruncated.ControlledBy != tt.expectedResult {
				t.Errorf("Expected ControlledBy to be '%s', got '%s'", tt.expectedResult, resourceDetailsTruncated.ControlledBy)
			}
		})
	}
}

func TestExtractCurrent(t *testing.T) {
	tests := []struct {
		name           string
		resource       unstructured.Unstructured
		resourceType   string
		expectedResult string
	}{
		{
			name:           "ResourceType not in transposedResourceListColumns['current']",
			resource:       unstructured.Unstructured{},
			resourceType:   "ConfigMap",
			expectedResult: "",
		},
		{
			name: "Status does not exist",
			resource: unstructured.Unstructured{
				Object: map[string]interface{}{},
			},
			resourceType:   "ReplicaSet",
			expectedResult: "0",
		},
		{
			name: "availableReplicas does not exist",
			resource: unstructured.Unstructured{
				Object: map[string]interface{}{
					"status": map[string]interface{}{},
				},
			},
			resourceType:   "ReplicaSet",
			expectedResult: "0",
		},
		{
			name: "availableReplicas exist",
			resource: unstructured.Unstructured{
				Object: map[string]interface{}{
					"status": map[string]interface{}{
						"availableReplicas": int64(3),
					},
				},
			},
			resourceType:   "ReplicaSet",
			expectedResult: "3",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resourceDetailsTruncated := &models.ResourceListResourceList{}

			extractCurrent(tt.resource, tt.resourceType, resourceDetailsTruncated)

			if resourceDetailsTruncated.Current != tt.expectedResult {
				t.Errorf("Expected Current to be '%s', got '%s'", tt.expectedResult, resourceDetailsTruncated.Current)
			}
		})
	}
}

func TestExtractDefault(t *testing.T) {
	tests := []struct {
		name           string
		resource       unstructured.Unstructured
		resourceType   string
		expectedResult string
	}{
		{
			name:           "ResourceType not in transposedResourceListColumns['default']",
			resource:       unstructured.Unstructured{},
			resourceType:   "ConfigMap",
			expectedResult: "",
		},
		{
			name: "Annotations do not exist",
			resource: unstructured.Unstructured{
				Object: map[string]interface{}{
					"metadata": map[string]interface{}{},
				},
			},
			resourceType:   "StorageClass",
			expectedResult: "No",
		},
		{
			name: "Is default class (new annotation)",
			resource: unstructured.Unstructured{
				Object: map[string]interface{}{
					"metadata": map[string]interface{}{
						"annotations": map[string]interface{}{
							"storageclass.kubernetes.io/is-default-class": "true",
						},
					},
				},
			},
			resourceType:   "StorageClass",
			expectedResult: "Yes",
		},
		{
			name: "Is default class (beta annotation)",
			resource: unstructured.Unstructured{
				Object: map[string]interface{}{
					"metadata": map[string]interface{}{
						"annotations": map[string]interface{}{
							"storageclass.beta.kubernetes.io/is-default-class": "true",
						},
					},
				},
			},
			resourceType:   "StorageClass",
			expectedResult: "Yes",
		},
		{
			name: "Not default class",
			resource: unstructured.Unstructured{
				Object: map[string]interface{}{
					"metadata": map[string]interface{}{
						"annotations": map[string]interface{}{
							"storageclass.kubernetes.io/is-default-class": "false",
						},
					},
				},
			},
			resourceType:   "StorageClass",
			expectedResult: "No",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resourceDetailsTruncated := &models.ResourceListResourceList{}

			extractDefault(tt.resource, tt.resourceType, resourceDetailsTruncated)

			if resourceDetailsTruncated.Default_ != tt.expectedResult {
				t.Errorf("Expected Default to be '%s', got '%s'", tt.expectedResult, resourceDetailsTruncated.Default_)
			}
		})
	}
}

func TestExtractDesired(t *testing.T) {
	tests := []struct {
		name           string
		resource       unstructured.Unstructured
		resourceType   string
		expectedResult string
	}{
		{
			name:           "ResourceType not in transposedResourceListColumns['desired']",
			resource:       unstructured.Unstructured{},
			resourceType:   "ConfigMap",
			expectedResult: "",
		},
		{
			name: "Spec does not exist",
			resource: unstructured.Unstructured{
				Object: map[string]interface{}{},
			},
			resourceType:   "ReplicaSet",
			expectedResult: "",
		},
		{
			name: "Replicas does not exist in spec",
			resource: unstructured.Unstructured{
				Object: map[string]interface{}{
					"spec": map[string]interface{}{},
				},
			},
			resourceType:   "ReplicaSet",
			expectedResult: "",
		},
		{
			name: "Replicas exist and are int64",
			resource: unstructured.Unstructured{
				Object: map[string]interface{}{
					"spec": map[string]interface{}{
						"replicas": int64(5),
					},
				},
			},
			resourceType:   "ReplicaSet",
			expectedResult: "5",
		},
		{
			name: "Replicas exist but not int64",
			resource: unstructured.Unstructured{
				Object: map[string]interface{}{
					"spec": map[string]interface{}{
						"replicas": "5",
					},
				},
			},
			resourceType:   "ReplicaSet",
			expectedResult: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resourceDetailsTruncated := &models.ResourceListResourceList{}

			extractDesired(tt.resource, tt.resourceType, resourceDetailsTruncated)

			if resourceDetailsTruncated.Desired != tt.expectedResult {
				t.Errorf("Expected Desired to be '%s', got '%s'", tt.expectedResult, resourceDetailsTruncated.Desired)
			}
		})
	}
}

func TestExtractExternalIp(t *testing.T) {
	tests := []struct {
		name           string
		resource       unstructured.Unstructured
		resourceType   string
		expectedResult string
	}{
		{
			name:           "ResourceType not in transposedResourceListColumns['external_ip']",
			resource:       unstructured.Unstructured{},
			resourceType:   "ConfigMap",
			expectedResult: "",
		},
		{
			name: "ServiceType not found in spec",
			resource: unstructured.Unstructured{
				Object: map[string]interface{}{
					"spec": map[string]interface{}{},
				},
			},
			resourceType:   "Service",
			expectedResult: "-",
		},
		{
			name: "ServiceType is LoadBalancer with no ingress",
			resource: unstructured.Unstructured{
				Object: map[string]interface{}{
					"spec": map[string]interface{}{
						"type": "LoadBalancer",
					},
					"status": map[string]interface{}{},
				},
			},
			resourceType:   "Service",
			expectedResult: "<pending>",
		},
		{
			name: "ServiceType is LoadBalancer with ingress IPs",
			resource: unstructured.Unstructured{
				Object: map[string]interface{}{
					"spec": map[string]interface{}{
						"type": "LoadBalancer",
					},
					"status": map[string]interface{}{
						"loadBalancer": map[string]interface{}{
							"ingress": []interface{}{
								map[string]interface{}{
									"ip": "192.168.1.1",
								},
								map[string]interface{}{
									"hostname": "example.com",
								},
							},
						},
					},
				},
			},
			resourceType:   "Service",
			expectedResult: "192.168.1.1,example.com",
		},
		{
			name: "ServiceType is LoadBalancer with empty ingress",
			resource: unstructured.Unstructured{
				Object: map[string]interface{}{
					"spec": map[string]interface{}{
						"type": "LoadBalancer",
					},
					"status": map[string]interface{}{
						"loadBalancer": map[string]interface{}{
							"ingress": []interface{}{},
						},
					},
				},
			},
			resourceType:   "Service",
			expectedResult: "<pending>",
		},
		{
			name: "ServiceType is NodePort with externalIPs",
			resource: unstructured.Unstructured{
				Object: map[string]interface{}{
					"spec": map[string]interface{}{
						"type":        "NodePort",
						"externalIPs": []interface{}{"10.0.0.1", "10.0.0.2"},
					},
				},
			},
			resourceType:   "Service",
			expectedResult: "10.0.0.1,10.0.0.2",
		},
		{
			name: "ServiceType is NodePort without externalIPs",
			resource: unstructured.Unstructured{
				Object: map[string]interface{}{
					"spec": map[string]interface{}{
						"type": "NodePort",
					},
				},
			},
			resourceType:   "Service",
			expectedResult: "-",
		},
		{
			name: "ServiceType is ClusterIP with externalIPs",
			resource: unstructured.Unstructured{
				Object: map[string]interface{}{
					"spec": map[string]interface{}{
						"type":        "ClusterIP",
						"externalIPs": []interface{}{"10.0.0.3"},
					},
				},
			},
			resourceType:   "Service",
			expectedResult: "10.0.0.3",
		},
		{
			name: "ServiceType is ClusterIP without externalIPs",
			resource: unstructured.Unstructured{
				Object: map[string]interface{}{
					"spec": map[string]interface{}{
						"type": "ClusterIP",
					},
				},
			},
			resourceType:   "Service",
			expectedResult: "-",
		},
		{
			name: "ServiceType is Unknown",
			resource: unstructured.Unstructured{
				Object: map[string]interface{}{
					"spec": map[string]interface{}{
						"type": "UnknownType",
					},
				},
			},
			resourceType:   "Service",
			expectedResult: "<unknown>",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resourceDetailsTruncated := &models.ResourceListResourceList{}

			extractExternalIp(tt.resource, tt.resourceType, resourceDetailsTruncated)

			if resourceDetailsTruncated.ExternalIp != tt.expectedResult {
				t.Errorf("Expected ExternalIp to be '%s', got '%s'", tt.expectedResult, resourceDetailsTruncated.ExternalIp)
			}
		})
	}
}

func TestExtractGroup(t *testing.T) {
	tests := []struct {
		name           string
		resource       unstructured.Unstructured
		resourceType   string
		expectedResult string
	}{
		{
			name:           "ResourceType not in transposedResourceListColumns['group']",
			resource:       unstructured.Unstructured{},
			resourceType:   "ConfigMap",
			expectedResult: "",
		},
		{
			name: "Spec does not exist",
			resource: unstructured.Unstructured{
				Object: map[string]interface{}{},
			},
			resourceType:   "CustomResourceDefinition",
			expectedResult: "",
		},
		{
			name: "group does not exist in spec",
			resource: unstructured.Unstructured{
				Object: map[string]interface{}{
					"spec": map[string]interface{}{},
				},
			},
			resourceType:   "CustomResourceDefinition",
			expectedResult: "",
		},
		{
			name: "group exists as string",
			resource: unstructured.Unstructured{
				Object: map[string]interface{}{
					"spec": map[string]interface{}{
						"group": "my-group",
					},
				},
			},
			resourceType:   "CustomResourceDefinition",
			expectedResult: "my-group",
		},
		{
			name: "group exists but not string",
			resource: unstructured.Unstructured{
				Object: map[string]interface{}{
					"spec": map[string]interface{}{
						"group": 123,
					},
				},
			},
			resourceType:   "CustomResourceDefinition",
			expectedResult: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resourceDetailsTruncated := &models.ResourceListResourceList{}

			extractGroup(tt.resource, tt.resourceType, resourceDetailsTruncated)

			if resourceDetailsTruncated.Group != tt.expectedResult {
				t.Errorf("Expected Group to be '%s', got '%s'", tt.expectedResult, resourceDetailsTruncated.Group)
			}
		})
	}
}

func TestExtractKeys(t *testing.T) {
	tests := []struct {
		name           string
		resource       unstructured.Unstructured
		resourceType   string
		expectedResult string
	}{
		{
			name:           "ResourceType not in transposedResourceListColumns['keys']",
			resource:       unstructured.Unstructured{},
			resourceType:   "Service",
			expectedResult: "",
		},
		{
			name: "data and binaryData do not exist",
			resource: unstructured.Unstructured{
				Object: map[string]interface{}{},
			},
			resourceType:   "ConfigMap",
			expectedResult: "",
		},
		{
			name: "data exists with keys",
			resource: unstructured.Unstructured{
				Object: map[string]interface{}{
					"data": map[string]interface{}{
						"key1": "value1",
						"key2": "value2",
					},
				},
			},
			resourceType:   "ConfigMap",
			expectedResult: "key1, key2",
		},
		{
			name: "binaryData exists with keys",
			resource: unstructured.Unstructured{
				Object: map[string]interface{}{
					"binaryData": map[string]interface{}{
						"key3": []byte{0x03, 0x04},
						"key4": []byte{0x01, 0x02},
					},
				},
			},
			resourceType:   "ConfigMap",
			expectedResult: "key3, key4",
		},
		{
			name: "both data and binaryData exist with keys",
			resource: unstructured.Unstructured{
				Object: map[string]interface{}{
					"data": map[string]interface{}{
						"key1": "value1",
					},
					"binaryData": map[string]interface{}{
						"key3": []byte{0x01, 0x02},
					},
				},
			},
			resourceType:   "ConfigMap",
			expectedResult: "key1, key3",
		},
		{
			name: "data and binaryData empty",
			resource: unstructured.Unstructured{
				Object: map[string]interface{}{
					"data":       map[string]interface{}{},
					"binaryData": map[string]interface{}{},
				},
			},
			resourceType:   "ConfigMap",
			expectedResult: "",
		},
		{
			name: "secret with both data and binaryData exist with keys",
			resource: unstructured.Unstructured{
				Object: map[string]interface{}{
					"data": map[string]interface{}{
						"key1": "value1",
					},
					"binaryData": map[string]interface{}{
						"key3": []byte{0x01, 0x02},
					},
				},
			},
			resourceType:   "Secret",
			expectedResult: "key1, key3",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resourceDetailsTruncated := &models.ResourceListResourceList{}

			extractKeys(tt.resource, tt.resourceType, resourceDetailsTruncated)

			if resourceDetailsTruncated.Keys != tt.expectedResult {
				t.Errorf("Expected Keys to be '%s', got '%s'", tt.expectedResult, resourceDetailsTruncated.Keys)
			}
		})
	}
}

func TestExtractLabels(t *testing.T) {
	tests := []struct {
		name           string
		resource       unstructured.Unstructured
		resourceType   string
		expectedResult string
	}{
		{
			name:           "ResourceType not in transposedResourceListColumns['labels']",
			resource:       unstructured.Unstructured{},
			resourceType:   "ConfigMap",
			expectedResult: "",
		},
		{
			name: "Labels do not exist",
			resource: unstructured.Unstructured{
				Object: map[string]interface{}{
					"metadata": map[string]interface{}{},
				},
			},
			resourceType:   "Secret",
			expectedResult: "",
		},
		{
			name: "Labels exist",
			resource: unstructured.Unstructured{
				Object: map[string]interface{}{
					"metadata": map[string]interface{}{
						"labels": map[string]interface{}{
							"app":     "my-app",
							"version": "v1",
						},
					},
				},
			},
			resourceType:   "Secret",
			expectedResult: "app=my-app, version=v1",
		},
		{
			name: "namespace Labels exist",
			resource: unstructured.Unstructured{
				Object: map[string]interface{}{
					"metadata": map[string]interface{}{
						"labels": map[string]interface{}{
							"app":     "my-app",
							"version": "v1",
						},
					},
				},
			},
			resourceType:   "Namespace",
			expectedResult: "app=my-app, version=v1",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resourceDetailsTruncated := &models.ResourceListResourceList{}

			extractLabels(tt.resource, tt.resourceType, resourceDetailsTruncated)

			if resourceDetailsTruncated.Labels != tt.expectedResult {
				t.Errorf("Expected Labels to be '%s', got '%s'", tt.expectedResult, resourceDetailsTruncated.Labels)
			}
		})
	}
}

func TestExtractLastSchedule(t *testing.T) {
	tests := []struct {
		name           string
		resource       unstructured.Unstructured
		resourceType   string
		expectedResult string
	}{
		{
			name:           "ResourceType not in transposedResourceListColumns['last_schedule']",
			resource:       unstructured.Unstructured{},
			resourceType:   "Job",
			expectedResult: "",
		},
		{
			name: "Status does not exist",
			resource: unstructured.Unstructured{
				Object: map[string]interface{}{},
			},
			resourceType:   "CronJob",
			expectedResult: "",
		},
		{
			name: "lastScheduleTime does not exist",
			resource: unstructured.Unstructured{
				Object: map[string]interface{}{
					"status": map[string]interface{}{},
				},
			},
			resourceType:   "CronJob",
			expectedResult: "",
		},
		{
			name: "lastScheduleTime exists",
			resource: unstructured.Unstructured{
				Object: map[string]interface{}{
					"status": map[string]interface{}{
						"lastScheduleTime": "2023-05-01T10:00:00Z",
					},
				},
			},
			resourceType:   "CronJob",
			expectedResult: "2023-05-01T10:00:00Z",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resourceDetailsTruncated := &models.ResourceListResourceList{}

			extractLastSchedule(tt.resource, tt.resourceType, resourceDetailsTruncated)

			if resourceDetailsTruncated.LastSchedule != tt.expectedResult {
				t.Errorf("Expected LastSchedule to be '%s', got '%s'", tt.expectedResult, resourceDetailsTruncated.LastSchedule)
			}
		})
	}
}

func TestExtractLoadbalancers(t *testing.T) {
	tests := []struct {
		name           string
		resource       unstructured.Unstructured
		resourceType   string
		expectedResult string
	}{
		{
			name:           "ResourceType not in transposedResourceListColumns['loadbalancers']",
			resource:       unstructured.Unstructured{},
			resourceType:   "ConfigMap",
			expectedResult: "",
		},
		{
			name: "Status does not exist",
			resource: unstructured.Unstructured{
				Object: map[string]interface{}{},
			},
			resourceType:   "Ingress",
			expectedResult: "",
		},
		{
			name: "loadBalancer does not exist in status",
			resource: unstructured.Unstructured{
				Object: map[string]interface{}{
					"status": map[string]interface{}{},
				},
			},
			resourceType:   "Ingress",
			expectedResult: "",
		},
		{
			name: "Ingress exists with IPs",
			resource: unstructured.Unstructured{
				Object: map[string]interface{}{
					"status": map[string]interface{}{
						"loadBalancer": map[string]interface{}{
							"ingress": []interface{}{
								map[string]interface{}{
									"ip": "192.168.1.1",
								},
							},
						},
					},
				},
			},
			resourceType:   "Ingress",
			expectedResult: "192.168.1.1",
		},
		{
			name: "Ingress exists but empty",
			resource: unstructured.Unstructured{
				Object: map[string]interface{}{
					"status": map[string]interface{}{
						"loadBalancer": map[string]interface{}{
							"ingress": []interface{}{},
						},
					},
				},
			},
			resourceType:   "Ingress",
			expectedResult: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resourceDetailsTruncated := &models.ResourceListResourceList{}

			extractLoadbalancers(tt.resource, tt.resourceType, resourceDetailsTruncated)

			if resourceDetailsTruncated.Loadbalancers != tt.expectedResult {
				t.Errorf("Expected Loadbalancers to be '%s', got '%s'", tt.expectedResult, resourceDetailsTruncated.Loadbalancers)
			}
		})
	}
}

func TestExtractName(t *testing.T) {
	tests := []struct {
		name           string
		resource       unstructured.Unstructured
		resourceType   string
		expectedResult string
	}{
		{
			name:           "ResourceType not in transposedResourceListColumns['name']",
			resource:       unstructured.Unstructured{},
			resourceType:   "UnknownType",
			expectedResult: "",
		},
		{
			name: "Name does not exist in metadata",
			resource: unstructured.Unstructured{
				Object: map[string]interface{}{
					"metadata": map[string]interface{}{},
				},
			},
			resourceType:   "Pod",
			expectedResult: "",
		},
		{
			name: "Name exists",
			resource: unstructured.Unstructured{
				Object: map[string]interface{}{
					"metadata": map[string]interface{}{
						"name": "my-pod",
					},
				},
			},
			resourceType:   "Pod",
			expectedResult: "my-pod",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resourceDetailsTruncated := &models.ResourceListResourceList{}

			extractName(tt.resource, tt.resourceType, resourceDetailsTruncated)

			if resourceDetailsTruncated.Name != tt.expectedResult {
				t.Errorf("Expected Name to be '%s', got '%s'", tt.expectedResult, resourceDetailsTruncated.Name)
			}
		})
	}
}

func TestExtractNamespace(t *testing.T) {
	tests := []struct {
		name           string
		resource       unstructured.Unstructured
		resourceType   string
		expectedResult string
	}{
		{
			name:           "ResourceType not in transposedResourceListColumns['namespace']",
			resource:       unstructured.Unstructured{},
			resourceType:   "UnknownType",
			expectedResult: "",
		},
		{
			name: "Namespace does not exist in metadata",
			resource: unstructured.Unstructured{
				Object: map[string]interface{}{
					"metadata": map[string]interface{}{},
				},
			},
			resourceType:   "Service",
			expectedResult: "",
		},
		{
			name: "Namespace exists",
			resource: unstructured.Unstructured{
				Object: map[string]interface{}{
					"metadata": map[string]interface{}{
						"namespace": "default",
					},
				},
			},
			resourceType:   "Service",
			expectedResult: "default",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resourceDetailsTruncated := &models.ResourceListResourceList{}

			extractNamespace(tt.resource, tt.resourceType, resourceDetailsTruncated)

			if resourceDetailsTruncated.Namespace != tt.expectedResult {
				t.Errorf("Expected Namespace to be '%s', got '%s'", tt.expectedResult, resourceDetailsTruncated.Namespace)
			}
		})
	}
}

func TestExtractNode(t *testing.T) {
	tests := []struct {
		name           string
		resource       unstructured.Unstructured
		resourceType   string
		expectedResult string
	}{
		{
			name:           "ResourceType not in transposedResourceListColumns['node']",
			resource:       unstructured.Unstructured{},
			resourceType:   "ConfigMap",
			expectedResult: "",
		},
		{
			name: "Spec does not exist",
			resource: unstructured.Unstructured{
				Object: map[string]interface{}{},
			},
			resourceType:   "Pod",
			expectedResult: "",
		},
		{
			name: "nodeName does not exist",
			resource: unstructured.Unstructured{
				Object: map[string]interface{}{
					"spec": map[string]interface{}{},
				},
			},
			resourceType:   "Pod",
			expectedResult: "",
		},
		{
			name: "nodeName exists",
			resource: unstructured.Unstructured{
				Object: map[string]interface{}{
					"spec": map[string]interface{}{
						"nodeName": "node-1",
					},
				},
			},
			resourceType:   "Pod",
			expectedResult: "node-1",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resourceDetailsTruncated := &models.ResourceListResourceList{}

			extractNode(tt.resource, tt.resourceType, resourceDetailsTruncated)

			if resourceDetailsTruncated.Node != tt.expectedResult {
				t.Errorf("Expected Node to be '%s', got '%s'", tt.expectedResult, resourceDetailsTruncated.Node)
			}
		})
	}
}

func TestExtractNodeSelector(t *testing.T) {
	tests := []struct {
		name           string
		resource       unstructured.Unstructured
		resourceType   string
		expectedResult string
	}{
		{
			name:           "ResourceType not in transposedResourceListColumns['node_selector']",
			resource:       unstructured.Unstructured{},
			resourceType:   "ConfigMap",
			expectedResult: "",
		},
		{
			name: "Spec does not exist",
			resource: unstructured.Unstructured{
				Object: map[string]interface{}{},
			},
			resourceType:   "DaemonSet",
			expectedResult: "None",
		},
		{
			name: "NodeSelector does not exist",
			resource: unstructured.Unstructured{
				Object: map[string]interface{}{
					"spec": map[string]interface{}{
						"template": map[string]interface{}{
							"spec": map[string]interface{}{},
						},
					},
				},
			},
			resourceType:   "DaemonSet",
			expectedResult: "None",
		},
		{
			name: "NodeSelector is empty",
			resource: unstructured.Unstructured{
				Object: map[string]interface{}{
					"spec": map[string]interface{}{
						"template": map[string]interface{}{
							"spec": map[string]interface{}{
								"nodeSelector": map[string]interface{}{},
							},
						},
					},
				},
			},
			resourceType:   "DaemonSet",
			expectedResult: "None",
		},
		{
			name: "NodeSelector exists with labels",
			resource: unstructured.Unstructured{
				Object: map[string]interface{}{
					"spec": map[string]interface{}{
						"template": map[string]interface{}{
							"spec": map[string]interface{}{
								"nodeSelector": map[string]interface{}{
									"disktype": "ssd",
									"region":   "us-west",
								},
							},
						},
					},
				},
			},
			resourceType:   "DaemonSet",
			expectedResult: "disktype=ssd, region=us-west",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resourceDetailsTruncated := &models.ResourceListResourceList{}

			extractNodeSelector(tt.resource, tt.resourceType, resourceDetailsTruncated)

			if resourceDetailsTruncated.NodeSelector != tt.expectedResult {
				t.Errorf("Expected NodeSelector to be '%s', got '%s'", tt.expectedResult, resourceDetailsTruncated.NodeSelector)
			}
		})
	}
}

func TestExtractPods(t *testing.T) {
	tests := []struct {
		name           string
		resource       unstructured.Unstructured
		resourceType   string
		expectedResult string
	}{
		{
			name:           "ResourceType not in transposedResourceListColumns['pods']",
			resource:       unstructured.Unstructured{},
			resourceType:   "ConfigMap",
			expectedResult: "",
		},
		{
			name: "Status does not exist",
			resource: unstructured.Unstructured{
				Object: map[string]interface{}{},
			},
			resourceType:   "Deployment",
			expectedResult: "",
		},
		{
			name: "Deployment with replicas and unavailableReplicas",
			resource: unstructured.Unstructured{
				Object: map[string]interface{}{
					"status": map[string]interface{}{
						"replicas":            int64(5),
						"unavailableReplicas": int64(2),
					},
				},
			},
			resourceType:   "Deployment",
			expectedResult: "3/5",
		},
		{
			name: "Deployment with only replicas",
			resource: unstructured.Unstructured{
				Object: map[string]interface{}{
					"status": map[string]interface{}{
						"replicas": int64(5),
					},
				},
			},
			resourceType:   "Deployment",
			expectedResult: "5/5",
		},
		{
			name: "StatefulSet with availableReplicas and replicas",
			resource: unstructured.Unstructured{
				Object: map[string]interface{}{
					"status": map[string]interface{}{
						"availableReplicas": int64(3),
						"replicas":          int64(5),
					},
				},
			},
			resourceType:   "StatefulSet",
			expectedResult: "3/5",
		},
		{
			name: "StatefulSet with only replicas",
			resource: unstructured.Unstructured{
				Object: map[string]interface{}{
					"status": map[string]interface{}{
						"replicas": int64(4),
					},
				},
			},
			resourceType:   "StatefulSet",
			expectedResult: "0/4",
		},
		{
			name: "DaemonSet with numberReady and desiredNumberScheduled",
			resource: unstructured.Unstructured{
				Object: map[string]interface{}{
					"status": map[string]interface{}{
						"numberReady":            int64(5),
						"desiredNumberScheduled": int64(5),
					},
				},
			},
			resourceType:   "DaemonSet",
			expectedResult: "5/5",
		},
		{
			name: "DaemonSet missing numberReady",
			resource: unstructured.Unstructured{
				Object: map[string]interface{}{
					"status": map[string]interface{}{
						"desiredNumberScheduled": int64(3),
					},
				},
			},
			resourceType:   "DaemonSet",
			expectedResult: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resourceDetailsTruncated := &models.ResourceListResourceList{}

			extractPods(tt.resource, tt.resourceType, resourceDetailsTruncated)

			if resourceDetailsTruncated.Pods != tt.expectedResult {
				t.Errorf("Expected Pods to be '%s', got '%s'", tt.expectedResult, resourceDetailsTruncated.Pods)
			}
		})
	}
}

func TestExtractPorts(t *testing.T) {
	tests := []struct {
		name           string
		resource       unstructured.Unstructured
		resourceType   string
		expectedResult string
	}{
		{
			name:           "ResourceType not in transposedResourceListColumns['ports']",
			resource:       unstructured.Unstructured{},
			resourceType:   "ConfigMap",
			expectedResult: "",
		},
		{
			name: "Spec does not exist",
			resource: unstructured.Unstructured{
				Object: map[string]interface{}{},
			},
			resourceType:   "Service",
			expectedResult: "",
		},
		{
			name: "Ports does not exist",
			resource: unstructured.Unstructured{
				Object: map[string]interface{}{
					"spec": map[string]interface{}{},
				},
			},
			resourceType:   "Service",
			expectedResult: "",
		},
		{
			name: "Ports exist with various fields",
			resource: unstructured.Unstructured{
				Object: map[string]interface{}{
					"spec": map[string]interface{}{
						"ports": []interface{}{
							map[string]interface{}{
								"port":       int64(80),
								"targetPort": int64(8080),
								"protocol":   "TCP",
							},
							map[string]interface{}{
								"port":     int64(443),
								"protocol": "TCP",
							},
							map[string]interface{}{
								"port":     int64(53),
								"nodePort": int64(30053),
								"protocol": "UDP",
							},
							map[string]interface{}{
								"port": int64(22),
							},
						},
					},
				},
			},
			resourceType:   "Service",
			expectedResult: "80:8080/TCP, 443/TCP, 53:30053/UDP, 22",
		},
		{
			name: "Ports with missing protocol",
			resource: unstructured.Unstructured{
				Object: map[string]interface{}{
					"spec": map[string]interface{}{
						"ports": []interface{}{
							map[string]interface{}{
								"port": int64(80),
							},
						},
					},
				},
			},
			resourceType:   "Service",
			expectedResult: "80",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resourceDetailsTruncated := &models.ResourceListResourceList{}

			extractPorts(tt.resource, tt.resourceType, resourceDetailsTruncated)

			if resourceDetailsTruncated.Ports != tt.expectedResult {
				t.Errorf("Expected Ports to be '%s', got '%s'", tt.expectedResult, resourceDetailsTruncated.Ports)
			}
		})
	}
}

func TestExtractProvisioner(t *testing.T) {
	tests := []struct {
		name           string
		resource       unstructured.Unstructured
		resourceType   string
		expectedResult string
	}{
		{
			name:           "ResourceType not in transposedResourceListColumns['provisioner']",
			resource:       unstructured.Unstructured{},
			resourceType:   "ConfigMap",
			expectedResult: "",
		},
		{
			name: "Provisioner field exists",
			resource: unstructured.Unstructured{
				Object: map[string]interface{}{
					"provisioner": "kubernetes.io/aws-ebs",
				},
			},
			resourceType:   "StorageClass",
			expectedResult: "kubernetes.io/aws-ebs",
		},
		{
			name: "Provisioner field does not exist",
			resource: unstructured.Unstructured{
				Object: map[string]interface{}{},
			},
			resourceType:   "StorageClass",
			expectedResult: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resourceDetailsTruncated := &models.ResourceListResourceList{}

			extractProvisioner(tt.resource, tt.resourceType, resourceDetailsTruncated)

			if resourceDetailsTruncated.Provisioner != tt.expectedResult {
				t.Errorf("Expected Provisioner to be '%s', got '%s'", tt.expectedResult, resourceDetailsTruncated.Provisioner)
			}
		})
	}
}

func TestExtractQos(t *testing.T) {
	tests := []struct {
		name           string
		resource       unstructured.Unstructured
		resourceType   string
		expectedResult string
	}{
		{
			name:           "ResourceType not in transposedResourceListColumns['qos']",
			resource:       unstructured.Unstructured{},
			resourceType:   "ConfigMap",
			expectedResult: "",
		},
		{
			name: "Status does not exist",
			resource: unstructured.Unstructured{
				Object: map[string]interface{}{},
			},
			resourceType:   "Pod",
			expectedResult: "Unknown",
		},
		{
			name: "ResourceType in 'qos', status.qosClass not present",
			resource: unstructured.Unstructured{
				Object: map[string]interface{}{},
			},
			resourceType:   "Pod",
			expectedResult: "Unknown",
		},
		{
			name: "ResourceType in 'qos', status.qosClass present",
			resource: unstructured.Unstructured{
				Object: map[string]interface{}{
					"status": map[string]interface{}{
						"qosClass": "BestEffort",
					},
				},
			},
			resourceType:   "Pod",
			expectedResult: "BestEffort",
		},
		{
			name: "ResourceType in 'qos', status.qosClass present as non-string (simulate error)",
			resource: unstructured.Unstructured{
				Object: map[string]interface{}{
					"status": map[string]interface{}{
						"qosClass": 123,
					},
				},
			},
			resourceType:   "Pod",
			expectedResult: "Unknown",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resourceDetailsTruncated := &models.ResourceListResourceList{}

			extractQos(tt.resource, tt.resourceType, resourceDetailsTruncated)

			if resourceDetailsTruncated.Qos != tt.expectedResult {
				t.Errorf("Expected Qos to be '%s', got '%s'", tt.expectedResult, resourceDetailsTruncated.Qos)
			}
		})
	}
}

func TestExtractReady(t *testing.T) {
	tests := []struct {
		name           string
		resource       unstructured.Unstructured
		resourceType   string
		expectedResult string
	}{
		{
			name:           "ResourceType not in transposedResourceListColumns['ready']",
			resource:       unstructured.Unstructured{},
			resourceType:   "ConfigMap",
			expectedResult: "",
		},
		{
			name: "Status does not exist",
			resource: unstructured.Unstructured{
				Object: map[string]interface{}{},
			},
			resourceType:   "ReplicaSet",
			expectedResult: "0",
		},
		{
			name: "readyReplicas does not exist",
			resource: unstructured.Unstructured{
				Object: map[string]interface{}{
					"status": map[string]interface{}{},
				},
			},
			resourceType:   "ReplicaSet",
			expectedResult: "0",
		},
		{
			name: "readyReplicas exists",
			resource: unstructured.Unstructured{
				Object: map[string]interface{}{
					"status": map[string]interface{}{
						"readyReplicas": int64(3),
					},
				},
			},
			resourceType:   "ReplicaSet",
			expectedResult: "3",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resourceDetailsTruncated := &models.ResourceListResourceList{}

			extractReady(tt.resource, tt.resourceType, resourceDetailsTruncated)

			if resourceDetailsTruncated.Ready != tt.expectedResult {
				t.Errorf("Expected Ready to be '%s', got '%s'", tt.expectedResult, resourceDetailsTruncated.Ready)
			}
		})
	}
}

func TestExtractReclaimPolicy(t *testing.T) {
	tests := []struct {
		name           string
		resource       unstructured.Unstructured
		resourceType   string
		expectedResult string
	}{
		{
			name:           "ResourceType not in transposedResourceListColumns['reclaim_policy']",
			resource:       unstructured.Unstructured{},
			resourceType:   "ConfigMap",
			expectedResult: "",
		},
		{
			name: "reclaimPolicy field exists",
			resource: unstructured.Unstructured{
				Object: map[string]interface{}{
					"reclaimPolicy": "Retain",
				},
			},
			resourceType:   "StorageClass",
			expectedResult: "Retain",
		},
		{
			name: "reclaimPolicy field does not exist",
			resource: unstructured.Unstructured{
				Object: map[string]interface{}{},
			},
			resourceType:   "StorageClass",
			expectedResult: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resourceDetailsTruncated := &models.ResourceListResourceList{}

			extractReclaimPolicy(tt.resource, tt.resourceType, resourceDetailsTruncated)

			if resourceDetailsTruncated.ReclaimPolicy != tt.expectedResult {
				t.Errorf("Expected ReclaimPolicy to be '%s', got '%s'", tt.expectedResult, resourceDetailsTruncated.ReclaimPolicy)
			}
		})
	}
}

func TestExtractReplicas(t *testing.T) {
	tests := []struct {
		name           string
		resource       unstructured.Unstructured
		resourceType   string
		expectedResult string
	}{
		{
			name:           "ResourceType not in transposedResourceListColumns['replicas']",
			resource:       unstructured.Unstructured{},
			resourceType:   "ConfigMap",
			expectedResult: "",
		},
		{
			name: "Spec does not exist",
			resource: unstructured.Unstructured{
				Object: map[string]interface{}{},
			},
			resourceType:   "Deployment",
			expectedResult: "0",
		},
		{
			name: "Replicas field does not exist",
			resource: unstructured.Unstructured{
				Object: map[string]interface{}{
					"spec": map[string]interface{}{},
				},
			},
			resourceType:   "Deployment",
			expectedResult: "0",
		},
		{
			name: "Replicas field exists",
			resource: unstructured.Unstructured{
				Object: map[string]interface{}{
					"spec": map[string]interface{}{
						"replicas": int64(3),
					},
				},
			},
			resourceType:   "Deployment",
			expectedResult: "3",
		},
		{
			name: "Replicas field exists but not int64",
			resource: unstructured.Unstructured{
				Object: map[string]interface{}{
					"spec": map[string]interface{}{
						"replicas": "3",
					},
				},
			},
			resourceType:   "Deployment",
			expectedResult: "0",
		},
		{
			name: "StatefulSet with replicas",
			resource: unstructured.Unstructured{
				Object: map[string]interface{}{
					"spec": map[string]interface{}{
						"replicas": int64(5),
					},
				},
			},
			resourceType:   "StatefulSet",
			expectedResult: "5",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resourceDetailsTruncated := &models.ResourceListResourceList{}

			extractReplicas(tt.resource, tt.resourceType, resourceDetailsTruncated)

			if resourceDetailsTruncated.Replicas != tt.expectedResult {
				t.Errorf("Expected Replicas to be '%s', got '%s'", tt.expectedResult, resourceDetailsTruncated.Replicas)
			}
		})
	}
}

func TestExtractResource(t *testing.T) {
	tests := []struct {
		name           string
		resource       unstructured.Unstructured
		resourceType   string
		expectedResult string
	}{
		{
			name:           "ResourceType not in transposedResourceListColumns['resource']",
			resource:       unstructured.Unstructured{},
			resourceType:   "Deployment",
			expectedResult: "",
		},
		{
			name: "Spec does not exist",
			resource: unstructured.Unstructured{
				Object: map[string]interface{}{},
			},
			resourceType:   "CustomResourceDefinition",
			expectedResult: "",
		},
		{
			name: "Names field does not exist",
			resource: unstructured.Unstructured{
				Object: map[string]interface{}{
					"spec": map[string]interface{}{},
				},
			},
			resourceType:   "CustomResourceDefinition",
			expectedResult: "",
		},
		{
			name: "Singular field exists",
			resource: unstructured.Unstructured{
				Object: map[string]interface{}{
					"spec": map[string]interface{}{
						"names": map[string]interface{}{
							"singular": "myresource",
						},
					},
				},
			},
			resourceType:   "CustomResourceDefinition",
			expectedResult: "Myresource",
		},
		{
			name: "Singular field missing",
			resource: unstructured.Unstructured{
				Object: map[string]interface{}{
					"spec": map[string]interface{}{
						"names": map[string]interface{}{},
					},
				},
			},
			resourceType:   "CustomResourceDefinition",
			expectedResult: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resourceDetailsTruncated := &models.ResourceListResourceList{}

			extractResource(tt.resource, tt.resourceType, resourceDetailsTruncated)

			if resourceDetailsTruncated.Resource != tt.expectedResult {
				t.Errorf("Expected Resource to be '%s', got '%s'", tt.expectedResult, resourceDetailsTruncated.Resource)
			}
		})
	}
}

func TestExtractRestarts(t *testing.T) {
	tests := []struct {
		name           string
		resource       unstructured.Unstructured
		resourceType   string
		expectedResult string
	}{
		{
			name:           "ResourceType not in transposedResourceListColumns['restarts']",
			resource:       unstructured.Unstructured{},
			resourceType:   "Service",
			expectedResult: "",
		},
		{
			name: "Status does not exist",
			resource: unstructured.Unstructured{
				Object: map[string]interface{}{},
			},
			resourceType:   "Pod",
			expectedResult: "",
		},
		{
			name: "ContainerStatuses does not exist",
			resource: unstructured.Unstructured{
				Object: map[string]interface{}{
					"status": map[string]interface{}{},
				},
			},
			resourceType:   "Pod",
			expectedResult: "",
		},
		{
			name: "ContainerStatuses exist with restart counts",
			resource: unstructured.Unstructured{
				Object: map[string]interface{}{
					"status": map[string]interface{}{
						"containerStatuses": []interface{}{
							map[string]interface{}{
								"restartCount": int64(2),
							},
							map[string]interface{}{
								"restartCount": int64(1),
							},
						},
					},
				},
			},
			resourceType:   "Pod",
			expectedResult: "3",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resourceDetailsTruncated := &models.ResourceListResourceList{}

			extractRestarts(tt.resource, tt.resourceType, resourceDetailsTruncated)

			if resourceDetailsTruncated.Restarts != tt.expectedResult {
				t.Errorf("Expected Restarts to be '%s', got '%s'", tt.expectedResult, resourceDetailsTruncated.Restarts)
			}
		})
	}
}

func TestExtractRoles(t *testing.T) {
	tests := []struct {
		name           string
		resource       unstructured.Unstructured
		resourceType   string
		expectedResult string
	}{
		{
			name:           "ResourceType not in transposedResourceListColumns['roles']",
			resource:       unstructured.Unstructured{},
			resourceType:   "Pod",
			expectedResult: "",
		},
		{
			name: "Labels do not exist",
			resource: unstructured.Unstructured{
				Object: map[string]interface{}{
					"metadata": map[string]interface{}{},
				},
			},
			resourceType:   "Node",
			expectedResult: "",
		},
		{
			name: "No node-role labels",
			resource: unstructured.Unstructured{
				Object: map[string]interface{}{
					"metadata": map[string]interface{}{
						"labels": map[string]interface{}{
							"app": "my-app",
						},
					},
				},
			},
			resourceType:   "Node",
			expectedResult: "",
		},
		{
			name: "Node-role labels exist",
			resource: unstructured.Unstructured{
				Object: map[string]interface{}{
					"metadata": map[string]interface{}{
						"labels": map[string]interface{}{
							"node-role.kubernetes.io/master": "",
							"node-role.kubernetes.io/worker": "",
						},
					},
				},
			},
			resourceType:   "Node",
			expectedResult: "master, worker",
		},
		{
			name: "Node-role label with value",
			resource: unstructured.Unstructured{
				Object: map[string]interface{}{
					"metadata": map[string]interface{}{
						"labels": map[string]interface{}{
							"node-role.kubernetes.io/custom-role": "true",
						},
					},
				},
			},
			resourceType:   "Node",
			expectedResult: "custom-role",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resourceDetailsTruncated := &models.ResourceListResourceList{}

			extractRoles(tt.resource, tt.resourceType, resourceDetailsTruncated)

			if resourceDetailsTruncated.Roles != tt.expectedResult {
				t.Errorf("Expected Roles to be '%s', got '%s'", tt.expectedResult, resourceDetailsTruncated.Roles)
			}
		})
	}
}

func TestExtractSchedule(t *testing.T) {
	tests := []struct {
		name           string
		resource       unstructured.Unstructured
		resourceType   string
		expectedResult string
	}{
		{
			name:           "ResourceType not in transposedResourceListColumns['schedule']",
			resource:       unstructured.Unstructured{},
			resourceType:   "Deployment",
			expectedResult: "",
		},
		{
			name: "Spec does not exist",
			resource: unstructured.Unstructured{
				Object: map[string]interface{}{},
			},
			resourceType:   "CronJob",
			expectedResult: "",
		},
		{
			name: "Schedule field does not exist",
			resource: unstructured.Unstructured{
				Object: map[string]interface{}{
					"spec": map[string]interface{}{},
				},
			},
			resourceType:   "CronJob",
			expectedResult: "",
		},
		{
			name: "Schedule field exists",
			resource: unstructured.Unstructured{
				Object: map[string]interface{}{
					"spec": map[string]interface{}{
						"schedule": "*/5 * * * *",
					},
				},
			},
			resourceType:   "CronJob",
			expectedResult: "*/5 * * * *",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resourceDetailsTruncated := &models.ResourceListResourceList{}

			extractSchedule(tt.resource, tt.resourceType, resourceDetailsTruncated)

			if resourceDetailsTruncated.Schedule != tt.expectedResult {
				t.Errorf("Expected Schedule to be '%s', got '%s'", tt.expectedResult, resourceDetailsTruncated.Schedule)
			}
		})
	}
}

func TestExtractScope(t *testing.T) {
	tests := []struct {
		name           string
		resource       unstructured.Unstructured
		resourceType   string
		expectedResult string
	}{
		{
			name:           "ResourceType not in transposedResourceListColumns['scope']",
			resource:       unstructured.Unstructured{},
			resourceType:   "Deployment",
			expectedResult: "",
		},
		{
			name: "Spec does not exist",
			resource: unstructured.Unstructured{
				Object: map[string]interface{}{},
			},
			resourceType:   "CustomResourceDefinition",
			expectedResult: "",
		},
		{
			name: "Scope field does not exist",
			resource: unstructured.Unstructured{
				Object: map[string]interface{}{
					"spec": map[string]interface{}{},
				},
			},
			resourceType:   "CustomResourceDefinition",
			expectedResult: "",
		},
		{
			name: "Scope field exists",
			resource: unstructured.Unstructured{
				Object: map[string]interface{}{
					"spec": map[string]interface{}{
						"scope": "Namespaced",
					},
				},
			},
			resourceType:   "CustomResourceDefinition",
			expectedResult: "Namespaced",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resourceDetailsTruncated := &models.ResourceListResourceList{}

			extractScope(tt.resource, tt.resourceType, resourceDetailsTruncated)

			if resourceDetailsTruncated.Scope != tt.expectedResult {
				t.Errorf("Expected Scope to be '%s', got '%s'", tt.expectedResult, resourceDetailsTruncated.Scope)
			}
		})
	}
}

func TestExtractSelector(t *testing.T) {
	tests := []struct {
		name           string
		resource       unstructured.Unstructured
		resourceType   string
		expectedResult string
	}{
		{
			name:           "ResourceType not in transposedResourceListColumns['selector']",
			resource:       unstructured.Unstructured{},
			resourceType:   "ConfigMap",
			expectedResult: "",
		},
		{
			name: "Spec does not exist",
			resource: unstructured.Unstructured{
				Object: map[string]interface{}{},
			},
			resourceType:   "Service",
			expectedResult: "",
		},
		{
			name: "Selector field does not exist",
			resource: unstructured.Unstructured{
				Object: map[string]interface{}{
					"spec": map[string]interface{}{},
				},
			},
			resourceType:   "Service",
			expectedResult: "",
		},
		{
			name: "Selector field exists",
			resource: unstructured.Unstructured{
				Object: map[string]interface{}{
					"spec": map[string]interface{}{
						"selector": map[string]interface{}{
							"app":   "my-app",
							"tier":  "frontend",
							"track": "stable",
						},
					},
				},
			},
			resourceType:   "Service",
			expectedResult: "app:my-app, tier:frontend, track:stable",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resourceDetailsTruncated := &models.ResourceListResourceList{}

			extractSelector(tt.resource, tt.resourceType, resourceDetailsTruncated)

			if resourceDetailsTruncated.Selector != tt.expectedResult {
				t.Errorf("Expected Selector to be '%s', got '%s'", tt.expectedResult, resourceDetailsTruncated.Selector)
			}
		})
	}
}

func TestExtractSize(t *testing.T) {
	tests := []struct {
		name           string
		resource       unstructured.Unstructured
		resourceType   string
		expectedResult string
	}{
		{
			name:           "ResourceType not in transposedResourceListColumns['size']",
			resource:       unstructured.Unstructured{},
			resourceType:   "ConfigMap",
			expectedResult: "",
		},
		{
			name: "Spec does not exist",
			resource: unstructured.Unstructured{
				Object: map[string]interface{}{},
			},
			resourceType:   "PersistentVolumeClaim",
			expectedResult: "",
		},
		{
			name: "Size field does not exist",
			resource: unstructured.Unstructured{
				Object: map[string]interface{}{
					"spec": map[string]interface{}{},
				},
			},
			resourceType:   "PersistentVolumeClaim",
			expectedResult: "",
		},
		{
			name: "Size field exists",
			resource: unstructured.Unstructured{
				Object: map[string]interface{}{
					"spec": map[string]interface{}{
						"resources": map[string]interface{}{
							"requests": map[string]interface{}{
								"storage": "10Gi",
							},
						},
					},
				},
			},
			resourceType:   "PersistentVolumeClaim",
			expectedResult: "10Gi",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resourceDetailsTruncated := &models.ResourceListResourceList{}

			extractSize(tt.resource, tt.resourceType, resourceDetailsTruncated)

			if resourceDetailsTruncated.Size != tt.expectedResult {
				t.Errorf("Expected Size to be '%s', got '%s'", tt.expectedResult, resourceDetailsTruncated.Size)
			}
		})
	}
}

func TestExtractStatus(t *testing.T) {
	tests := []struct {
		name           string
		resource       unstructured.Unstructured
		resourceType   string
		expectedResult string
	}{
		{
			name:           "ResourceType not in transposedResourceListColumns['status']",
			resource:       unstructured.Unstructured{},
			resourceType:   "ConfigMap",
			expectedResult: "",
		},
		{
			name:           "ResourceType is Pod, status does not exist",
			resource:       unstructured.Unstructured{},
			resourceType:   "Pod",
			expectedResult: "",
		},
		{
			name: "ResourceType is Pod, phase does not exist",
			resource: unstructured.Unstructured{
				Object: map[string]interface{}{
					"status": map[string]interface{}{},
				},
			},
			resourceType:   "Pod",
			expectedResult: "",
		},
		{
			name: "ResourceType is Pod, phase exists",
			resource: unstructured.Unstructured{
				Object: map[string]interface{}{
					"status": map[string]interface{}{
						"phase": "Running",
					},
				},
			},
			resourceType:   "Pod",
			expectedResult: "Running",
		},
		{
			name: "ResourceType is PersistentVolumeClaim, phase exists",
			resource: unstructured.Unstructured{
				Object: map[string]interface{}{
					"status": map[string]interface{}{
						"phase": "Bound",
					},
				},
			},
			resourceType:   "PersistentVolumeClaim",
			expectedResult: "Bound",
		},
		{
			name: "ResourceType is Namespace, phase exists",
			resource: unstructured.Unstructured{
				Object: map[string]interface{}{
					"status": map[string]interface{}{
						"phase": "Active",
					},
				},
			},
			resourceType:   "Namespace",
			expectedResult: "Active",
		},
		{
			name: "ResourceType is PersistentVolume, phase exists",
			resource: unstructured.Unstructured{
				Object: map[string]interface{}{
					"status": map[string]interface{}{
						"phase": "Bound",
					},
				},
			},
			resourceType:   "PersistentVolume",
			expectedResult: "Bound",
		},
		{
			name:           "ResourceType is Service, status does not exist",
			resource:       unstructured.Unstructured{},
			resourceType:   "Service",
			expectedResult: "",
		},
		{
			name: "Service type is LoadBalancer, no ingresses",
			resource: unstructured.Unstructured{
				Object: map[string]interface{}{
					"spec": map[string]interface{}{
						"type": "LoadBalancer",
					},
					"status": map[string]interface{}{},
				},
			},
			resourceType:   "Service",
			expectedResult: "Pending",
		},
		{
			name: "Service type is LoadBalancer, with ingresses",
			resource: unstructured.Unstructured{
				Object: map[string]interface{}{
					"spec": map[string]interface{}{
						"type": "LoadBalancer",
					},
					"status": map[string]interface{}{
						"loadBalancer": map[string]interface{}{
							"ingress": []interface{}{
								map[string]interface{}{
									"ip": "192.168.1.1",
								},
							},
						},
					},
				},
			},
			resourceType:   "Service",
			expectedResult: "Active",
		},
		{
			name: "Service type is ClusterIP",
			resource: unstructured.Unstructured{
				Object: map[string]interface{}{
					"spec": map[string]interface{}{
						"type": "ClusterIP",
					},
					"status": map[string]interface{}{},
				},
			},
			resourceType:   "Service",
			expectedResult: "Active",
		},
		{
			name: "Service type unknown",
			resource: unstructured.Unstructured{
				Object: map[string]interface{}{
					"spec": map[string]interface{}{
						"type": "Unknown",
					},
					"status": map[string]interface{}{},
				},
			},
			resourceType:   "Service",
			expectedResult: "Active",
		},
		{
			name: "Service type LoadBalancer, error in unstructured.NestedSlice",
			resource: unstructured.Unstructured{
				Object: map[string]interface{}{
					"spec": map[string]interface{}{
						"type": "LoadBalancer",
					},
					"status": map[string]interface{}{
						"loadBalancer": "invalid",
					},
				},
			},
			resourceType:   "Service",
			expectedResult: "Unknown",
		},
		{
			name: "Service type LoadBalancer, ingresses is empty",
			resource: unstructured.Unstructured{
				Object: map[string]interface{}{
					"spec": map[string]interface{}{
						"type": "LoadBalancer",
					},
					"status": map[string]interface{}{
						"loadBalancer": map[string]interface{}{
							"ingress": []interface{}{},
						},
					},
				},
			},
			resourceType:   "Service",
			expectedResult: "Pending",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resourceDetailsTruncated := &models.ResourceListResourceList{}

			extractStatus(tt.resource, tt.resourceType, resourceDetailsTruncated)

			if resourceDetailsTruncated.Status != tt.expectedResult {
				t.Errorf("Expected Status to be '%s', got '%s'", tt.expectedResult, resourceDetailsTruncated.Status)
			}
		})
	}
}

func TestExtractStorageClass(t *testing.T) {
	tests := []struct {
		name           string
		resource       unstructured.Unstructured
		resourceType   string
		expectedResult string
	}{
		{
			name:           "ResourceType not in transposedResourceListColumns['storage_class']",
			resource:       unstructured.Unstructured{},
			resourceType:   "Deployment",
			expectedResult: "",
		},
		{
			name: "Spec does not exist",
			resource: unstructured.Unstructured{
				Object: map[string]interface{}{},
			},
			resourceType:   "PersistentVolumeClaim",
			expectedResult: "",
		},
		{
			name: "storageClassName field does not exist",
			resource: unstructured.Unstructured{
				Object: map[string]interface{}{
					"spec": map[string]interface{}{},
				},
			},
			resourceType:   "PersistentVolumeClaim",
			expectedResult: "",
		},
		{
			name: "storageClassName field exists",
			resource: unstructured.Unstructured{
				Object: map[string]interface{}{
					"spec": map[string]interface{}{
						"storageClassName": "fast-storage",
					},
				},
			},
			resourceType:   "PersistentVolumeClaim",
			expectedResult: "fast-storage",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resourceDetailsTruncated := &models.ResourceListResourceList{}

			extractStorageClass(tt.resource, tt.resourceType, resourceDetailsTruncated)

			if resourceDetailsTruncated.StorageClass != tt.expectedResult {
				t.Errorf("Expected StorageClass to be '%s', got '%s'", tt.expectedResult, resourceDetailsTruncated.StorageClass)
			}
		})
	}
}

func TestExtractSuspend(t *testing.T) {
	tests := []struct {
		name           string
		resource       unstructured.Unstructured
		resourceType   string
		expectedResult string
	}{
		{
			name:           "ResourceType not in transposedResourceListColumns['suspend']",
			resource:       unstructured.Unstructured{},
			resourceType:   "Deployment",
			expectedResult: "",
		},
		{
			name: "Spec does not exist",
			resource: unstructured.Unstructured{
				Object: map[string]interface{}{},
			},
			resourceType:   "CronJob",
			expectedResult: "",
		},
		{
			name: "suspend field does not exist",
			resource: unstructured.Unstructured{
				Object: map[string]interface{}{
					"spec": map[string]interface{}{},
				},
			},
			resourceType:   "CronJob",
			expectedResult: "",
		},
		{
			name: "suspend field exists (true)",
			resource: unstructured.Unstructured{
				Object: map[string]interface{}{
					"spec": map[string]interface{}{
						"suspend": true,
					},
				},
			},
			resourceType:   "CronJob",
			expectedResult: "true",
		},
		{
			name: "suspend field exists (false)",
			resource: unstructured.Unstructured{
				Object: map[string]interface{}{
					"spec": map[string]interface{}{
						"suspend": false,
					},
				},
			},
			resourceType:   "CronJob",
			expectedResult: "false",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resourceDetailsTruncated := &models.ResourceListResourceList{}

			extractSuspend(tt.resource, tt.resourceType, resourceDetailsTruncated)

			if resourceDetailsTruncated.Suspend != tt.expectedResult {
				t.Errorf("Expected Suspend to be '%s', got '%s'", tt.expectedResult, resourceDetailsTruncated.Suspend)
			}
		})
	}
}

func TestExtractTaints(t *testing.T) {
	tests := []struct {
		name           string
		resource       unstructured.Unstructured
		resourceType   string
		expectedResult string
	}{
		{
			name:           "ResourceType not in transposedResourceListColumns['taints']",
			resource:       unstructured.Unstructured{},
			resourceType:   "Pod",
			expectedResult: "",
		},
		{
			name: "Spec does not exist",
			resource: unstructured.Unstructured{
				Object: map[string]interface{}{},
			},
			resourceType:   "Node",
			expectedResult: "",
		},
		{
			name: "taints field does not exist",
			resource: unstructured.Unstructured{
				Object: map[string]interface{}{
					"spec": map[string]interface{}{},
				},
			},
			resourceType:   "Node",
			expectedResult: "0",
		},
		{
			name: "taints field exists with taints",
			resource: unstructured.Unstructured{
				Object: map[string]interface{}{
					"spec": map[string]interface{}{
						"taints": []interface{}{
							map[string]interface{}{
								"key":    "key1",
								"value":  "value1",
								"effect": "NoSchedule",
							},
							map[string]interface{}{
								"key":    "key2",
								"value":  "value2",
								"effect": "NoExecute",
							},
						},
					},
				},
			},
			resourceType:   "Node",
			expectedResult: "2",
		},
		{
			name: "taints field exists but empty",
			resource: unstructured.Unstructured{
				Object: map[string]interface{}{
					"spec": map[string]interface{}{
						"taints": []interface{}{},
					},
				},
			},
			resourceType:   "Node",
			expectedResult: "0",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resourceDetailsTruncated := &models.ResourceListResourceList{}

			extractTaints(tt.resource, tt.resourceType, resourceDetailsTruncated)

			if resourceDetailsTruncated.Taints != tt.expectedResult {
				t.Errorf("Expected Taints to be '%s', got '%s'", tt.expectedResult, resourceDetailsTruncated.Taints)
			}
		})
	}
}

func TestExtractType(t *testing.T) {
	tests := []struct {
		name           string
		resource       unstructured.Unstructured
		resourceType   string
		expectedResult string
	}{
		{
			name:           "ResourceType not in transposedResourceListColumns['type']",
			resource:       unstructured.Unstructured{},
			resourceType:   "ConfigMap",
			expectedResult: "",
		},
		{
			name: "Spec does not exist",
			resource: unstructured.Unstructured{
				Object: map[string]interface{}{},
			},
			resourceType:   "Service",
			expectedResult: "",
		},
		{
			name: "ResourceType is Secret, type field exists",
			resource: unstructured.Unstructured{
				Object: map[string]interface{}{
					"type": "Opaque",
				},
			},
			resourceType:   secretString,
			expectedResult: "Opaque",
		},
		{
			name: "ResourceType is Secret, type field does not exist",
			resource: unstructured.Unstructured{
				Object: map[string]interface{}{},
			},
			resourceType:   secretString,
			expectedResult: "",
		},
		{
			name: "ResourceType is not Secret, spec.type exists",
			resource: unstructured.Unstructured{
				Object: map[string]interface{}{
					"spec": map[string]interface{}{
						"type": "LoadBalancer",
					},
				},
			},
			resourceType:   "Service",
			expectedResult: "LoadBalancer",
		},
		{
			name: "ResourceType is not Secret, spec.type does not exist",
			resource: unstructured.Unstructured{
				Object: map[string]interface{}{
					"spec": map[string]interface{}{},
				},
			},
			resourceType:   "Service",
			expectedResult: "",
		},
		{
			name: "ResourceType is not Secret, spec does not exist",
			resource: unstructured.Unstructured{
				Object: map[string]interface{}{},
			},
			resourceType:   "Service",
			expectedResult: "",
		},
		{
			name: "Type field exists but wrong type",
			resource: unstructured.Unstructured{
				Object: map[string]interface{}{
					"type": 123,
				},
			},
			resourceType:   secretString,
			expectedResult: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resourceDetailsTruncated := &models.ResourceListResourceList{}

			extractType(tt.resource, tt.resourceType, resourceDetailsTruncated)

			if resourceDetailsTruncated.Type_ != tt.expectedResult {
				t.Errorf("Expected Type_ to be '%s', got '%s'", tt.expectedResult, resourceDetailsTruncated.Type_)
			}
		})
	}
}

func TestExtractVersion(t *testing.T) {
	tests := []struct {
		name           string
		resource       unstructured.Unstructured
		resourceType   string
		expectedResult string
	}{
		{
			name:           "ResourceType not in transposedResourceListColumns['version']",
			resource:       unstructured.Unstructured{},
			resourceType:   "Service",
			expectedResult: "",
		},
		{
			name: "ResourceType is Node, status.nodeInfo.kubeletVersion exists",
			resource: unstructured.Unstructured{
				Object: map[string]interface{}{
					"status": map[string]interface{}{
						"nodeInfo": map[string]interface{}{
							"kubeletVersion": "v1.20.4",
						},
					},
				},
			},
			resourceType:   nodeString,
			expectedResult: "v1.20.4",
		},
		{
			name: "ResourceType is Node, status.nodeInfo does not exist",
			resource: unstructured.Unstructured{
				Object: map[string]interface{}{
					"status": map[string]interface{}{},
				},
			},
			resourceType:   nodeString,
			expectedResult: "",
		},
		{
			name: "ResourceType is Node, status does not exist",
			resource: unstructured.Unstructured{
				Object: map[string]interface{}{},
			},
			resourceType:   nodeString,
			expectedResult: "",
		},
		{
			name: "ResourceType is CustomResourceDefinition, spec.versions exists with storage true",
			resource: unstructured.Unstructured{
				Object: map[string]interface{}{
					"spec": map[string]interface{}{
						"versions": []interface{}{
							map[string]interface{}{
								"name":    "v1",
								"storage": false,
							},
							map[string]interface{}{
								"name":    "v1beta1",
								"storage": true,
							},
						},
					},
				},
			},
			resourceType:   "CustomResourceDefinition",
			expectedResult: "v1beta1",
		},
		{
			name: "ResourceType is CustomResourceDefinition, spec.versions exists but no storage true",
			resource: unstructured.Unstructured{
				Object: map[string]interface{}{
					"spec": map[string]interface{}{
						"versions": []interface{}{
							map[string]interface{}{
								"name":    "v1",
								"storage": false,
							},
							map[string]interface{}{
								"name":    "v1beta1",
								"storage": false,
							},
						},
					},
				},
			},
			resourceType:   "CustomResourceDefinition",
			expectedResult: "",
		},
		{
			name: "ResourceType is CustomResourceDefinition, spec.versions does not exist",
			resource: unstructured.Unstructured{
				Object: map[string]interface{}{
					"spec": map[string]interface{}{},
				},
			},
			resourceType:   "CustomResourceDefinition",
			expectedResult: "",
		},
		{
			name: "ResourceType is CustomResourceDefinition, spec does not exist",
			resource: unstructured.Unstructured{
				Object: map[string]interface{}{},
			},
			resourceType:   "CustomResourceDefinition",
			expectedResult: "",
		},
		{
			name: "ResourceType is CustomResourceDefinition, version name missing",
			resource: unstructured.Unstructured{
				Object: map[string]interface{}{
					"spec": map[string]interface{}{
						"versions": []interface{}{
							map[string]interface{}{
								"storage": true,
							},
						},
					},
				},
			},
			resourceType:   "CustomResourceDefinition",
			expectedResult: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resourceDetailsTruncated := &models.ResourceListResourceList{}

			extractVersion(tt.resource, tt.resourceType, resourceDetailsTruncated)

			if resourceDetailsTruncated.Version != tt.expectedResult {
				t.Errorf("Expected Version to be '%s', got '%s'", tt.expectedResult, resourceDetailsTruncated.Version)
			}
		})
	}
}
