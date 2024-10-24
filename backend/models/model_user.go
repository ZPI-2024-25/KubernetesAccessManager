/*
 * KubernetesAccessManager - API
 *
 * This is a backend API server documentation for KubernetesAccessManager  Some useful links: - [Jira](https://samuelus.atlassian.net/jira/software/projects/ZPI/boards/4) - [Confluence](https://samuelus.atlassian.net/wiki/spaces/ZPI/overview)
 *
 * API version: 0.0.3
 * Generated by: Swagger Codegen (https://github.com/swagger-api/swagger-codegen.git)
 */
package models

import (
	"time"
)

// User information
type User struct {
	// Unique identifier of the user
	Id string `json:"id,omitempty"`
	// Username of the user
	Username string `json:"username,omitempty"`
	// Email address of the user
	Email string `json:"email,omitempty"`
	// List of user permissions
	Permissions []string `json:"permissions,omitempty"`
	// Account creation timestamp
	CreatedAt time.Time `json:"createdAt,omitempty"`
}
