package cluster

import (
	"context"
	"fmt"
	"github.com/ZPI-2024-25/KubernetesAccessManager/common"
	"github.com/ZPI-2024-25/KubernetesAccessManager/models"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
)

//func GetResource(resourceType string, namespace string, resourceName string) (models.Resource, *models.ModelError) {
//	gvr, httpErr := GetResourceGroupVersion(resourceType)
//	if httpErr != nil {
//		return models.Resource{}, httpErr
//	}
//
//	singleton, err := common.GetInstance()
//	if err != nil {
//		return models.Resource{}, &models.ModelError{Code: 500, Message: fmt.Sprintf("Failed to get client instance: %s", err)}
//	}
//	dynamicClient := singleton.GetClientSet()
//
//	var resource *unstructured.Unstructured
//	if namespace == "" {
//		resource, err = dynamicClient.Resource(gvr).Get(context.TODO(), resourceName, metav1.GetOptions{})
//	} else {
//		resource, err = dynamicClient.Resource(gvr).Namespace(namespace).Get(context.TODO(), resourceName, metav1.GetOptions{})
//	}
//
//	if err != nil {
//		if errors.IsNotFound(err) {
//			return models.Resource{}, &models.ModelError{Code: 404, Message: fmt.Sprintf("Resource not found: %s", err)}
//		} else {
//			return models.Resource{}, &models.ModelError{Code: 500, Message: fmt.Sprintf("Error: %s", err)}
//		}
//	}
//
//	metadata, _, err := unstructured.NestedMap(resource.Object, "metadata")
//	if err != nil {
//		return models.Resource{}, &models.ModelError{Code: 500, Message: fmt.Sprintf("Error: %s", err)}
//	}
//
//	spec, _, err := unstructured.NestedMap(resource.Object, "spec")
//	if err != nil {
//		return models.Resource{}, &models.ModelError{Code: 500, Message: fmt.Sprintf("Error: %s", err)}
//	}
//
//	status, _, err := unstructured.NestedMap(resource.Object, "status")
//	if err != nil {
//		return models.Resource{}, &models.ModelError{Code: 500, Message: fmt.Sprintf("Error: %s", err)}
//	}
//
//	var metadataSwagger interface{} = metadata
//	var specSwagger interface{} = spec
//	var statusSwagger interface{} = status
//
//	return models.Resource{
//		ApiVersion: resource.GetAPIVersion(),
//		Kind:       resource.GetKind(),
//		Metadata:   &metadataSwagger,
//		Spec:       &specSwagger,
//		Status:     &statusSwagger,
//	}, nil
//}

func CreateResource(resourceType string, namespace string, resource models.ResourceDetails) (models.ResourceDetails, *models.ModelError) {
	gvr, httpErr := GetResourceGroupVersion(resourceType)
	if httpErr != nil {
		return models.ResourceDetails{}, httpErr
	}

	singleton, err := common.GetInstance()
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
	if namespace == "" {
		createdResource, err = dynamicClient.Resource(gvr).Create(context.TODO(), resourceDefinition, metav1.CreateOptions{})
	} else {
		createdResource, err = dynamicClient.Resource(gvr).Namespace(namespace).Create(context.TODO(), resourceDefinition, metav1.CreateOptions{})
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

//func UpdateResource(resourceType string, namespace string, name string, resource models.Resource) (models.Resource, *models.ModelError) {
//	gvr, httpErr := GetResourceGroupVersion(resourceType)
//	if httpErr != nil {
//		return models.Resource{}, httpErr
//	}
//
//	singleton, err := common.GetInstance()
//	if err != nil {
//		return models.Resource{}, &models.ModelError{Code: 500, Message: fmt.Sprintf("Failed to get client instance: %s", err)}
//	}
//	dynamicClient := singleton.GetClientSet()
//
//	var currentResource *unstructured.Unstructured
//	if namespace == "" {
//		currentResource, err = dynamicClient.Resource(gvr).Get(context.TODO(), name, metav1.GetOptions{})
//	} else {
//		currentResource, err = dynamicClient.Resource(gvr).Namespace(namespace).Get(context.TODO(), name, metav1.GetOptions{})
//	}
//
//	if err != nil {
//		fmt.Println("Here 1")
//		return models.Resource{}, &models.ModelError{Code: 404, Message: fmt.Sprintf("Resource not found: %s", err)}
//	}
//
//	if resource.Metadata != nil {
//		metadataMap, ok := (*resource.Metadata).(map[string]interface{})
//		if ok {
//			currentResource.Object["metadata"] = metadataMap
//		} else {
//			return models.Resource{}, &models.ModelError{Code: 400, Message: "Invalid metadata format"}
//		}
//	} else {
//		return models.Resource{}, &models.ModelError{Code: 400, Message: "Missing metadata"}
//	}
//
//	if resource.Spec != nil {
//		specMap, ok := (*resource.Spec).(map[string]interface{})
//		if ok {
//			currentResource.Object["spec"] = specMap
//		} else {
//			return models.Resource{}, &models.ModelError{Code: 400, Message: "Invalid spec format"}
//		}
//	} else {
//		return models.Resource{}, &models.ModelError{Code: 400, Message: "Missing spec"}
//	}
//
//	var updatedResource *unstructured.Unstructured
//
//	if namespace == "" {
//		updatedResource, err = dynamicClient.Resource(gvr).Update(context.TODO(), currentResource, metav1.UpdateOptions{})
//	} else {
//		updatedResource, err = dynamicClient.Resource(gvr).Namespace(namespace).Update(context.TODO(), currentResource, metav1.UpdateOptions{})
//	}
//	if err != nil {
//		fmt.Println("Here 2")
//		return models.Resource{}, &models.ModelError{Code: 500, Message: fmt.Sprintf("Error: %s", err)}
//	}
//
//	metadata, _, err := unstructured.NestedMap(updatedResource.Object, "metadata")
//	if err != nil {
//		fmt.Println("Here 3")
//		return models.Resource{}, &models.ModelError{Code: 500, Message: fmt.Sprintf("Error: %s", err)}
//	}
//
//	spec, _, err := unstructured.NestedMap(updatedResource.Object, "spec")
//	if err != nil {
//		return models.Resource{}, &models.ModelError{Code: 500, Message: fmt.Sprintf("Error: %s", err)}
//	}
//
//	status, _, err := unstructured.NestedMap(updatedResource.Object, "status")
//	if err != nil {
//		return models.Resource{}, &models.ModelError{Code: 500, Message: fmt.Sprintf("Error: %s", err)}
//	}
//
//	var metadataSwagger interface{} = metadata
//	var specSwagger interface{} = spec
//	var statusSwagger interface{} = status
//
//	return models.Resource{
//		ApiVersion: updatedResource.GetAPIVersion(),
//		Kind:       updatedResource.GetKind(),
//		Metadata:   &metadataSwagger,
//		Spec:       &specSwagger,
//		Status:     &statusSwagger,
//	}, nil
//}
