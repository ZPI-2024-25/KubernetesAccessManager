package cluster

import (
	"context"
	"fmt"
	"github.com/ZPI-2024-25/KubernetesUserManager/common"
	"github.com/ZPI-2024-25/KubernetesUserManager/models"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/discovery"
)

func getAllowedResourceTypes() [16]string {
	return [16]string{
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
		"CronJob",
		"ServiceAccount",
	}
}

func isResourceTypeAllowed(resourceType string) bool {
	allowedResourceTypes := getAllowedResourceTypes()
	for _, allowedResourceType := range allowedResourceTypes {
		if allowedResourceType == resourceType {
			return true
		}
	}
	return false
}

func GetResourceGroupVersion(resourceType string) (schema.GroupVersionResource, *models.ModelError) {
	if !isResourceTypeAllowed(resourceType) {
		return schema.GroupVersionResource{}, &models.ModelError{Code: 400, Message: fmt.Sprintf("Resource type '%s' not allowed", resourceType)}
	}

	dynamicClientSingleton, _ := common.GetInstance()
	config := dynamicClientSingleton.GetConfig()

	discoveryClient, _ := discovery.NewDiscoveryClientForConfig(config)

	apiResourceLists, _ := discoveryClient.ServerPreferredResources()

	for _, apiResourceList := range apiResourceLists {
		for _, apiResource := range apiResourceList.APIResources {
			if apiResource.Kind == resourceType {
				groupVersion, err := schema.ParseGroupVersion(apiResourceList.GroupVersion)
				if err != nil {
					return schema.GroupVersionResource{}, &models.ModelError{Code: 500, Message: fmt.Sprintf("%s", err.Error())}
				}

				return schema.GroupVersionResource{
					Group:    groupVersion.Group,
					Version:  groupVersion.Version,
					Resource: apiResource.Name,
				}, nil
			}
		}
	}

	return schema.GroupVersionResource{}, &models.ModelError{Code: 400, Message: fmt.Sprintf("Resource type '%s' not found", resourceType)}
}

func GetResource(resourceType string, namespace string, resourceName string) (models.Resource, *models.ModelError) {
	gvr, httpErr := GetResourceGroupVersion(resourceType)
	if httpErr != nil {
		return models.Resource{}, httpErr
	}

	singleton, err := common.GetInstance()
	if err != nil {
		return models.Resource{}, &models.ModelError{Code: 500, Message: fmt.Sprintf("Failed to get client instance: %s", err)}
	}
	dynamicClient := singleton.GetClientSet()

	var resource *unstructured.Unstructured
	if namespace == "" {
		resource, err = dynamicClient.Resource(gvr).Get(context.TODO(), resourceName, metav1.GetOptions{})
	} else {
		resource, err = dynamicClient.Resource(gvr).Namespace(namespace).Get(context.TODO(), resourceName, metav1.GetOptions{})
	}

	if err != nil {
		if errors.IsNotFound(err) {
			return models.Resource{}, &models.ModelError{Code: 404, Message: fmt.Sprintf("Resource not found: %s", err)}
		} else {
			return models.Resource{}, &models.ModelError{Code: 500, Message: fmt.Sprintf("Error: %s", err)}
		}
	}

	metadata, _, err := unstructured.NestedMap(resource.Object, "metadata")
	if err != nil {
		return models.Resource{}, &models.ModelError{Code: 500, Message: fmt.Sprintf("Error: %s", err)}
	}

	spec, _, err := unstructured.NestedMap(resource.Object, "spec")
	if err != nil {
		return models.Resource{}, &models.ModelError{Code: 500, Message: fmt.Sprintf("Error: %s", err)}
	}

	status, _, err := unstructured.NestedMap(resource.Object, "status")
	if err != nil {
		return models.Resource{}, &models.ModelError{Code: 500, Message: fmt.Sprintf("Error: %s", err)}
	}

	var metadataSwagger interface{} = metadata
	var specSwagger interface{} = spec
	var statusSwagger interface{} = status

	return models.Resource{
		ApiVersion: resource.GetAPIVersion(),
		Kind:       resource.GetKind(),
		Metadata:   &metadataSwagger,
		Spec:       &specSwagger,
		Status:     &statusSwagger,
	}, nil
}

