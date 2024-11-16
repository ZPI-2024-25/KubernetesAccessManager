/*
 * KubernetesAccessManager - API
 *
 * This is a backend API server documentation for KubernetesAccessManager  Some useful links: - [Jira](https://samuelus.atlassian.net/jira/software/projects/ZPI/boards/4) - [Confluence](https://samuelus.atlassian.net/wiki/spaces/ZPI/overview)
 *
 * API version: 0.0.5
 * Generated by: Swagger Codegen (https://github.com/swagger-api/swagger-codegen.git)
 */
package api

import (
	"github.com/ZPI-2024-25/KubernetesAccessManager/controllers"
	"net/http"
)

func CreateResource(w http.ResponseWriter, r *http.Request) {
	controllers.CreateResourceController(w, r)
}

func DeleteResource(w http.ResponseWriter, r *http.Request) {
	controllers.DeleteResourceController(w, r)
}

func GetResource(w http.ResponseWriter, r *http.Request) {
	controllers.GetResourceController(w, r)
}

func ListResources(w http.ResponseWriter, r *http.Request) {
	controllers.ListResourcesController(w, r)
}

func UpdateResource(w http.ResponseWriter, r *http.Request) {
	controllers.UpdateResourceController(w, r)
}
