package auth

import (
	"container/list"
	"fmt"
	"log"
	"math"
	"math/rand"
	"sort"
	"testing"

	"github.com/ZPI-2024-25/KubernetesAccessManager/models"
	"github.com/stretchr/testify/assert"
)

func TestHasCycle(t *testing.T) {
	tests := []struct {
		name     string
		roleMap  map[string]*models.Role
		expected bool
	}{
		{
			name: "No cycle",
			roleMap: map[string]*models.Role{
				"admin":   {Subroles: []string{"user", "manager"}},
				"manager": {Subroles: []string{"user"}},
				"user":    {Subroles: []string{}},
			},
			expected: false,
		},
		{
			name: "Cycle exists",
			roleMap: map[string]*models.Role{
				"cyclic1": {Subroles: []string{"cyclic2"}},
				"cyclic2": {Subroles: []string{"cyclic1"}},
			},
			expected: true,
		},
		{
			name: "Complex no cycle",
			roleMap: map[string]*models.Role{
				"admin":   {Subroles: []string{"manager", "auditor"}},
				"manager": {Subroles: []string{"user"}},
				"user":    {Subroles: []string{}},
				"auditor": {Subroles: []string{"user"}},
			},
			expected: false,
		},
		{
			name: "Complex cycle",
			roleMap: map[string]*models.Role{
				"admin":   {Subroles: []string{"manager"}},
				"manager": {Subroles: []string{"auditor"}},
				"auditor": {Subroles: []string{"admin"}},
			},
			expected: true,
		},
		{
			name: "Self cycle",
			roleMap: map[string]*models.Role{
				"self": {Subroles: []string{"self"}},
			},
			expected: true,
		},
		{
			name: "Disconnected graph no cycle",
			roleMap: map[string]*models.Role{
				"admin":   {Subroles: []string{"manager"}},
				"manager": {Subroles: []string{"user"}},
				"user":    {Subroles: []string{}},
				"guest":   {Subroles: []string{"visitor"}},
				"visitor": {Subroles: []string{}},
			},
			expected: false,
		},
		{
			name: "Disconnected graph with cycle",
			roleMap: map[string]*models.Role{
				"admin":   {Subroles: []string{"manager"}},
				"manager": {Subroles: []string{"user"}},
				"user":    {Subroles: []string{}},
				"guest":   {Subroles: []string{"visitor"}},
				"visitor": {Subroles: []string{"guest"}},
			},
			expected: true,
		},
		{
			name: "Multiple cycles",
			roleMap: map[string]*models.Role{
				"role1": {Subroles: []string{"role2"}},
				"role2": {Subroles: []string{"role3"}},
				"role3": {Subroles: []string{"role1"}},
				"role4": {Subroles: []string{"role5"}},
				"role5": {Subroles: []string{"role6"}},
				"role6": {Subroles: []string{"role4"}},
			},
			expected: true,
		},
		{
			name: "Large acyclic graph",
			roleMap: map[string]*models.Role{
				"role1":  {Subroles: []string{"role2", "role3"}},
				"role2":  {Subroles: []string{"role4"}},
				"role3":  {Subroles: []string{"role4", "role5"}},
				"role4":  {Subroles: []string{"role6"}},
				"role5":  {Subroles: []string{"role6"}},
				"role6":  {Subroles: []string{"role7"}},
				"role7":  {Subroles: []string{"role8"}},
				"role8":  {Subroles: []string{"role9"}},
				"role9":  {Subroles: []string{"role10"}},
				"role10": {Subroles: []string{}},
			},
			expected: false,
		},
		{
			name: "Large cyclic graph",
			roleMap: map[string]*models.Role{
				"role1":  {Subroles: []string{"role2", "role3"}},
				"role2":  {Subroles: []string{"role4"}},
				"role3":  {Subroles: []string{"role4", "role5"}},
				"role4":  {Subroles: []string{"role6"}},
				"role5":  {Subroles: []string{"role6"}},
				"role6":  {Subroles: []string{"role7"}},
				"role7":  {Subroles: []string{"role8"}},
				"role8":  {Subroles: []string{"role9"}},
				"role9":  {Subroles: []string{"role10"}},
				"role10": {Subroles: []string{"role1"}},
			},
			expected: true,
		},
		{
			name: "Non-existing subrole",
			roleMap: map[string]*models.Role{
				"admin": {Subroles: []string{"nonexistent"}},
			},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := hasCycle(tt.roleMap)
			assert.Equal(t, tt.expected, result)
		})
	}
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

