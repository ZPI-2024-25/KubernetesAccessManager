package auth

import (
	"github.com/ZPI-2024-25/KubernetesAccessManager/models"
)

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

func deepCopy(m PermissionMatrix) PermissionMatrix {
	copy := make(PermissionMatrix)
	for namespace, resources := range m {
		copy[namespace] = make(map[string]map[models.OperationType]struct{})
		for resource, operations := range resources {
			copy[namespace][resource] = make(map[models.OperationType]struct{})
			for opType := range operations {
				copy[namespace][resource][opType] = struct{}{}
			}
		}
	}
	return copy
}

func pruneResourcesNamespaces(matrix PermissionMatrix) bool {
	// delete unnecessary resources, namespaces if all operations are the same as in *
	wasPruned := false
	for namespace, resources := range matrix { //prune namespaces
		if namespace != "*" {
			allOpsSame := true
			for resource, operations := range resources {
				wildcardOps := matrix["*"][resource]
				if !sameOps(operations, wildcardOps) {
					allOpsSame = false
				}
			}
			if allOpsSame {
				delete(matrix, namespace)
				wasPruned = true
			}
		}
	}
	for resource := range matrix["*"] { //prune resources
		if resource != "*" {
			allOpsSame := true
			for namespace, resources := range matrix {
				wildcardOps := matrix[namespace]["*"]
				if !sameOps(resources[resource], wildcardOps) {
					allOpsSame = false
				}
			}
			if allOpsSame {
				for _, resources := range matrix {
					delete(resources, resource)
				}
				wasPruned = true
			}
		}
	}
	return wasPruned
}

func PrunePermissions(matrix PermissionMatrix) int {
	deleted := 0
	for namespace, resources := range matrix {
		for resource, operations := range resources {
			if resource != "*" {
				if sameOps(operations, matrix[namespace]["*"]) {
					delete(matrix[namespace], resource)
					deleted++
				}
			}
		}
	}
	return deleted
}

func sameOps(ops1 map[models.OperationType]struct{}, ops2 map[models.OperationType]struct{}) bool {
	if len(ops1) != len(ops2) {
		return false
	}
	for op := range ops1 {
		if _, exists := ops2[op]; !exists {
			return false
		}
	}
	return true
}

func toMatrix(role *models.Role, subroleMap map[string]*models.Role) PermissionMatrix {
	var matrix PermissionMatrix
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
		matrix = make(PermissionMatrix)
		matrix["*"] = make(map[string]map[models.OperationType]struct{})
		matrix["*"]["*"] = make(map[models.OperationType]struct{})
	}
	for _, permit := range role.Permit {
		addPermitToMatrix(matrix, permit)
	}
	for _, deny := range role.Deny {
		restrictMatrix(matrix, deny)
	}
	pruneResourcesNamespaces(matrix)
	return matrix
}

func addMatrix(m1 PermissionMatrix,
	m2 PermissionMatrix) PermissionMatrix {
	x1, y1 := len(m1), len(m1["*"])
	x2, y2 := len(m2), len(m2["*"])
	if x1*y1 < x2*y2 {
		m1, m2 = m2, m1
	}

	sum := deepCopy(m1)
	for namespace := range m2 {
		if _, exists := sum[namespace]; !exists {
			expandNamespaces(namespace, sum)
		}
	}
	for resource := range m2["*"] {
		if _, exists := sum["*"][resource]; !exists {
			expandResources(resource, sum)
		}
	}
	var fromNs, fromRes string
	for namespace, resources := range sum {
		if _, exists := m2[namespace]; !exists {
			fromNs = "*"
		} else {
			fromNs = namespace
		}
		for resource, operations := range resources {
			if _, exists := m2["*"][resource]; !exists {
				fromRes = "*"
			} else {
				fromRes = resource
			}
			for opType := range m2[fromNs][fromRes] {
				operations[opType] = struct{}{}
			}
		}
	}
	return sum
}

func addPermitToMatrix(matrix PermissionMatrix, permit models.Operation) {
	regulateMatrix(matrix, permit, addOp)
}

func restrictMatrix(matrix PermissionMatrix, deny models.Operation) {
	regulateMatrix(matrix, deny, removeOp)
}

func regulateMatrix(matrix PermissionMatrix, op models.Operation,
	action func(map[models.OperationType]struct{}, models.OperationType)) {
	if _, hasNamespace := matrix[op.Namespace]; !hasNamespace {
		expandNamespaces(op.Namespace, matrix)
	}
	if _, hasResource := matrix["*"][op.Resource]; !hasResource {
		expandResources(op.Resource, matrix)
	}

	if op.Namespace == "*" {
		if op.Resource == "*" { // * * op
			for _, resources := range matrix {
				for _, operations := range resources {
					action(operations, op.Type)
				}
			}
		} else { // * r op
			for _, resources := range matrix {
				action(resources[op.Resource], op.Type)
			}
		}
	} else {
		if op.Resource == "*" { // n * op
			for _, operations := range matrix[op.Namespace] {
				action(operations, op.Type)
			}
		} else { // n r op
			action(matrix[op.Namespace][op.Resource], op.Type)
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

func expandNamespaces(namespace string, matrix PermissionMatrix) {
	matrix[namespace] = make(map[string]map[models.OperationType]struct{})
	for resource, ops := range matrix["*"] {
		matrix[namespace][resource] = make(map[models.OperationType]struct{})
		for opType := range ops {
			matrix[namespace][resource][opType] = struct{}{}
		}
	}
}

func expandResources(resource string, matrix PermissionMatrix) {
	for _, namespaced := range matrix {
		namespaced[resource] = make(map[models.OperationType]struct{})
		for opType := range namespaced["*"] {
			namespaced[resource][opType] = struct{}{}
		}
	}
}