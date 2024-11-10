package auth

import (
	"errors"
	"log"
	"os"
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
	flattenedMap map[string]map[string]map[string][]models.OperationType
	flattenedCompiled bool
}

var (
	instance      *RoleMapRepository
	once          sync.Once
)

func GetRoleMapInstance() (*RoleMapRepository, error) {
	once.Do(func() {
		roleMapNamespace := getOrDefaultEnv("ROLEMAP_NAMESPACE", "default")
		roleMapName := getOrDefaultEnv("ROLEMAP_NAME", "role-mapper")

		roleMap, subroleMap := GetRoleMapConfig(roleMapNamespace, roleMapName)
		if roleMap == nil {
			log.Printf("Failed to initialize RoleMapRepository")
			return
		}
		instance = &RoleMapRepository{
			RoleMap: roleMap, 
			SubroleMap: subroleMap,
		}

		log.Printf("RoleMapRepository initialized with %d roles %d subroles", len(instance.RoleMap), len(instance.SubroleMap))
	})

	if instance == nil {
		return nil, errors.New("failed to initialize RoleMapRepository")
	}

	return instance, nil
}

func getOrDefaultEnv(key, defaultValue string) string {
	val := os.Getenv(key)
	if val == "" {
		log.Printf("Environment variable %s not set, using default value %s", key, defaultValue)
		return defaultValue
	}
	return val
}

func (rmr *RoleMapRepository) HasPermission(rolenames []string, operation *models.Operation) bool {
	visited := make(map[string]struct{})
	for _, role := range rolenames {
		if rmr.flattenedCompiled {
			if flatHasPermission(operation, rmr.flattenedMap[role]) {
				return true
			}
		} else {
			role := rmr.RoleMap[role]
			if hasPermission(role, rmr.SubroleMap, operation, visited) {
				return true
			}
		}
	}
	return false
}

func flatHasPermission(op *models.Operation, matrix map[string]map[string][]models.OperationType) bool {
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

	for _, opType := range matrix[namespace][resource] {
		if opType == op.Type {
			return true
		}
	}
	return false
}

func createPermissionMatrix(
	roleMap map[string] *models.Role,
	subroleMap map[string]*models.Role,
	) map[string]map[string]map[string]map[models.OperationType]struct{} {
	superMatrix := make(map[string]map[string]map[string]map[models.OperationType]struct{})
	for _, role := range roleMap {
		matrix := toMatrix(role, subroleMap)
		superMatrix[role.Name] = matrix
	}
	return superMatrix
}

func toMatrix(role *models.Role, subroleMap map[string] *models.Role) map[string]map[string]map[models.OperationType]struct{} {
	var matrix map[string]map[string]map[models.OperationType]struct{}
	first := true
	for _, child := range role.Subroles {
		if childRole, exists := subroleMap[child]; exists {
			if first {
				matrix = toMatrix(childRole, subroleMap)
				first = false
			} else {
				matrix = addMatrix(matrix, toMatrix(childRole, subroleMap))
			}
		}
	}
	if first { //no subroles init the matrix
		matrix = make(map[string]map[string]map[models.OperationType]struct{})
		matrix["*"] = make(map[string]map[models.OperationType]struct{})
		matrix["*"]["*"] = make(map[models.OperationType]struct{})
	}
	for _, permit := range role.Permit {
		addPermitToMatrix(matrix, permit)
	}
	for _, deny := range role.Deny {
		restrictMatrix(matrix, deny)
	}
	return matrix
}

func addMatrix(m1 map[string]map[string]map[models.OperationType]struct{}, 
	m2 map[string]map[string]map[models.OperationType]struct{}) (
		map[string]map[string]map[models.OperationType]struct{}) {
	x1, y1 := len(m1), len(m1["*"])
	x2, y2 := len(m2), len(m2["*"])
	if x1 * y1 < x2 * y2 {
		m1, m2 = m2, m1
	}
	for namespace := range m2 {
		if _, exists := m1[namespace]; !exists {
			expandNamespaces(namespace, m1)
		}
	}
	for resource := range m2["*"] {
		if _, exists := m1["*"][resource]; !exists {
			expandResources(resource, m1)
		}
	}
	for namespace, resources := range m2 {
		for resource, operations := range resources {
			for opType := range operations {
				m1[namespace][resource][opType] = struct{}{}
			}
		}
	}
	return m1
}

func addPermitToMatrix(matrix map[string]map[string]map[models.OperationType]struct{}, permit models.Operation) {
	_, hasNamespace := matrix[permit.Namespace]
	_, hasResource := matrix["*"][permit.Resource]

	if !hasNamespace {
		expandNamespaces(permit.Namespace, matrix)
	}
	if !hasResource {
		expandResources(permit.Resource, matrix)
	}

	if permit.Namespace == "*" {
		if permit.Resource == "*" { // * * op
			for _, resources := range matrix {
				for _, operations := range resources {
					addOp(operations, permit.Type)
				}
			}
		} else { // * r op
			for _, resources := range matrix {
				addOp(resources[permit.Resource], permit.Type)
			}
		}
	} else {
		if permit.Resource == "*" { // n * op
			for _, operations := range matrix[permit.Namespace] {
				addOp(operations, permit.Type)
			}
		} else { // n r op
			addOp(matrix[permit.Namespace][permit.Resource], permit.Type)
		}
	}
}
func restrictMatrix(matrix map[string]map[string]map[models.OperationType]struct{}, deny models.Operation) {
	_, hasNamespace := matrix[deny.Namespace]
	_, hasResource := matrix["*"][deny.Resource]

	if !hasNamespace {
		expandNamespaces(deny.Namespace, matrix)
	}
	if !hasResource {
		expandResources(deny.Resource, matrix)
	}

	if deny.Namespace == "*" {
		if deny.Resource == "*" { // * * op
			for _, resources := range matrix {
				for _, operations := range resources {
					removeOp(operations, deny.Type)
				}
			}
		} else { // * r op
			for _, resources := range matrix {
				removeOp(resources[deny.Resource], deny.Type)
			}
		}
	} else {
		if deny.Resource == "*" { // n * op
			for _, operations := range matrix[deny.Namespace] {
				removeOp(operations, deny.Type)
			}
		} else { // n r op
			removeOp(matrix[deny.Namespace][deny.Resource], deny.Type)
		}
	}
}

