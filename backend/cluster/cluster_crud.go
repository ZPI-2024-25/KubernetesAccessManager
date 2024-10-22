package cluster

import (
	"context"
	"github.com/ZPI-2024-25/KubernetesAccessManager/models"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
)

const DefaultNamespace = "default"

var getResourceI = getResourceInterface

func GetResource(resourceType string, namespace string, resourceName string) (models.ResourceDetails, *models.ModelError) {
	resourceInterface, err := getResourceI(resourceType, namespace, DefaultNamespace)
	if err != nil {
		return models.ResourceDetails{}, err
	}

	var resource interface{}
	resource, getErr := resourceInterface.Get(context.TODO(), resourceName, metav1.GetOptions{})
	if getErr != nil {
		return models.ResourceDetails{}, handleKubernetesError(getErr)
	}

	return models.ResourceDetails{ResourceDetails: &resource}, nil
}

func CreateResource(resourceType string, namespace string, resource models.ResourceDetails) (models.ResourceDetails, *models.ModelError) {
	resourceInterface, err := getResourceI(resourceType, namespace, DefaultNamespace)
	if err != nil {
		return models.ResourceDetails{}, err
	}

	resourceMap, ok := (*resource.ResourceDetails).(map[string]interface{})
	if !ok {
		return models.ResourceDetails{}, &models.ModelError{Code: 400, Message: "Invalid resource format"}
	}

	resourceDefinition := &unstructured.Unstructured{
		Object: resourceMap,
	}

	if namespace != "" {
		resourceDefinition.SetNamespace(namespace)
	}

	var createdResource interface{}
	createdResource, createErr := resourceInterface.Create(context.TODO(), resourceDefinition, metav1.CreateOptions{})
	if createErr != nil {
		return models.ResourceDetails{}, handleKubernetesError(createErr)
	}

	return models.ResourceDetails{ResourceDetails: &createdResource}, nil
}

func DeleteResource(resourceType string, namespace string, resourceName string) *models.ModelError {
	resourceInterface, err := getResourceI(resourceType, namespace, DefaultNamespace)
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
	resourceInterface, err := getResourceI(resourceType, namespace, DefaultNamespace)
	if err != nil {
		return models.ResourceDetails{}, err
	}

	resourceMap, ok := (*resource.ResourceDetails).(map[string]interface{})
	if !ok {
		return models.ResourceDetails{}, &models.ModelError{Code: 400, Message: "Invalid resource format"}
	}

	unstructuredResource := &unstructured.Unstructured{Object: resourceMap}

	updatedResourceName, _, _ := unstructured.NestedString(resourceMap, "metadata", "name")
	if updatedResourceName != resourceName {
		return models.ResourceDetails{}, &models.ModelError{Code: 400, Message: "Invalid Input: Different resource names"}
	}

	var updatedResource interface{}
	updatedResource, updateErr := resourceInterface.Update(context.TODO(), unstructuredResource, metav1.UpdateOptions{})
	if updateErr != nil {
		return models.ResourceDetails{}, handleKubernetesError(updateErr)
	}

	return models.ResourceDetails{ResourceDetails: &updatedResource}, nil
}
