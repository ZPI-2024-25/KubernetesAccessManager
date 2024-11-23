package auth

import (
	"errors"
	"log"
	"sync"

	"github.com/ZPI-2024-25/KubernetesAccessManager/cluster"
	"github.com/ZPI-2024-25/KubernetesAccessManager/common"
	"github.com/ZPI-2024-25/KubernetesAccessManager/models"
	"gopkg.in/yaml.v2"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
)

type PermissionMatrix map[string]map[string]map[models.OperationType]struct{}

type RoleMapRepository struct {
	// Subrole map is optional, if not provided, subroles will be ignored. Roles received in token are checked only with RoleMap
	RoleMap      map[string]*models.Role
	SubroleMap   map[string]*models.Role
	flattenedMap map[string]PermissionMatrix
}

type operationConfig struct {
	Namespace  string                 `json:"namespace,omitempty"`
	Resource   string                 `json:"resource,omitempty"`
	Operations []models.OperationType `json:"operations,omitempty"`
}

type roleConfig struct {
	Name     string            `json:"name,omitempty"`
	Permit   []operationConfig `json:"permit,omitempty"`
	Deny     []operationConfig `json:"deny,omitempty"`
	Subroles []string          `json:"subroles,omitempty"`
}

var (
	instance *RoleMapRepository
	once     sync.Once
)

func GetRoleMapInstance() (*RoleMapRepository, error) {
	once.Do(func() {
		roleMapNamespace := common.GetOrDefaultEnv("ROLEMAP_NAMESPACE", common.DEFAULT_ROLEMAP_NAMESPACE)
		roleMapName := common.GetOrDefaultEnv("ROLEMAP_NAME", common.DEFAULT_ROLEMAP_NAME)

		roleMap, subroleMap := GetRoleMapConfig(roleMapNamespace, roleMapName)
		if roleMap == nil {
			return
		}
		permissionMatrix := createPermissionMatrix(roleMap, subroleMap)
		instance = &RoleMapRepository{
			RoleMap:      roleMap,
			SubroleMap:   subroleMap,
			flattenedMap: permissionMatrix,
		}
		log.Printf("RoleMapRepository initialized with %d roles %d subroles", len(instance.RoleMap), len(instance.SubroleMap))
	})
	if instance == nil {
		return nil, errors.New("RoleMapRepository is not initialized")
	}
	return instance, nil
}

func (rmr *RoleMapRepository) HasPermission(rolenames []string, operation *models.Operation) bool {
	for _, role := range rolenames {
		if flatHasPermission(operation, rmr.flattenedMap[role]) {
			return true
		}
	}
	return false
}

func (rmr *RoleMapRepository) GetAllPermissions(roles []string) PermissionMatrix {
	pmatrix := make(PermissionMatrix)
	first := true
	for _, r := range roles {
		if _, exists := rmr.flattenedMap[r]; exists {
			if first {
				pmatrix = deepCopy(rmr.flattenedMap[r])
				pruneResourcesNamespaces(pmatrix)
				first = false
			} else {
				pmatrix = addMatrix(pmatrix, rmr.flattenedMap[r])
				pruneResourcesNamespaces(pmatrix)
			}
		}
	}
	return pmatrix
}

func (rmr *RoleMapRepository) HasPermissionInAnyNamespace(rolenames []string, resource string, op models.OperationType) bool {
	for _, role := range rolenames {
		for _, namespace := range rmr.flattenedMap[role] {
			if _, exists := namespace[resource]; exists {
				if _, exists := namespace[resource][op]; exists {
					return true
				}
			} else if _, exists := namespace["*"][op]; exists {
				return true
			}
		}
	}
	return false
}

