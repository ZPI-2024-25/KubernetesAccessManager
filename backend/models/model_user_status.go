/*
 * KubernetesAccessManager - API
 *
 * This is a backend API server documentation for KubernetesAccessManager  Some useful links: - [Jira](https://samuelus.atlassian.net/jira/software/projects/ZPI/boards/4) - [Confluence](https://samuelus.atlassian.net/wiki/spaces/ZPI/overview)
 *
 * API version: 0.0.5
 * Generated by: Swagger Codegen (https://github.com/swagger-api/swagger-codegen.git)
 */
package models

// User session details and privileges
type UserStatus struct {
	Permissions *UserStatusPermissions `json:"permissions,omitempty"`

	User *UserStatusUser `json:"user,omitempty"`
}
