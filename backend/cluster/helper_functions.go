package cluster

import (
	"context"
	"fmt"
	"log"
	"sync"

	"github.com/ZPI-2024-25/KubernetesAccessManager/models"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/discovery"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/kubernetes"
)

func getAllowedResourceTypes() map[string]struct{} {
	return map[string]struct{}{
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
}

func isResourceTypeAllowed(resourceType string) bool {
	allowedResourceTypes := getAllowedResourceTypes()
	_, exists := allowedResourceTypes[resourceType]
	return exists
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

func GetResourceInterface(resourceType string, namespace string, emptyNamespace string) (dynamic.ResourceInterface, *models.ModelError) {
	gvr, namespaced, httpErr := GetResourceGroupVersion(resourceType)
	if httpErr != nil {
		return nil, httpErr
	}

	dynamicClient, err := GetClientSet()
	if err != nil {
		return nil, &models.ModelError{Code: 500, Message: fmt.Sprintf("Failed to get client: %s", err)}
	}

	if namespaced {
		if namespace == "" {
			namespace = emptyNamespace
		}
		return dynamicClient.Resource(gvr).Namespace(namespace), nil
	} else {
		return dynamicClient.Resource(gvr), nil
	}
}

func handleKubernetesError(err error) *models.ModelError {
	if errors.IsNotFound(err) {
		return &models.ModelError{Code: 404, Message: fmt.Sprintf("Resource not found: %s", err)}
	} else if errors.IsForbidden(err) {
		return &models.ModelError{Code: 403, Message: fmt.Sprintf("Forbidden: %s", err)}
	} else if errors.IsUnauthorized(err) {
		return &models.ModelError{Code: 401, Message: fmt.Sprintf("Unauthorized: %s", err)}
	}
	return &models.ModelError{Code: 500, Message: fmt.Sprintf("Internal server error: %s", err)}
}

func WatchForChanges(namespace, resourceName string, mutex *sync.Mutex, updateFunc func(<-chan watch.Event, *sync.Mutex, string, string)) {
	for {
		config, err := GetConfig()
		if err != nil {
			log.Printf("Watch: Unable to get config")
			return
		}	
		clientset, err := kubernetes.NewForConfig(config)
		if err != nil {
			log.Printf("Watch: Unable to create clientset")
			return
		}	
		watcher, err := clientset.CoreV1().ConfigMaps(namespace).Watch(context.TODO(),
			metav1.SingleObject(metav1.ObjectMeta{Name: resourceName, Namespace: namespace}))
		if err != nil {
			log.Printf("Error creating watcher: %v", err)
			return
		}
		updateFunc(watcher.ResultChan(), mutex, namespace, resourceName)
	}
}
