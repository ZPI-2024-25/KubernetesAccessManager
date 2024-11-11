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
	"fmt"
	"github.com/ZPI-2024-25/KubernetesAccessManager/common"
	"github.com/gorilla/mux"
	"net/http"
	"strings"
)

type Route struct {
	Name        string
	Method      string
	Pattern     string
	HandlerFunc http.HandlerFunc
}

type Routes []Route

func NewRouter() *mux.Router {
	router := mux.NewRouter().StrictSlash(true)
	for _, route := range routes {
		var handler http.Handler
		handler = route.HandlerFunc
		handler = common.Logger(handler, route.Name)

		router.
			Methods(route.Method).
			Path(route.Pattern).
			Name(route.Name).
			Handler(handler)
	}

	return router
}

func Index(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello World!")
}

var routes = Routes{
	Route{
		"Index",
		"GET",
		"/api/v1/",
		Index,
	},

	Route{
		"GetHelmRelease",
		strings.ToUpper("Get"),
		"/api/v1/helm/releases/{releaseName}",
		GetHelmRelease,
	},

	Route{
		"GetHelmReleaseHistory",
		strings.ToUpper("Get"),
		"/api/v1/helm/releases/{releaseName}/history",
		GetHelmReleaseHistory,
	},

	Route{
		"ListHelmReleases",
		strings.ToUpper("Get"),
		"/api/v1/helm/releases",
		ListHelmReleases,
	},

	Route{
		"RollbackHelmRelease",
		strings.ToUpper("Post"),
		"/api/v1/helm/releases/{releaseName}/rollback",
		RollbackHelmRelease,
	},

	Route{
		"UninstallHelmRelease",
		strings.ToUpper("Delete"),
		"/api/v1/helm/releases/{releaseName}",
		UninstallHelmRelease,
	},

	Route{
		"CreateResource",
		strings.ToUpper("Post"),
		"/api/v1/k8s/{resourceType}",
		CreateResource,
	},

	Route{
		"DeleteResource",
		strings.ToUpper("Delete"),
		"/api/v1/k8s/{resourceType}/{resourceName}",
		DeleteResource,
	},

	Route{
		"GetResource",
		strings.ToUpper("Get"),
		"/api/v1/k8s/{resourceType}/{resourceName}",
		GetResource,
	},

	Route{
		"ListResources",
		strings.ToUpper("Get"),
		"/api/v1/k8s/{resourceType}",
		ListResources,
	},

	Route{
		"UpdateResource",
		strings.ToUpper("Put"),
		"/api/v1/k8s/{resourceType}/{resourceName}",
		UpdateResource,
	},

	Route{
		"CheckLoginStatus",
		strings.ToUpper("Get"),
		"/api/v1/auth/status",
		CheckLoginStatus,
	},
}