func TestHasPermission(t *testing.T) {
	tests := []struct {
		name       string
		roleMap    map[string]*models.Role
		subroleMap map[string]*models.Role
		rolename   string
		operation  *models.Operation
		expected   bool
	}{
		{
			name: "Permission granted",
			roleMap: map[string]*models.Role{
				"admin": {Subroles: []string{"user", "manager"}},
			},
			subroleMap: map[string]*models.Role{
				"user":    {Permit: []models.Operation{{Type: "read", Resource: "resource1", Namespace: "*"}}},
				"manager": {Subroles: []string{"user"}},
			},
			rolename:  "admin",
			operation: &models.Operation{Type: "read", Resource: "resource1", Namespace: "namespace"},
			expected:  true,
		},
		{
			name: "Permission denied",
			roleMap: map[string]*models.Role{
				"admin": {Subroles: []string{"user", "manager"}},
			},
			subroleMap: map[string]*models.Role{
				"user":    {Deny: []models.Operation{{Type: "read", Resource: "resource1", Namespace: "*"}}, Permit: []models.Operation{{Type: "read", Resource: "resource", Namespace: "hihi"}}},
				"manager": {Subroles: []string{"user"}},
			},
			rolename:  "admin",
			operation: &models.Operation{Type: "read", Resource: "resource1"},
			expected:  false,
		},
		{
			name: "Role not found",
			roleMap: map[string]*models.Role{
				"admin": {Subroles: []string{"user", "manager"}},
			},
			subroleMap: map[string]*models.Role{
				"user":    {Permit: []models.Operation{{Type: "read", Resource: "resource1"}}},
				"manager": {Subroles: []string{"user"}},
			},
			rolename:  "nonexistent",
			operation: &models.Operation{Type: "read", Resource: "resource1"},
			expected:  false,
		},
		{
			name: "Permission granted from alternative path",
			roleMap: map[string]*models.Role{
				"0": {Subroles: []string{"1", "2"}},
			},
			subroleMap: map[string]*models.Role{
				"1": {Deny: []models.Operation{{Type: "read", Resource: "resource1", Namespace: "*"}}, Subroles: []string{"3"}},
				"2": {Subroles: []string{"3"}},
				"3": {Permit: []models.Operation{{Type: "read", Resource: "resource1", Namespace: "*"}}},
			},
			rolename:  "0",
			operation: &models.Operation{Type: "read", Resource: "resource1", Namespace: "namespace"},
			expected:  true,
		},
		{
			name: "Permission granted from alternative path",
			roleMap: map[string]*models.Role{
				"0": {Subroles: []string{"1", "2"}},
			},
			subroleMap: map[string]*models.Role{
				"1": {Deny: []models.Operation{{Type: "read", Resource: "resource1", Namespace: "*"}}, Subroles: []string{"3"}},
				"2": {Subroles: []string{"3"}},
				"3": {Permit: []models.Operation{{Type: "read", Resource: "resource1", Namespace: "*"}}},
			},
			rolename:  "0",
			operation: &models.Operation{Type: "read", Resource: "resource1", Namespace: "namespace"},
			expected:  true,
		},
		{
			name: "Nil subrole map",
			roleMap: map[string]*models.Role{
				"admin": {Subroles: []string{"user", "manager"}},
			},
			subroleMap: nil,
			rolename:   "admin",
			operation:  &models.Operation{Type: "read", Namespace: "n1", Resource: "resource1"},
			expected:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rmr := &RoleMapRepository{
				RoleMap:    tt.roleMap,
				SubroleMap: tt.subroleMap,
			}
			visited := make(map[string]struct{})
			role := rmr.RoleMap[tt.rolename]
			result := hasPermission(role, rmr.SubroleMap, tt.operation, visited)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestFlatHasPermission(t *testing.T) {
	tests := []struct {
		name   string
		matrix PermissionMatrix
		tests  map[*models.Operation]bool
	}{
		{
			name: "Permission granted specific namespace and resource",
			matrix: PermissionMatrix{
				"namespace1": {
					"resource1": {
						"read": struct{}{},
					},
				},
			},
			tests: map[*models.Operation]bool{
				{Type: "read", Resource: "resource1", Namespace: "namespace1"}: true,
			},
		},
		{
			name: "Permission granted wildcard namespace",
			matrix: PermissionMatrix{
				"*": {
					"resource1": {
						"read": struct{}{},
					},
				},
			},
			tests: map[*models.Operation]bool{
				{Type: "read", Resource: "resource1", Namespace: "namespace1"}: true,
			},
		},
		{
			name: "Permission granted wildcard resource",
			matrix: PermissionMatrix{
				"namespace1": {
					"*": {
						"read": struct{}{},
					},
				},
			},
			tests: map[*models.Operation]bool{
				{Type: "read", Resource: "resource1", Namespace: "namespace1"}: true,
			},
		},
		{
			name: "Permission denied no matching namespace",
			matrix: PermissionMatrix{
				"namespace2": {
					"resource1": {
						"read": struct{}{},
					},
				},
			},
			tests: map[*models.Operation]bool{
				{Type: "read", Resource: "resource1", Namespace: "namespace1"}: false,
			},
		},
		{
			name: "Permission denied no matching resource",
			matrix: PermissionMatrix{
				"namespace1": {
					"resource2": {
						"read": struct{}{},
					},
				},
			},
			tests: map[*models.Operation]bool{
				{Type: "read", Resource: "resource1", Namespace: "namespace1"}: false,
			},
		},
		{
			name: "Permission denied no matching operation type",
			matrix: PermissionMatrix{
				"namespace1": {
					"resource1": {
						"write": struct{}{},
					},
				},
			},
			tests: map[*models.Operation]bool{
				{Type: "read", Resource: "resource1", Namespace: "namespace1"}: false,
			},
		},
		{
			name: "Test all restrictions of wildcards",
			matrix: PermissionMatrix{
				"*": {
					"*":  {"create": struct{}{}, "read": struct{}{}, "write": struct{}{}},
					"r1": {"create": struct{}{}, "read": struct{}{}},
				},
				"n1": {
					"*":  {"read": struct{}{}, "write": struct{}{}},
					"r1": {"list": struct{}{}},
				},
			},
			tests: map[*models.Operation]bool{
				// n1, r1
				{Type: "create", Resource: "r1", Namespace: "n1"}: false,
				{Type: "read", Resource: "r1", Namespace: "n1"}:   false,
				{Type: "write", Resource: "r1", Namespace: "n1"}:  false,
				{Type: "list", Resource: "r1", Namespace: "n1"}:   true,
				// n1, r2 = *
				{Type: "create", Resource: "r2", Namespace: "n1"}: false,
				{Type: "read", Resource: "r2", Namespace: "n1"}:   true,
				{Type: "write", Resource: "r2", Namespace: "n1"}:  true,
				{Type: "list", Resource: "r2", Namespace: "n1"}:   false,
				// n2 = *, r1
				{Type: "create", Resource: "r1", Namespace: "n2"}: true,
				{Type: "read", Resource: "r1", Namespace: "n2"}:   true,
				{Type: "write", Resource: "r1", Namespace: "n2"}:  false,
				{Type: "list", Resource: "r1", Namespace: "n2"}:   false,
				// n2 = *, r2 = *
				{Type: "create", Resource: "r2", Namespace: "n2"}: true,
				{Type: "read", Resource: "r2", Namespace: "n2"}:   true,
				{Type: "write", Resource: "r2", Namespace: "n2"}:  true,
				{Type: "list", Resource: "r2", Namespace: "n2"}:   false,
			},
		},
		// Test case from fuzy:
		// Namespace: *
		//   Resource: *
		//   Resource: r1
		//   Resource: r0
		//     Operation: delete
		//   Resource: r2
		// Namespace: n2
		//   Resource: *
		//   Resource: r1
		//     Operation: delete
		//     Operation: update
		//   Resource: r0
		//     Operation: delete
		//   Resource: r2
		// Namespace: n0
		//   Resource: *
		//   Resource: r1
		//   Resource: r0
		//     Operation: delete
		//   Resource: r2
		// Namespace: n1
		//   Resource: r2
		//   Resource: *
		//     Operation: delete
		//   Resource: r1
		//     Operation: delete
		//   Resource: r0
		//     Operation: delete
		{
			name: "Test case from fuzzy",
			matrix: PermissionMatrix{
				"*": {
					"*":  {},
					"r1": {},
					"r0": {"delete": {}},
					"r2": {},
				},
				"n2": {
					"*":  {},
					"r1": {"delete": {}, "update": {}},
					"r0": {"delete": {}},
					"r2": {},
				},
				"n0": {
					"*":  {},
					"r1": {},
					"r0": {"delete": {}},
					"r2": {},
				},
				"n1": {
					"r2": {},
					"*":  {"delete": {}},
					"r1": {"delete": {}},
					"r0": {"delete": {}},
				},
			},
			tests: map[*models.Operation]bool{
				{Type: "delete", Resource: "r0", Namespace: "n1"}: true,
				{Type: "delete", Resource: "r1", Namespace: "n1"}: true,
				{Type: "delete", Resource: "r2", Namespace: "n1"}: false,
				{Type: "delete", Resource: "r3", Namespace: "n1"}: true,
				{Type: "delete", Resource: "r0", Namespace: "n2"}: true,
				{Type: "delete", Resource: "r0", Namespace: "n3"}: true,
				{Type: "update", Resource: "r1", Namespace: "n2"}: true,
				{Type: "delete", Resource: "r1", Namespace: "n2"}: true,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			for operation, expected := range tt.tests {
				result := flatHasPermission(operation, tt.matrix)
				if expected == result {
					assert.Equal(t, expected, result)
				} else {
					t.Errorf("Test failed for operation: %v", operation)
					assert.Equal(t, expected, result)
				}
			}
		})
	}
}

func TestRoleMapFlattening(t *testing.T) {
	tests := []struct {
		name       string
		roleMap    map[string]*models.Role
		subroleMap map[string]*models.Role
	}{
		{
			name: "Nil subrole map",
			roleMap: map[string]*models.Role{
				"root": {
					Permit: []models.Operation{
						{Resource: "r1", Type: "create", Namespace: "*"},
						{Resource: "r0", Type: "read", Namespace: "n1"},
					},
					Deny: []models.Operation{
						{Resource: "r1", Type: "*", Namespace: "n0"},
					},
					Name: "root", Subroles: []string{"user", "manager"}},
			},
			subroleMap: nil,
		},
		{
			name: "Test case from fuzzy",
			subroleMap: map[string]*models.Role{
				"sub1": {
					Subroles: []string{"sub4"},
					Permit:   []models.Operation{},
					Deny: []models.Operation{
						{Resource: "r2", Type: "*", Namespace: "n0"},
					},
				},
				"sub2": {
					Subroles: []string{},
					Permit: []models.Operation{
						{Resource: "r1", Type: "create", Namespace: "*"},
						{Resource: "r0", Type: "read", Namespace: "n1"},
					},
					Deny: []models.Operation{},
				},
				"sub3": {
					Subroles: []string{"sub4"},
					Permit: []models.Operation{
						{Resource: "r1", Type: "update", Namespace: "n2"},
						{Resource: "*", Type: "delete", Namespace: "n1"},
					},
					Deny: []models.Operation{},
				},
				"sub4": {
					Subroles: []string{},
					Permit: []models.Operation{
						{Resource: "r1", Type: "delete", Namespace: "n2"},
						{Resource: "r0", Type: "delete", Namespace: "*"},
					},
					Deny: []models.Operation{},
				},
				"sub0": {
					Subroles: []string{"sub3", "sub1"},
					Permit:   []models.Operation{},
					Deny:     []models.Operation{},
				},
			},
			roleMap: map[string]*models.Role{
				"root": {
					Subroles: []string{"sub0"},
					Name:     "root",
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			flattened := createPermissionMatrix(tt.roleMap, tt.subroleMap)
			for _, ns := range []string{"n0", "n1", "n2", "n3", "n4"} {
				for _, rs := range []string{"r0", "r1", "r2", "r3", "r4"} {
					for _, opType := range []models.OperationType{"create", "read", "update", "delete", "list"} {
						op := &models.Operation{
							Namespace: ns,
							Type:      opType,
							Resource:  rs,
						}
						expected := hasPermission(tt.roleMap["root"], tt.subroleMap, op, make(map[string]struct{}))
						result := flatHasPermission(op, flattened["root"])
						if expected != result {
							t.Errorf("Test failed for operation %v. Expected: %v, Actual: %v", op, expected, result)
						}
					}
				}
			}
		})
	}
}

func TestRoleMapFlatteningFuzzy(t *testing.T) {
	operations := []models.OperationType{"create", "read", "update", "delete"}
	tests := []struct {
		name              string
		permits           int
		denies            int
		resources         int
		namespaces        int
		pstar             float64
		dstar             float64
		nodes             int
		edges             int
		allowedOperations []models.OperationType
	}{
		{
			name:              "Small graph with few permits and denies",
			permits:           10,
			denies:            5,
			resources:         3,
			namespaces:        3,
			pstar:             0.2,
			dstar:             0.1,
			nodes:             5,
			edges:             4,
			allowedOperations: []models.OperationType{"create", "read", "update", "delete", "*"},
		},
		{
			name:              "Medium graph with moderate permits and denies",
			permits:           50,
			denies:            25,
			resources:         10,
			namespaces:        5,
			pstar:             0.3,
			dstar:             0.2,
			nodes:             10,
			edges:             15,
			allowedOperations: []models.OperationType{"create", "read", "update", "delete", "*"},
		},
		{
			name:              "Large graph with many permits and denies",
			permits:           100,
			denies:            50,
			resources:         20,
			namespaces:        10,
			pstar:             0.4,
			dstar:             0.3,
			nodes:             20,
			edges:             30,
			allowedOperations: []models.OperationType{"create", "read", "update", "delete", "*"},
		},
		{
			name:              "Large graph with many permits and denies",
			permits:           2500,
			denies:            1000,
			resources:         40,
			namespaces:        100,
			pstar:             0.2,
			dstar:             0.15,
			nodes:             100,
			edges:             200,
			allowedOperations: []models.OperationType{"create", "read", "update", "delete", "list", "*"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			log.Printf("Creating graph")
			graph := randomAcyclicGraph(tt.nodes, tt.edges)
			log.Printf("Graph created")
			repo := createFuzzyRoleRepo(tt.permits, tt.denies, tt.resources, tt.namespaces, tt.pstar, tt.dstar, tt.allowedOperations, graph)
			log.Printf("Role map created")
			assert.NotNil(t, repo)

			// Check if the flattened map is created correctly
			repo.flattenedMap = createPermissionMatrix(repo.RoleMap, repo.SubroleMap)
			log.Printf("Flattened map created")

			allPassed := true
			// Compare hasPermission with flatHasPermission
			for i := 0; i < tt.namespaces+1; i++ {
				for j := 0; j < tt.resources+1; j++ {
					for _, op := range operations {
						operation := &models.Operation{Type: op, Resource: fmt.Sprintf("r%d", j), Namespace: fmt.Sprintf("n%d", i)}
						expected := hasPermission(repo.RoleMap["root"], repo.SubroleMap, operation, make(map[string]struct{}))
						result := flatHasPermission(operation, repo.flattenedMap["root"])
						assert.Equal(t, expected, result)
						if expected != result {
							log.Printf("Test failed for operation: %v", operation)
							allPassed = false
						}
					}
				}
			}
			if !allPassed {
				// print the matrix for debugging
				log.Printf("Matrix:")
				for k, v := range repo.flattenedMap {
					log.Printf("Role: %s", k)
					for k2, v2 := range v {
						log.Printf("  Namespace: %s", k2)
						for k3, v3 := range v2 {
							log.Printf("    Resource: %s", k3)
							for k4 := range v3 {
								log.Printf("      Operation: %s", k4)
							}
						}
					}
				}
				// print the role map for debugging
				log.Printf("Role map:")
				for k, v := range repo.RoleMap {
					log.Printf("Role: %s", k)
					log.Printf("  Subroles: %v", v.Subroles)
					log.Printf("  Permit: %v", v.Permit)
					log.Printf("  Deny: %v", v.Deny)
				}
				log.Printf("Subrole map:")
				// make a list of [(subroleName, subrole)] from repo.SubroleMap, sort it by subroleName, and print it
				subroleList := list.New()
				for k, v := range repo.SubroleMap {
					subroleList.PushBack(struct {
						key   string
						value *models.Role
					}{k, v})
				}
				// Convert list to slice for sorting
				subroleSlice := make([]struct {
					key   string
					value *models.Role
				}, subroleList.Len())
				i := 0
				for e := subroleList.Front(); e != nil; e = e.Next() {
					subroleSlice[i] = e.Value.(struct {
						key   string
						value *models.Role
					})
					i++
				}
				sort.Slice(subroleSlice, func(i, j int) bool {
					return subroleSlice[i].key < subroleSlice[j].key
				})
				for _, subrole := range subroleSlice {
					log.Printf("Subrole: %s", subrole.key)
					log.Printf("  Subroles: %v", subrole.value.Subroles)
					log.Printf("  Permit: %v", subrole.value.Permit)
					log.Printf("  Deny: %v", subrole.value.Deny)
				}
			}
		})
	}
}

func TestPruning(t *testing.T) {
	operations := []models.OperationType{"create", "read", "update", "delete"}
	tests := []struct {
		name              string
		permits           int
		denies            int
		resources         int
		namespaces        int
		pstar             float64
		dstar             float64
		nodes             int
		edges             int
		allowedOperations []models.OperationType
	}{
		{
			name:              "Small graph with few permits and denies",
			permits:           10,
			denies:            5,
			resources:         3,
			namespaces:        3,
			pstar:             0.2,
			dstar:             0.1,
			nodes:             5,
			edges:             4,
			allowedOperations: []models.OperationType{"create", "read", "update", "delete", "*"},
		},
		{
			name:              "Medium graph with moderate permits and denies",
			permits:           50,
			denies:            25,
			resources:         10,
			namespaces:        5,
			pstar:             0.3,
			dstar:             0.2,
			nodes:             10,
			edges:             15,
			allowedOperations: []models.OperationType{"create", "read", "update", "delete", "*"},
		},
		{
			name:              "Large graph with many permits and denies",
			permits:           100,
			denies:            50,
			resources:         20,
			namespaces:        10,
			pstar:             0.4,
			dstar:             0.3,
			nodes:             20,
			edges:             30,
			allowedOperations: []models.OperationType{"create", "read", "update", "delete", "*"},
		},
		{
			name:              "Large graph with many permits and denies",
			permits:           2500,
			denies:            1000,
			resources:         40,
			namespaces:        100,
			pstar:             0.2,
			dstar:             0.15,
			nodes:             100,
			edges:             200,
			allowedOperations: []models.OperationType{"create", "read", "update", "delete", "list", "*"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			graph := randomAcyclicGraph(tt.nodes, tt.edges)
			repo := createFuzzyRoleRepo(tt.permits, tt.denies, tt.resources, tt.namespaces, tt.pstar, tt.dstar, tt.allowedOperations, graph)
			assert.NotNil(t, repo)

			// Check if the flattened map is created correctly
			repo.flattenedMap = createPermissionMatrix(repo.RoleMap, repo.SubroleMap)
			ogSize := len(repo.flattenedMap["root"]) * len(repo.flattenedMap["root"]["*"])
			pruned := PrunePermissions(repo.flattenedMap["root"])
			log.Printf("Flattened map created and pruned %d out of %d", pruned, ogSize)

			// Compare hasPermission with flatHasPermission
			for i := 0; i < tt.namespaces+1; i++ {
				for j := 0; j < tt.resources+1; j++ {
					for _, op := range operations {
						operation := &models.Operation{Type: op, Resource: fmt.Sprintf("r%d", j), Namespace: fmt.Sprintf("n%d", i)}
						expected := hasPermission(repo.RoleMap["root"], repo.SubroleMap, operation, make(map[string]struct{}))
						result := flatHasPermission(operation, repo.flattenedMap["root"])
						assert.Equal(t, expected, result)
						if expected != result {
							log.Printf("Test failed for operation: %v", operation)
						}
					}
				}
			}
		})
	}
}

func createFuzzyRoleRepo(permits, denies, resources, namespaces int, pstar, dstar float64, operations []models.OperationType, graph map[string][]string) *RoleMapRepository {
	roleMap := map[string]*models.Role{"root": {Name: "root", Subroles: []string{"sub0"}}}
	subroleMap := make(map[string]*models.Role)
	subroles := len(graph)
	avgP := permits / subroles
	avgD := denies / subroles

	for node, children := range graph {
		subrole := &models.Role{
			Permit:   make([]models.Operation, 0),
			Deny:     make([]models.Operation, 0),
			Subroles: children,
		}
		subroleMap[node] = subrole

		pnr := rand.Intn(avgP*2 + 1)
		dnr := rand.Intn(avgD*2 + 1)
		permits -= pnr
		denies -= dnr
		avgP = permits / subroles
		avgD = denies / subroles
		for j := 0; j < pnr; j++ {
			op := makeRandomOperation(pstar, dstar, namespaces, resources, operations)
			subrole.Permit = append(subrole.Permit, *op)
		}
		for j := 0; j < dnr; j++ {
			op := makeRandomOperation(pstar, dstar, namespaces, resources, operations)
			subrole.Deny = append(subrole.Deny, *op)
		}
	}

	return &RoleMapRepository{
		RoleMap:    roleMap,
		SubroleMap: subroleMap,
	}
}

func makeRandomOperation(pstar, dstar float64, namespaces, resources int, operations []models.OperationType) *models.Operation {
	wildcard := "*"
	var n, r string
	if rand.Float64() < pstar {
		n = wildcard
	} else {
		n = fmt.Sprintf("n%d", rand.Intn(namespaces))
	}
	if rand.Float64() < dstar {
		r = wildcard
	} else {
		r = fmt.Sprintf("r%d", rand.Intn(resources))
	}
	op := models.Operation{
		Type:      operations[rand.Intn(len(operations))],
		Resource:  r,
		Namespace: n,
	}
	return &op
}

func randomAcyclicGraph(nodes, edges int) map[string][]string {
	graph := make(map[string]map[string]struct{})
	if edges > (nodes * (nodes - 1) / 3) {
		edges = nodes * (nodes - 1) / 3
	}
	for i := 0; i < nodes; i++ {
		graph[fmt.Sprintf("sub%d", i)] = make(map[string]struct{})
	}
	for edges > 0 {
		nr := rand.Float64()
		skewed := math.Pow(nr, 2) // skew the distribution towards head
		from := int(skewed * float64(nodes-1))
		to := rand.Intn(nodes-1-from) + from + 1
		// log.Printf("from: %d, to: %d", from, to)
		if _, ok := graph[fmt.Sprintf("sub%d", from)][fmt.Sprintf("sub%d", to)]; !ok {
			edges = edges - 1
			graph[fmt.Sprintf("sub%d", from)][fmt.Sprintf("sub%d", to)] = struct{}{}
		}
	}
	graphAsList := make(map[string][]string)
	for k, v := range graph {
		graphAsList[k] = make([]string, 0)
		for k2 := range v {
			graphAsList[k] = append(graphAsList[k], k2)
		}
	}
	return graphAsList
}

func TestDeepCopy(t *testing.T) {
	tests := []struct {
		matrice PermissionMatrix
	}{
		{
			matrice: PermissionMatrix{
				"n1": {
					"r1": {
						"create": struct{}{},
						"read":   struct{}{},
					},
					"r2": {
						"delete": struct{}{},
					},
				},
				"n2": {
					"r1": {
						"update": struct{}{},
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run("Deep copy", func(t *testing.T) {
			copy := deepCopy(tt.matrice)
			assert.EqualValues(t, tt.matrice, copy)
			tt.matrice["n1"]["r1"]["update"] = struct{}{}
			assert.NotEqual(t, tt.matrice, copy)
		})
	}
}
