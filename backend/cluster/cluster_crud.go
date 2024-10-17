package cluster

import (
	"context"
	"github.com/ZPI-2024-25/KubernetesAccessManager/models"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
)

func GetResource(resourceType string, namespace string, resourceName string) (models.ResourceDetails, *models.ModelError) {
	resourceInterface, err := getResourceInterface(resourceType, namespace)
	if err != nil {
		return models.ResourceDetails{}, err
	}

	resource, getErr := resourceInterface.Get(context.TODO(), resourceName, metav1.GetOptions{})
	if getErr != nil {
		return models.ResourceDetails{}, handleKubernetesError(getErr)
	}

	var outputInterface interface{} = resource.Object
	return models.ResourceDetails{ResourceDetails: &outputInterface}, nil
}

func CreateResource(resourceType string, namespace string, resource models.ResourceDetails) (models.ResourceDetails, *models.ModelError) {
	resourceInterface, err := getResourceInterface(resourceType, namespace)
	if err != nil {
		return models.ResourceDetails{}, err
	}

	resourceDetails := resource.ResourceDetails
	resourceMap, ok := (*resourceDetails).(map[string]interface{})
	if !ok {
		return models.ResourceDetails{}, &models.ModelError{Code: 400, Message: "Invalid resource format"}
	}

	resourceDefinition := &unstructured.Unstructured{
		Object: resourceMap,
	}

	if namespace != "" {
		resourceDefinition.SetNamespace(namespace)
	}

	createdResource, createErr := resourceInterface.Create(context.TODO(), resourceDefinition, metav1.CreateOptions{})
	if createErr != nil {
		return models.ResourceDetails{}, handleKubernetesError(createErr)
	}

	var details interface{} = createdResource.Object
	return models.ResourceDetails{ResourceDetails: &details}, nil
}

func DeleteResource(resourceType string, namespace string, resourceName string) *models.ModelError {
	resourceInterface, err := getResourceInterface(resourceType, namespace)
	if err != nil {
		return err
	}

	deleteErr := resourceInterface.Delete(context.TODO(), resourceName, metav1.DeleteOptions{})
	if deleteErr != nil {
		return handleKubernetesError(deleteErr)
	}

	return nil
}

func UpdateResource(resourceType string, namespace string, resourceName string, resource models.ResourceDetails) (models.ResourceDetails, *models.ModelError) {
	resourceInterface, err := getResourceInterface(resourceType, namespace)
	if err != nil {
		return models.ResourceDetails{}, err
	}

	currentResource, getErr := resourceInterface.Get(context.TODO(), resourceName, metav1.GetOptions{})
	if getErr != nil {
		return models.ResourceDetails{}, handleKubernetesError(getErr)
	}

	resourceMap, ok := (*resource.ResourceDetails).(map[string]interface{})
	if !ok {
		return models.ResourceDetails{}, &models.ModelError{Code: 400, Message: "Invalid resource format"}
	}

	updatedResourceName, _, _ := unstructured.NestedString(resourceMap, "metadata", "name")
	if updatedResourceName != resourceName {
		return models.ResourceDetails{}, &models.ModelError{Code: 400, Message: "Invalid Input: Different resource names"}
	}

	currentResource.Object = resourceMap

	updatedResource, updateErr := resourceInterface.Update(context.TODO(), currentResource, metav1.UpdateOptions{})
	if updateErr != nil {
		return models.ResourceDetails{}, handleKubernetesError(updateErr)
	}

	var details interface{} = updatedResource.Object
	return models.ResourceDetails{ResourceDetails: &details}, nil
}
