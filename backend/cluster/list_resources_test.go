package cluster

import (
	"github.com/ZPI-2024-25/KubernetesAccessManager/models"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"testing"
)

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