func flatHasPermission(op *models.Operation, matrix PermissionMatrix) bool {
	var namespace string
	if _, exists := matrix[op.Namespace]; exists {
		namespace = op.Namespace
	} else {
		namespace = "*"
	}
	var resource string
	if _, exists := matrix[namespace][op.Resource]; exists {
		resource = op.Resource
	} else {
		resource = "*"
	}

	_, exists := matrix[namespace][resource][op.Type]
	return exists
}

func createPermissionMatrix(
	roleMap map[string]*models.Role,
	subroleMap map[string]*models.Role,
) map[string]PermissionMatrix {
	superMatrix := make(map[string]PermissionMatrix)
	for _, role := range roleMap {
		matrix := toMatrix(role, subroleMap)
		superMatrix[role.Name] = matrix
	}
	return superMatrix
}

func GetRoleMapConfig(namespace string, name string) (map[string]*models.Role, map[string]*models.Role) {
	res, err := cluster.GetResource("ConfigMap", namespace, name, cluster.GetResourceInterface)
	if err != nil {
		return nil, nil
	}

	details := (*res.ResourceDetails).(*unstructured.Unstructured)

	roleMapConfigData, foundRoleMapConfig, err2 := unstructured.NestedString(details.Object, "data", "role-map")
	if err2 != nil || !foundRoleMapConfig {
		log.Printf("Error retrieving roleMap data: %v", err)
		return nil, nil
	}

	roleMapConfig := make(map[string]*roleConfig)
	err2 = yaml.Unmarshal([]byte(roleMapConfigData), &roleMapConfig)
	if err2 != nil {
		log.Printf("Error parsing roleMap data: %v", err)
		return nil, nil
	}

	subroleMapConfigData, foundSubroleMapConfig, err2 := unstructured.NestedString(details.Object, "data", "subrole-map")
	subroleMapConfig := make(map[string]*roleConfig)

	var subroleMap map[string]*models.Role

	if !foundSubroleMapConfig || err2 != nil {
		log.Printf("Error retrieving subroleMap data")
	} else {
		err2 = yaml.Unmarshal([]byte(subroleMapConfigData), &subroleMapConfig)
		if err2 != nil {
			log.Printf("Error parsing subroleMap data: %v", err)
			// No return as subroleMap is optional roleMap can be used without it
		}
		subroleMap = fromRoleMapConfig(subroleMapConfig)
		if hasCycle(subroleMap) {
			log.Printf("Cycle detected in subrole map")
			subroleMap = make(map[string]*models.Role) // clear subrole map, can't use it
		}
	}

	return fromRoleMapConfig(roleMapConfig), subroleMap
}

func fromRoleMapConfig(config map[string]*roleConfig) map[string]*models.Role {
	roleMap := make(map[string]*models.Role)
	for name, roleConfig := range config {
		role := fromRoleConfig(roleConfig)
		if role.Name == "" {
			role.Name = name
		}
		roleMap[name] = role
	}
	return roleMap
}

func fromRoleConfig(config *roleConfig) *models.Role {
	permit := fromOperationConfigList(config.Permit)
	deny := fromOperationConfigList(config.Deny)

	role := &models.Role{
		Name:     config.Name,
		Subroles: config.Subroles,
	}
	if len(permit) > 0 {
		role.Permit = permit
	}
	if len(deny) > 0 {
		role.Deny = deny
	}
	return role
}

func fromOperationConfigList(operations []operationConfig) []models.Operation {
	ops := make([]models.Operation, 0)
	for _, opConfig := range operations {
		namespace := opConfig.Namespace
		if namespace == "" {
			namespace = "*"
		}
		resource := opConfig.Resource
		if resource == "" {
			resource = "*"
		}
		if len(opConfig.Operations) == 0 {
			ops = append(ops, models.Operation{
				Namespace: namespace,
				Resource:  resource,
				Type:      models.All,
			})
		} else {
			for _, opType := range opConfig.Operations {
				ops = append(ops, models.Operation{
					Namespace: namespace,
					Resource:  resource,
					Type:      opType,
				})
			}
		}
	}
	return ops
}
