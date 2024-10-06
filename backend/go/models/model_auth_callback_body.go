/*
 * KubernetesUserManager - API
 *
 * This is a backend API server documentation for KubernetesUserManager  Some useful links: - [Jira](https://samuelus.atlassian.net/jira/software/projects/ZPI/boards/4) - [Confluence](https://samuelus.atlassian.net/wiki/spaces/ZPI/overview)
 *
 * API version: 0.0.1
 * Generated by: Swagger Codegen (https://github.com/swagger-api/swagger-codegen.git)
 */
package models

type AuthCallbackBody struct {
	// Authorization code.
	Code string `json:"code,omitempty"`
}