func addOp(ops map[models.OperationType]struct{}, opType models.OperationType) {
		if opType == models.All {
			ops[models.Create] = struct{}{}
			ops[models.Read] = struct{}{}
			ops[models.Update] = struct{}{}
			ops[models.Delete] = struct{}{}
			ops[models.List] = struct{}{}
		} else {
			ops[opType] = struct{}{}
		}
	}

func removeOp(ops map[models.OperationType]struct{}, opType models.OperationType) {
	if opType == models.All {
		delete(ops, models.Create)
		delete(ops, models.Read)
		delete(ops, models.Update)
		delete(ops, models.Delete)
		delete(ops, models.List)
	} else {
		delete(ops, opType)
	}
}

func expandNamespaces(namespace string, matrix map[string]map[string]map[models.OperationType]struct{}) {
	matrix[namespace] = make(map[string]map[models.OperationType]struct{})
	for resource, ops := range matrix["*"] {
		matrix[namespace][resource] = make(map[models.OperationType]struct{})
		for opType := range ops {
			matrix[namespace][resource][opType] = struct{}{}
		}
	}
}

func expandResources(resource string, matrix map[string]map[string]map[models.OperationType]struct{}) {
	for _, namespaced := range matrix {
		namespaced[resource] = make(map[models.OperationType]struct{})
		for opType := range namespaced["*"] {
			namespaced[resource][opType] = struct{}{}
		}
	}
}




// 1. Macierz pozwoleń -> zaczynamy * * kopiujemy uprawnienia z ** do poszczególnych nam, resource type
// 2. Deny -> odejmujemy od odpowiednich kolumn/wiersz uprawnienia

// Dodawanie macierzy
//

func flattenRoleMap (roleMap map[string]*models.Role, subroleMap map[string]*models.Role) map[string]map[string]map[string][]models.OperationType{
	flattenedMap := make(map[string]map[string]map[string][]models.OperationType)
	for rolename, role := range roleMap {
		oneRoleMap := flattenRole(role, subroleMap)
		flattenedMap[rolename] = oneRoleMap
	}
	return flattenedMap
}

func flattenRole(role *models.Role, subroleMap map[string]*models.Role) map[string]map[string][]models.OperationType {
	namespaces := findUsedOperationsInGraph(
		role, 
		subroleMap, 
		func (op *models.Operation) string { return op.Namespace },
		make(map[string]struct{}), 
		make(map[string]struct{}),
	)

	resources := findUsedOperationsInGraph(
		role, 
		subroleMap, 
		func (op *models.Operation) string { return op.Resource },
		make(map[string]struct{}), 
		make(map[string]struct{}),
	)

	opTypes := models.GetAllOperationTypes()

	result := make(map[string]map[string][]models.OperationType)

	for namespace := range namespaces {
		result[namespace] = make(map[string][]models.OperationType)
		for resource := range resources {
			result[namespace][resource] = make([]models.OperationType, 0, len(opTypes))
			for _, opType := range opTypes {
				op := models.Operation{
					Namespace: namespace,
					Type: opType,
					Resource: resource,
				}

				if hasPermission(role, subroleMap, &op, make(map[string]struct{})) {
					result[namespace][resource] = append(result[namespace][resource], opType)
				}
			}
		}
	}

	return result
}

type getFieldFromOp func(*models.Operation) string

func findUsedOperationsInGraph(
	role *models.Role,
	subroleMap map[string]*models.Role, 
	fieldFunction getFieldFromOp, 
	values map[string]struct{}, 
	visited map[string]struct{},
) map[string]struct{}{
	for _, op := range role.Deny {
		value := fieldFunction(&op)
		values[value] = struct{}{}
	}

	for _, op := range role.Permit {
		value := fieldFunction(&op)
		values[value] = struct{}{}
	}

	for _, subroleName := range role.Subroles {
		subrole := subroleMap[subroleName]
		findUsedOperationsInGraph(subrole, subroleMap, fieldFunction, values, visited)
	}

	return values
}
	

func hasPermission(role *models.Role, subroleMap map[string]*models.Role, operation *models.Operation, visited map[string]struct{}) bool {
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
	for _, subroleName := range role.Subroles {
		if _, exists := visited[subroleName]; !exists {
			subrole := subroleMap[subroleName]
			visited[subroleName] = struct{}{}
			if hasPermission(subrole, subroleMap, operation, visited) {
				return true
			}
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
	role, exists := roleMap[roleName]
    if !exists || role == nil {
        visitState[roleName] = visited
        return false
    }

    for _, subrole := range role.Subroles {
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

	if !foundSubRoleMap || err2 != nil {
		log.Printf("Error retrieving subroleMap data")
	} else {
		err2 = yaml.Unmarshal([]byte(subroleMapData), &subroleMap)
		if err2 != nil {
			log.Printf("Error parsing subroleMap data: %v", err)
			// No return as subroleMap is optional roleMap can be used without it
		}

		if hasCycle(subroleMap) {
			log.Printf("Cycle detected in subrole map")
			subroleMap = make(map[string]*models.Role) // clear subrole map, can't use it
		}
	}

	return roleMap, subroleMap
}


