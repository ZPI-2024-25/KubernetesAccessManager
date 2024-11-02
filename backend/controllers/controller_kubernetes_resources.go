package controllers

import (
	"fmt"
	"github.com/ZPI-2024-25/KubernetesAccessManager/cluster"
	"github.com/ZPI-2024-25/KubernetesAccessManager/models"
	"net/http"
)

func GetResourceController(w http.ResponseWriter, r *http.Request) {
	setJSONContentType(w)

	resourceType := getResourceType(r)
	resourceName := getResourceName(r)
	namespace := getNamespace(r)

	resource, err := cluster.GetResource(resourceType, namespace, resourceName, cluster.GetResourceInterface)
	if err != nil {
		writeJSONResponse(w, int(err.Code), err)
		return
	}

	writeJSONResponse(w, http.StatusOK, resource)
}

func ListResourcesController(w http.ResponseWriter, r *http.Request) {
	setJSONContentType(w)

	resourceType := getResourceType(r)
	namespace := getNamespace(r)

	resources, err := cluster.ListResources(resourceType, namespace, cluster.GetResourceInterface)
	if err != nil {
		writeJSONResponse(w, int(err.Code), err)
		return
	}

	writeJSONResponse(w, http.StatusOK, resources)
}

func CreateResourceController(w http.ResponseWriter, r *http.Request) {
	setJSONContentType(w)

	resourceType := getResourceType(r)
	namespace := getNamespace(r)

	var resource models.ResourceDetails
	if !decodeJSONBody(w, r, &resource.ResourceDetails) {
		return
	}

	resource, err := cluster.CreateResource(resourceType, namespace, resource, cluster.GetResourceInterface)
	if err != nil {
		fmt.Println(err)
		writeJSONResponse(w, int(err.Code), err)
		return
	}

	writeJSONResponse(w, http.StatusCreated, resource)
}

func DeleteResourceController(w http.ResponseWriter, r *http.Request) {
	setJSONContentType(w)

	resourceType := getResourceType(r)
	resourceName := getResourceName(r)
	namespace := getNamespace(r)

	err := cluster.DeleteResource(resourceType, namespace, resourceName, cluster.GetResourceInterface)
	if err != nil {
		writeJSONResponse(w, int(err.Code), err)
		return
	}

	status := models.Status{
		Status:  "Success",
		Code:    http.StatusOK,
		Message: fmt.Sprintf("Resource %s deleted successfully", resourceName),
	}
	writeJSONResponse(w, http.StatusOK, status)
}

func UpdateResourceController(w http.ResponseWriter, r *http.Request) {
	setJSONContentType(w)

	resourceType := getResourceType(r)
	resourceName := getResourceName(r)
	namespace := getNamespace(r)

	var resource models.ResourceDetails
	if !decodeJSONBody(w, r, &resource.ResourceDetails) {
		return
	}

	resource, err := cluster.UpdateResource(resourceType, namespace, resourceName, resource, cluster.GetResourceInterface)
	if err != nil {
		writeJSONResponse(w, int(err.Code), err)
		return
	}

	writeJSONResponse(w, http.StatusOK, resource)
}
