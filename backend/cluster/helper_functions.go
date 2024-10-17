package cluster

import (
	"fmt"
	"github.com/ZPI-2024-25/KubernetesAccessManager/models"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/discovery"
)

func getAllowedResourceTypes() [20]string {
	return [20]string{
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

func GetResourceGroupVersion(resourceType string) (output schema.GroupVersionResource, namespaced bool, error *models.ModelError) {
	if !isResourceTypeAllowed(resourceType) {
		return schema.GroupVersionResource{}, false, &models.ModelError{Code: 400, Message: fmt.Sprintf("Invalid Resource Type")}
	}

	config, err := GetConfig()
	if err != nil {
		return schema.GroupVersionResource{}, false, &models.ModelError{Code: 500, Message: fmt.Sprintf("Failed to get config: %s", err)}
	}

	discoveryClient, _ := discovery.NewDiscoveryClientForConfig(config)

	apiResourceLists, _ := discoveryClient.ServerPreferredResources()

	for _, apiResourceList := range apiResourceLists {
		for _, apiResource := range apiResourceList.APIResources {
			if apiResource.Kind == resourceType {
				groupVersion, err := schema.ParseGroupVersion(apiResourceList.GroupVersion)
				if err != nil {
					return schema.GroupVersionResource{}, false, &models.ModelError{Code: 500, Message: fmt.Sprintf("%s", err.Error())}
				}

				return schema.GroupVersionResource{
					Group:    groupVersion.Group,
					Version:  groupVersion.Version,
					Resource: apiResource.Name,
				}, apiResource.Namespaced, nil
			}
		}
	}

	return schema.GroupVersionResource{}, false, &models.ModelError{Code: 400, Message: fmt.Sprintf("Invalid Resource Type")}
}
