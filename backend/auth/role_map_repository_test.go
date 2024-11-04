package auth

import (
	"testing"

	"github.com/ZPI-2024-25/KubernetesAccessManager/models"
	"github.com/stretchr/testify/assert"
)

func TestHasCycle(t *testing.T) {
	// Define test cases
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
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := hasCycle(tt.roleMap)
			assert.Equal(t, tt.expected, result)
		})
	}
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
	}

	for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            rmr := &RoleMapRepository{
                RoleMap:    tt.roleMap,
                SubroleMap: tt.subroleMap,
            }
            visited := make(map[string]interface{})
            role := rmr.RoleMap[tt.rolename]
            result := rmr.hasPermission(role, tt.operation, visited)
            assert.Equal(t, tt.expected, result)
        })
    }
}
