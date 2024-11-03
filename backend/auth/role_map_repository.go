package auth

import (
	"errors"
	"log"
	"sync"

	"github.com/ZPI-2024-25/KubernetesAccessManager/cluster"
	"github.com/ZPI-2024-25/KubernetesAccessManager/models"
	"gopkg.in/yaml.v2"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
)

type RoleMapRepository struct {
	RoleMap map[string]*models.Role
}

var (
	instance      *RoleMapRepository
	once          sync.Once
)

func GetInstance() (*RoleMapRepository, error) {
	once.Do(func() {
		// TODO: Load role mappings from Kubernetes resources
		roleMap := GetRoleMapConfig("default", "role-map")
		if roleMap == nil {
			log.Printf("Failed to initialize RoleMapRepository")
			return
		}
		instance := &RoleMapRepository{RoleMap: *roleMap}

		log.Printf("RoleMapRepository initialized with %d roles", len(instance.RoleMap))
	})

	if instance == nil {
		return nil, errors.New("failed to initialize RoleMapRepository")
	}

	return instance, nil
}




func (rmr *RoleMapRepository) HasPermission(rolename string, operation *models.Operation) bool {
	role := rmr.RoleMap[rolename]
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

func GetRoleMapConfig (namespace string, name string) *map[string]*models.Role {
	res, err := cluster.GetResource("ConfigMap", namespace, name, cluster.GetResourceInterface)
	if err != nil {
		return nil
	}

	details := (*res.ResourceDetails).(*unstructured.Unstructured)
	
	roleMapData, found, err2 := unstructured.NestedString(details.Object, "data", "role-map")

	if err2 != nil || !found {
		log.Printf("Error retrieving roleMap data: %v", err)
		return nil
	}

	// Parse the YAML data in roleMap into the Role map structure
	roleMap := make(map[string]*models.Role)
	err2 = yaml.Unmarshal([]byte(roleMapData), &roleMap)
	if err2 != nil {
		log.Printf("Error parsing roleMap data: %v", err)
		return nil
	}

	return &roleMap
}


