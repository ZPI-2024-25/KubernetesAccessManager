package cluster

import (
	"context"
	"fmt"
	"github.com/ZPI-2024-25/KubernetesAccessManager/models"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
)

func GetResource(resourceType string, namespace string, resourceName string) (models.ResourceDetails, *models.ModelError) {
	gvr, namespaced, httpErr := GetResourceGroupVersion(resourceType)
	if httpErr != nil {
		return models.ResourceDetails{}, httpErr
	}

	singleton, err := GetInstance()
	if err != nil {
		return models.ResourceDetails{}, &models.ModelError{Code: 500, Message: fmt.Sprintf("Failed to get client instance: %s", err)}
	}
	dynamicClient := singleton.GetClientSet()

	var resource *unstructured.Unstructured
	if namespaced {
		if namespace == "" {
			resource, err = dynamicClient.Resource(gvr).Namespace("default").Get(context.TODO(), resourceName, metav1.GetOptions{})
		} else {
			resource, err = dynamicClient.Resource(gvr).Namespace(namespace).Get(context.TODO(), resourceName, metav1.GetOptions{})
		}
	} else {
		resource, err = dynamicClient.Resource(gvr).Get(context.TODO(), resourceName, metav1.GetOptions{})
	}

	if err != nil {
		if errors.IsNotFound(err) {
			return models.ResourceDetails{}, &models.ModelError{Code: 404, Message: fmt.Sprintf("Resource not found: %s", err)}
		} else {
			return models.ResourceDetails{}, &models.ModelError{Code: 500, Message: fmt.Sprintf("Error: %s", err)}
		}
	}

	output := make(map[string]interface{})

	fieldsToKeep := []string{"apiVersion", "kind", "metadata", "spec", "status"}

	for _, field := range fieldsToKeep {
		if value, found, _ := unstructured.NestedFieldCopy(resource.Object, field); found {
			output[field] = value
		}
	}

	var outputInterface interface{} = output

	return models.ResourceDetails{ResourceDetails: &outputInterface}, nil
}

func CreateResource(resourceType string, namespace string, resource models.ResourceDetails) (models.ResourceDetails, *models.ModelError) {
	gvr, namespaced, httpErr := GetResourceGroupVersion(resourceType)
	if httpErr != nil {
		return models.ResourceDetails{}, httpErr
	}

	singleton, err := GetInstance()
	if err != nil {
		return models.ResourceDetails{}, &models.ModelError{Code: 500, Message: fmt.Sprintf("Failed to get client instance: %s", err)}
	}
	dynamicClient := singleton.GetClientSet()

	resourceDetails := resource.ResourceDetails

	resourceMap, _ := (*resourceDetails).(map[string]interface{})

	resourceDefinition := &unstructured.Unstructured{
		Object: resourceMap,
	}

	var createdResource *unstructured.Unstructured

	if namespaced {
		if namespace == "" {
			createdResource, err = dynamicClient.Resource(gvr).Namespace("default").Create(context.TODO(), resourceDefinition, metav1.CreateOptions{})
		} else {
			createdResource, err = dynamicClient.Resource(gvr).Namespace(namespace).Create(context.TODO(), resourceDefinition, metav1.CreateOptions{})
		}
	} else {
		createdResource, err = dynamicClient.Resource(gvr).Create(context.TODO(), resourceDefinition, metav1.CreateOptions{})
	}

	if err != nil {
		return models.ResourceDetails{}, &models.ModelError{Code: 500, Message: fmt.Sprintf("Error: %s", err)}
	}

	var details interface{} = createdResource.Object

	return models.ResourceDetails{
		ResourceDetails: &details,
	}, nil
}

func DeleteResource(resourceType string, namespace string, resourceName string) *models.ModelError {
	gvr, namespaced, httpErr := GetResourceGroupVersion(resourceType)
	if httpErr != nil {
		return httpErr
	}

	singleton, err := GetInstance()
	if err != nil {
		return &models.ModelError{Code: 500, Message: fmt.Sprintf("Failed to get client instance: %s", err)}
	}
	dynamicClient := singleton.GetClientSet()

	if namespaced {
		if namespace == "" {
			err = dynamicClient.Resource(gvr).Namespace("default").Delete(context.TODO(), resourceName, metav1.DeleteOptions{})
		} else {
			err = dynamicClient.Resource(gvr).Namespace(namespace).Delete(context.TODO(), resourceName, metav1.DeleteOptions{})
		}
	} else {
		err = dynamicClient.Resource(gvr).Delete(context.TODO(), resourceName, metav1.DeleteOptions{})
	}

	if err != nil {
		if errors.IsNotFound(err) {
			return &models.ModelError{Code: 404, Message: fmt.Sprintf("Resource not found: %s", err)}
		} else if errors.IsForbidden(err) {
			return &models.ModelError{Code: 403, Message: fmt.Sprintf("Forbidden: %s", err)}
		} else if errors.IsUnauthorized(err) {
			return &models.ModelError{Code: 401, Message: fmt.Sprintf("Unauthorized: %s", err)}
		}
		return &models.ModelError{Code: 500, Message: fmt.Sprintf("Internal server error: %s", err)}
	}

	return nil
}

func UpdateResource(resourceType string, namespace string, resourceName string, resource models.ResourceDetails) (models.ResourceDetails, *models.ModelError) {
	gvr, namespaced, httpErr := GetResourceGroupVersion(resourceType)
	if httpErr != nil {
		return models.ResourceDetails{}, httpErr
	}

	singleton, err := GetInstance()
	if err != nil {
		return models.ResourceDetails{}, &models.ModelError{Code: 500, Message: fmt.Sprintf("Failed to get client instance: %s", err)}
	}
	dynamicClient := singleton.GetClientSet()

	var currentResource *unstructured.Unstructured
	if namespaced {
		if namespace == "" {
			currentResource, err = dynamicClient.Resource(gvr).Namespace("default").Get(context.TODO(), resourceName, metav1.GetOptions{})
		} else {
			currentResource, err = dynamicClient.Resource(gvr).Namespace(namespace).Get(context.TODO(), resourceName, metav1.GetOptions{})
		}
	} else {
		currentResource, err = dynamicClient.Resource(gvr).Get(context.TODO(), resourceName, metav1.GetOptions{})
	}

	if err != nil {
		return models.ResourceDetails{}, &models.ModelError{Code: 404, Message: fmt.Sprintf("Resource not found: %s", err)}
	}

	currentResource.Object = (*resource.ResourceDetails).(map[string]interface{})

	updatedResourceName, _, _ := unstructured.NestedString(currentResource.Object, "metadata", "resourceName")
	if updatedResourceName != resourceName {
		return models.ResourceDetails{}, &models.ModelError{Code: 400, Message: "Invalid Input: Different resource names"}
	}

	var updatedResource *unstructured.Unstructured

	if namespaced {
		if namespace == "" {
			updatedResource, err = dynamicClient.Resource(gvr).Namespace("default").Update(context.TODO(), currentResource, metav1.UpdateOptions{})
		} else {
			updatedResource, err = dynamicClient.Resource(gvr).Namespace(namespace).Update(context.TODO(), currentResource, metav1.UpdateOptions{})
		}
	} else {
		updatedResource, err = dynamicClient.Resource(gvr).Update(context.TODO(), currentResource, metav1.UpdateOptions{})
	}

	if err != nil {
		return models.ResourceDetails{}, &models.ModelError{Code: 500, Message: fmt.Sprintf("Error: %s", err)}
	}

	var details interface{} = updatedResource.Object

	return models.ResourceDetails{
		ResourceDetails: &details,
	}, nil
}