func ListClusterResources(resourceType string) {
}

func ListNamespacedResources(resourceType string, namespace string) {
}

func CreateResource(resourceType string, namespace string, resource models.Resource) (models.Resource, *models.ModelError) {
	gvr, httpErr := GetResourceGroupVersion(resourceType)
	if httpErr != nil {
		return models.Resource{}, httpErr
	}

	singleton, err := common.GetInstance()
	if err != nil {
		return models.Resource{}, &models.ModelError{Code: 500, Message: fmt.Sprintf("Failed to get client instance: %s", err)}
	}
	dynamicClient := singleton.GetClientSet()

	resourceDefinition := &unstructured.Unstructured{
		Object: map[string]interface{}{
			"apiVersion": resource.ApiVersion,
			"kind":       resource.Kind,
			"metadata":   resource.Metadata,
			"spec":       resource.Spec,
		},
	}

	var createdResource *unstructured.Unstructured
	if namespace == "" {
		createdResource, err = dynamicClient.Resource(gvr).Create(context.TODO(), resourceDefinition, metav1.CreateOptions{})
	} else {
		createdResource, err = dynamicClient.Resource(gvr).Namespace(namespace).Create(context.TODO(), resourceDefinition, metav1.CreateOptions{})
	}

	if err != nil {
		return models.Resource{}, &models.ModelError{Code: 500, Message: fmt.Sprintf("Error: %s", err)}
	}

	metadata, _, err := unstructured.NestedMap(createdResource.Object, "metadata")
	if err != nil {
		return models.Resource{}, &models.ModelError{Code: 500, Message: fmt.Sprintf("Error: %s", err)}
	}

	spec, _, err := unstructured.NestedMap(createdResource.Object, "spec")
	if err != nil {
		return models.Resource{}, &models.ModelError{Code: 500, Message: fmt.Sprintf("Error: %s", err)}
	}

	status, _, err := unstructured.NestedMap(createdResource.Object, "status")
	if err != nil {
		return models.Resource{}, &models.ModelError{Code: 500, Message: fmt.Sprintf("Error: %s", err)}
	}

	var metadataSwagger interface{} = metadata
	var specSwagger interface{} = spec
	var statusSwagger interface{} = status

	return models.Resource{
		ApiVersion: createdResource.GetAPIVersion(),
		Kind:       createdResource.GetKind(),
		Metadata:   &metadataSwagger,
		Spec:       &specSwagger,
		Status:     &statusSwagger,
	}, nil
}

func DeleteResource(resourceType string, namespace string, resourceName string) *models.ModelError {
	gvr, httpErr := GetResourceGroupVersion(resourceType)
	if httpErr != nil {
		return httpErr
	}

	singleton, err := common.GetInstance()
	if err != nil {
		return &models.ModelError{Code: 500, Message: fmt.Sprintf("Failed to get client instance: %w", err)}
	}
	dynamicClient := singleton.GetClientSet()

	if namespace != "" {
		err = dynamicClient.Resource(gvr).Delete(context.TODO(), resourceName, metav1.DeleteOptions{})
	} else {
		err = dynamicClient.Resource(gvr).Namespace(namespace).Delete(context.TODO(), resourceName, metav1.DeleteOptions{})
	}
	if err != nil {
		if errors.IsNotFound(err) {
			return &models.ModelError{Code: 404, Message: fmt.Sprintf("Resource not found: %w", err)}
		} else if errors.IsForbidden(err) {
			return &models.ModelError{Code: 403, Message: fmt.Sprintf("Forbidden: %w", err)}
		} else if errors.IsUnauthorized(err) {
			return &models.ModelError{Code: 401, Message: fmt.Sprintf("Unauthorized: %w", err)}
		}
		return &models.ModelError{Code: 500, Message: fmt.Sprintf("Internal server error: %w", err)}
	}

	return nil
}
