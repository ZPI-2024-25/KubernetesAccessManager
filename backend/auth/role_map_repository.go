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
	// Subrole map is optional, if not provided, subroles will be ignored. Roles received in token are checked only with RoleMap
	RoleMap map[string]*models.Role
	SubroleMap map[string]*models.Role
}

var (
	instance      *RoleMapRepository
	once          sync.Once
)

func GetInstance() (*RoleMapRepository, error) {
	once.Do(func() {
		roleMap, subroleMap := GetRoleMapConfig("default", "role-map")
		if roleMap == nil {
			log.Printf("Failed to initialize RoleMapRepository")
			return
		}
		instance := &RoleMapRepository{RoleMap: roleMap, SubroleMap: subroleMap}

		log.Printf("RoleMapRepository initialized with %d roles %d subroles", len(instance.RoleMap), len(instance.SubroleMap))
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
		if rmr.subHasPermission(subrole, operation) {
			return true
		}
	}

	return false
}

func (rmr *RoleMapRepository) subHasPermission(rolename string, operation *models.Operation) bool {
	role := rmr.SubroleMap[rolename]
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
		if rmr.subHasPermission(subrole, operation) {
			return true
		}
	}

	return false
}

func hasCycle(roleMap map[string]*models.Role) bool {
	// Map to track visit state of each role
	visitState := make(map[string]int)

	for roleName := range roleMap {
		if visitState[roleName] == 0 { // unvisited
			if dfs(roleName, roleMap, visitState) {
				return true // Cycle detected
			}
		}
	}

	return false // No cycles found
}

func dfs(roleName string, roleMap map[string]*models.Role, visitState map[string]int) bool {
	const (
		unvisited = 0
		visiting  = 1
		visited   = 2
	)

	// If this role is currently being visited, a cycle is detected.
	if visitState[roleName] == visiting {
		return true
	}
	if visitState[roleName] == visited {
		return false
	}

	visitState[roleName] = visiting

	// visit all subroles
	for _, subrole := range roleMap[roleName].Subroles {
		if dfs(subrole, roleMap, visitState) {
			return true
		}
	}

	visitState[roleName] = visited
	return false
}



func GetRoleMapConfig (namespace string, name string) (map[string]*models.Role, map[string]*models.Role) {
	res, err := cluster.GetResource("ConfigMap", namespace, name, cluster.GetResourceInterface)
	if err != nil {
		return nil, nil
	}

	details := (*res.ResourceDetails).(*unstructured.Unstructured)
	
	roleMapData, foundRoleMap, err2 := unstructured.NestedString(details.Object, "data", "role-map")

	if err2 != nil || !foundRoleMap {
		log.Printf("Error retrieving roleMap data: %v", err)
		return nil, nil
	}

	roleMap := make(map[string]*models.Role)
	err2 = yaml.Unmarshal([]byte(roleMapData), &roleMap)
	if err2 != nil {
		log.Printf("Error parsing roleMap data: %v", err)
		return nil, nil
	}

	subroleMapData, foundSubRoleMap, err2 := unstructured.NestedString(details.Object, "data", "subrole-map")
	subroleMap := make(map[string]*models.Role)

	if foundSubRoleMap && err2 == nil {
		err2 = yaml.Unmarshal([]byte(subroleMapData), &subroleMap)
		if err2 != nil {
			log.Printf("Error parsing subroleMap data: %v", err)
			// No return as subroleMap is optional roleMap can be used without it
		}

		// clear subrole map, can't use it
		if hasCycle(subroleMap) {
			log.Printf("Cycle detected in subrole map")
			subroleMap = make(map[string]*models.Role)
		}
	}

	return roleMap, subroleMap
}


