package cluster

import (
	"errors"
	"fmt"
	"github.com/ZPI-2024-25/KubernetesAccessManager/models"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"reflect"
	"testing"
)

func TestGetAllowedResourceTypes(t *testing.T) {
	expected := map[string]struct{}{
		"Pod":                      {},
		"Service":                  {},
		"Deployment":               {},
		"ConfigMap":                {},
		"StatefulSet":              {},
		"DaemonSet":                {},
		"Secret":                   {},
		"Ingress":                  {},
		"PersistentVolumeClaim":    {},
		"ReplicaSet":               {},
		"Node":                     {},
		"Namespace":                {},
		"CustomResourceDefinition": {},
		"PersistentVolume":         {},
		"Job":                      {},
		"CronJob":                  {},
		"ServiceAccount":           {},
		"StorageClass":             {},
		"ClusterRole":              {},
		"ClusterRoleBinding":       {},
	}

	result := getAllowedResourceTypes()

	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Expected %v, got %v", expected, result)
	}
}

func TestIsResourceTypeAllowed(t *testing.T) {
	allowedTypes := []string{
		"Pod",
		"Service",
		"Deployment",
		"ConfigMap",
		"StatefulSet",
		"DaemonSet",
		"Secret",
		"Ingress",
		"PersistentVolumeClaim",
		"ReplicaSet",
		"Node",
		"Namespace",
		"CustomResourceDefinition",
		"PersistentVolume",
		"Job",
		"CronJob",
		"ServiceAccount",
		"StorageClass",
		"ClusterRole",
		"ClusterRoleBinding",
	}

	for _, resourceType := range allowedTypes {
		if !isResourceTypeAllowed(resourceType) {
			t.Errorf("Expected resource type '%s' to be allowed", resourceType)
		}
	}

	disallowedTypes := []string{
		"InvalidType",
		"AnotherType",
		"",
	}

	for _, resourceType := range disallowedTypes {
		if isResourceTypeAllowed(resourceType) {
			t.Errorf("Expected resource type '%s' to be disallowed", resourceType)
		}
	}
}

func TestHandleKubernetesError(t *testing.T) {
	notFoundErr := apierrors.NewNotFound(schema.GroupResource{Group: "testGroup", Resource: "testResource"}, "testName")
	expectedNotFoundError := &models.ModelError{Code: 404, Message: fmt.Sprintf("Resource not found: %s", notFoundErr.Error())}

	result := handleKubernetesError(notFoundErr)
	if result.Code != expectedNotFoundError.Code || result.Message != expectedNotFoundError.Message {
		t.Errorf("Expected %v, got %v", expectedNotFoundError, result)
	}

	forbiddenErr := apierrors.NewForbidden(schema.GroupResource{Group: "testGroup", Resource: "testResource"}, "testName", errors.New("forbidden"))
	expectedForbiddenError := &models.ModelError{Code: 403, Message: fmt.Sprintf("Forbidden: %s", forbiddenErr.Error())}

	result = handleKubernetesError(forbiddenErr)
	if result.Code != expectedForbiddenError.Code || result.Message != expectedForbiddenError.Message {
		t.Errorf("Expected %v, got %v", expectedForbiddenError, result)
	}

	unauthorizedErr := apierrors.NewUnauthorized("unauthorized")
	expectedUnauthorizedError := &models.ModelError{Code: 401, Message: fmt.Sprintf("Unauthorized: %s", unauthorizedErr.Error())}

	result = handleKubernetesError(unauthorizedErr)
	if result.Code != expectedUnauthorizedError.Code || result.Message != expectedUnauthorizedError.Message {
		t.Errorf("Expected %v, got %v", expectedUnauthorizedError, result)
	}

	otherErr := errors.New("some other error")
	expectedOtherError := &models.ModelError{Code: 500, Message: fmt.Sprintf("Internal server error: %s", otherErr.Error())}

	result = handleKubernetesError(otherErr)
	if result.Code != expectedOtherError.Code || result.Message != expectedOtherError.Message {
		t.Errorf("Expected %v, got %v", expectedOtherError, result)
	}
}
