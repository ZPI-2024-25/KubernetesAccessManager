package auth

import (
	"errors"
	"sync"
	"github.com/ZPI-2024-25/KubernetesAccessManager/models"
)

type RoleMapRepository struct {
	roleMap map[string]*models.Role
}

var (
	instance      *RoleMapRepository
	once          sync.Once
)

func GetInstance() (*RoleMapRepository, error) {
	once.Do(func() {
		// TODO: Load role mappings from Kubernetes resources
		instance = &RoleMapRepository{
			roleMap: map[string]*models.Role{
				"admin": {
					Name: "admin",
					Permit: []models.Operation{
						{
							Resource: "all",
							Type: models.All,
							Namespace: "all",
						},
					},
					Subroles: []string{"user"},
				},
				"user": {
					Name: "user",
					Permit: []models.Operation{
						{
							Resource: "Pod",
							Type: models.Read,
							Namespace: "all",
						},
						{
							Resource: "Service",
							Type: models.Read,
							Namespace: "all",
						},
						{
							Resource: "Deployment",
							Type: models.Read,
							Namespace: "all",
						},
					},
					Deny: []models.Operation{
						{
							Resource: "all",
							Type: models.All,
							Namespace: "forbidden",
						},
					},
					Subroles: []string{"guest"},
				},
			},
		}
	})

	if instance == nil {
		return nil, errors.New("failed to initialize RoleMapRepository")
	}

	return instance, nil
}




func (rmr *RoleMapRepository) HasPermission(rolename string, operation *models.Operation) bool {
	role := rmr.roleMap[rolename]
	if role == nil {
		return false
	}

	for _, deny := range role.Deny {
		if deny.IsSuper(operation) {
			return false
		}
	}

	for _, permit := range role.Permit {
		if permit.IsSuper(operation) {
			return true
		}
	}

	// Recursively check subroles, if any matches, return true
	for _, subrole := range role.Subroles {
		if rmr.HasPermission(subrole, operation) {
			return true
		}
	}

	return false
}


