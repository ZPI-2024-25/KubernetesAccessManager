package cluster

import (
	"context"
	"fmt"
	"github.com/ZPI-2024-25/KubernetesUserManager/go/common"
	"github.com/ZPI-2024-25/KubernetesUserManager/go/models"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/discovery"
	"strings"
)

func GetResourceGroupVersion(resourceType string) (schema.GroupVersionResource, error) {
	dynamicClientSingleton, _ := common.GetInstance()
	config := dynamicClientSingleton.GetConfig()

	discoveryClient, _ := discovery.NewDiscoveryClientForConfig(config)

	apiResourceLists, _ := discoveryClient.ServerPreferredResources()

	for _, apiResourceList := range apiResourceLists {
		for _, apiResource := range apiResourceList.APIResources {
			if strings.ToLower(apiResource.Kind) == strings.ToLower(resourceType) {
				groupVersion, err := schema.ParseGroupVersion(apiResourceList.GroupVersion)
				if err != nil {
					return schema.GroupVersionResource{}, fmt.Errorf("failed to parse GroupVersion: %v", err)
				}

				return schema.GroupVersionResource{
					Group:    groupVersion.Group,
					Version:  groupVersion.Version,
					Resource: apiResource.Name,
				}, nil
			}
		}
	}

	return schema.GroupVersionResource{}, fmt.Errorf("resource type '%s' not found in the cluster", resourceType)
}

func GetClusterResource(resourceType string, resourceName string) (models.Resource, error) {
	gvr, _ := GetResourceGroupVersion(resourceType)

	singleton, err := common.GetInstance()
	dynamicClient := singleton.GetClientSet()

	resource, err := dynamicClient.Resource(gvr).Get(context.TODO(), resourceName, metav1.GetOptions{})
	if err != nil {
		if errors.IsNotFound(err) {
			return models.Resource{}, fmt.Errorf("pod not found")
		} else {
			return models.Resource{}, fmt.Errorf("error: %s", err.Error())
		}
	}

	metadata, _, _ := unstructured.NestedMap(resource.Object, "metadata")

	spec, _, _ := unstructured.NestedMap(resource.Object, "spec")

	status, _, _ := unstructured.NestedMap(resource.Object, "status")

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

func GetNamespacedResource(resourceType string, namespace string, resourceName string) (models.Resource, error) {
	gvr, _ := GetResourceGroupVersion(resourceType)

	singleton, err := common.GetInstance()
	dynamicClient := singleton.GetClientSet()

	resource, err := dynamicClient.Resource(gvr).Namespace(namespace).Get(context.TODO(), resourceName, metav1.GetOptions{})
	if err != nil {
		if errors.IsNotFound(err) {
			return models.Resource{}, fmt.Errorf("pod not found")
		} else {
			return models.Resource{}, fmt.Errorf("error: %s", err.Error())
		}
	}

	metadata, _, _ := unstructured.NestedMap(resource.Object, "metadata")

	spec, _, _ := unstructured.NestedMap(resource.Object, "spec")

	status, _, _ := unstructured.NestedMap(resource.Object, "status")

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

func CreateClusterResource(resourceType string, resource models.Resource) {
}

func CreateNamespacedResource(resourceType string, namespace string, resource models.Resource) {
}

func DeleteClusterResource(resourceType string, resourceName string) error {
	gvr, _ := GetResourceGroupVersion(resourceType)

	singleton, err := common.GetInstance()
	dynamicClient := singleton.GetClientSet()

	err = dynamicClient.Resource(gvr).Delete(context.TODO(), resourceName, metav1.DeleteOptions{})
	if err != nil {
		if errors.IsNotFound(err) {
			return fmt.Errorf("pod not found")
		} else {
			return fmt.Errorf("error: %s", err.Error())
		}
	}

	return nil
}

func DeleteNamespacedResource(resourceType string, namespace string, resourceName string) error {
	gvr, _ := GetResourceGroupVersion(resourceType)

	singleton, err := common.GetInstance()
	dynamicClient := singleton.GetClientSet()

	err = dynamicClient.Resource(gvr).Namespace(namespace).Delete(context.TODO(), resourceName, metav1.DeleteOptions{})
	if err != nil {
		if errors.IsNotFound(err) {
			return fmt.Errorf("pod not found")
		} else {
			return fmt.Errorf("error: %s", err.Error())
		}
	}

	return nil
}
