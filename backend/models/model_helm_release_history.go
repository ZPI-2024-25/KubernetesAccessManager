/*
 * KubernetesAccessManager - API
 *
 * This is a backend API server documentation for KubernetesAccessManager  Some useful links: - [Jira](https://samuelus.atlassian.net/jira/software/projects/ZPI/boards/4) - [Confluence](https://samuelus.atlassian.net/wiki/spaces/ZPI/overview)
 *
 * API version: 0.0.2
 * Generated by: Swagger Codegen (https://github.com/swagger-api/swagger-codegen.git)
 */
package models

import (
	"time"
)

// Helm release history entry
type HelmReleaseHistory struct {
	// Revision number
	Revision int32 `json:"revision,omitempty"`
	// Update timestamp
	Updated time.Time `json:"updated,omitempty"`
	// Status of the release at this revision
	Status string `json:"status,omitempty"`
	// Chart version used
	Chart string `json:"chart,omitempty"`
	// App version of the release
	AppVersion string `json:"app_version,omitempty"`
	// Description of the revision
	Description string `json:"description,omitempty"`
}