/*
 * KubernetesAccessManager - API
 *
 * This is a backend API server documentation for KubernetesAccessManager  Some useful links: - [Jira](https://samuelus.atlassian.net/jira/software/projects/ZPI/boards/4) - [Confluence](https://samuelus.atlassian.net/wiki/spaces/ZPI/overview)
 *
 * API version: 0.0.5
 * Generated by: Swagger Codegen (https://github.com/swagger-api/swagger-codegen.git)
 */
package models

import (
	"time"
)

// Helm release information
type HelmRelease struct {
	// Name of the Helm release
	Name string `json:"name,omitempty"`
	// Namespace where the release is installed
	Namespace string `json:"namespace,omitempty"`
	// Name and version of the Helm chart
	Chart string `json:"chart,omitempty"`
	// Current status of the release
	Status string `json:"status,omitempty"`
	// Last update timestamp
	Updated time.Time `json:"updated,omitempty"`
	// Revision number of the release
	Revision string `json:"revision,omitempty"`
	// App version of the release
	AppVersion string `json:"app_version,omitempty"`
}
