package auth

import (
	"regexp"
	"strings"
	"testing"
	"github.com/golang-jwt/jwt/v4"
	"github.com/stretchr/testify/assert"
)


func TestAlternativeWayOfExtractingRoles(t *testing.T) {
	claimsStr := `{
  "exp": 1730123468,
  "iat": 1730123168,
  "auth_time": 1730123100,
  "jti": "df306998-45d3-4a4d-918e-a3d9a2037938",
  "iss": "http://localhost:8081/realms/access-manager",
  "aud": "account",
  "sub": "dd967421-a04e-4c3a-a74c-57e483dad1a8",
  "typ": "Bearer",
  "azp": "account-console",
  "sid": "4586f354-45f6-4e36-a87b-a1aa7f5cd873",
  "acr": "0",
  "resource_access": {
    "account-console": {
      "roles": [
        "pod-reader"
      ]
    },
    "account": {
      "roles": [
        "manage-account",
        "manage-account-links"
      ]
    }
  },
  "scope": "openid profile email",
  "email_verified": false,
  "name": "Marek Fiuk",
  "preferred_username": "marefek1@gmail.com",
  "given_name": "Marek",
  "family_name": "Fiuk",
  "email": "marefek1@gmail.com"
}`
t.Run("TestExtractRoles", func(t *testing.T) {
		// Regex to find all roles within resource_access
		re := regexp.MustCompile(`"roles":\s*\[([^\]]+)\]`)
		matches := re.FindAllStringSubmatch(claimsStr, -1)

		var roles []string
		for _, match := range matches {
			if len(match) > 1 {
				// Extract roles from the matched group
				roleItems := strings.Split(match[1], ",")
				for _, role := range roleItems {
					// Clean up leading/trailing whitespace and quotes around each role
					role = strings.TrimSpace(role)
					role = strings.Trim(role, `"`)
					roles = append(roles, role)
				}
			}
		}
		// Assert that the roles match the expected output
		assert.Equal(t, []string{"pod-reader", "manage-account", "manage-account-links"}, roles)
	})
}

func TestJsonTokenRoleExtraction(t *testing.T) {
	t.Run("TestJsonToken", func(t *testing.T) {
		tokenStr := `eyJhbGciOiJSUzI1NiIsInR5cCIgOiAiSldUIiwia2lkIiA6ICI3SERLbTBsSHJLY18ybHc0eFo1S0NBR0JObndCTDJsOUlucFJ5VVU4ZHBjIn0.eyJleHAiOjE3MzAxMjM0NjgsImlhdCI6MTczMDEyMzE2OCwiYXV0aF90aW1lIjoxNzMwMTIzMTAwLCJqdGkiOiJkZjMwNjk5OC00NWQzLTRhNGQtOTE4ZS1hM2Q5YTIwMzc5MzgiLCJpc3MiOiJodHRwOi8vbG9jYWxob3N0OjgwODEvcmVhbG1zL2FjY2Vzcy1tYW5hZ2VyIiwiYXVkIjoiYWNjb3VudCIsInN1YiI6ImRkOTY3NDIxLWEwNGUtNGMzYS1hNzRjLTU3ZTQ4M2RhZDFhOCIsInR5cCI6IkJlYXJlciIsImF6cCI6ImFjY291bnQtY29uc29sZSIsInNpZCI6IjQ1ODZmMzU0LTQ1ZjYtNGUzNi1hODdiLWExYWE3ZjVjZDg3MyIsImFjciI6IjAiLCJyZXNvdXJjZV9hY2Nlc3MiOnsiYWNjb3VudC1jb25zb2xlIjp7InJvbGVzIjpbInBvZC1yZWFkZXIiXX0sImFjY291bnQiOnsicm9sZXMiOlsibWFuYWdlLWFjY291bnQiLCJtYW5hZ2UtYWNjb3VudC1saW5rcyJdfX0sInNjb3BlIjoib3BlbmlkIHByb2ZpbGUgZW1haWwiLCJlbWFpbF92ZXJpZmllZCI6ZmFsc2UsIm5hbWUiOiJNYXJlayBGaXVrIiwicHJlZmVycmVkX3VzZXJuYW1lIjoibWFyZWZlazFAZ21haWwuY29tIiwiZ2l2ZW5fbmFtZSI6Ik1hcmVrIiwiZmFtaWx5X25hbWUiOiJGaXVrIiwiZW1haWwiOiJtYXJlZmVrMUBnbWFpbC5jb20ifQ.zVe1FBnNkx7OlYveHZVG9vNJqwEJTtua5rDFekFJ9sNFAXK7e-xahcuEoOAy4_YTAjfGtgvQMHq2hy61_30Xe1cp6okmH0YnXZ-w4WXaxKdB7tHNcpduFiQSeCFBp4COImTEyuvOqv4PjLjLu5N0wkyfXClhoTIjvn932e_QEpeAjCeG5nDTePk3SqDbVYKo3cK0Ymzap7U4-H1OmM_YGPoYTGzC1Qri2rspPtfoaFP3Uv3jYUmGA1dl8_b90QDRalOq8AZxzrnTJbm1VfHH0tbEfUZqQV8ok_Wjf7PQ27M8dajkXcYDNneFoCVlaFwrfXJcJDdFfOvTS4ryRy1ZyA`
		claims := jwt.MapClaims{}
		_, _ = jwt.ParseWithClaims(tokenStr, &claims, nil)
		roles, _ := ExtractRoles(&claims)
		expectedRoles := []string{"manage-account", "manage-account-links"}
		assert.ElementsMatch(t, expectedRoles, roles)
	})
}
func TestExtractRoles(t *testing.T) {
	t.Run("TestExtractRolesWithResourceAccess", func(t *testing.T) {
		claims := jwt.MapClaims{
			"resource_access": map[string]interface{}{
				"account-console": map[string]interface{}{
					"roles": []interface{}{"pod-reader"},
				},
				"account": map[string]interface{}{
					"roles": []interface{}{"manage-account", "manage-account-links"},
				},
			},
		}
		expectedRoles := []string{ "manage-account", "manage-account-links"}
		roles, _ := ExtractRoles(&claims)
		assert.ElementsMatch(t, expectedRoles, roles)
	})

	t.Run("TestExtractRolesWithRealmAccess", func(t *testing.T) {
		claims := jwt.MapClaims{
			"realm_access": map[string]interface{}{
				"roles": []interface{}{"admin", "user"},
			},
		}
		expectedRoles := []string{"admin", "user"}
		roles, _ := ExtractRoles(&claims)
		assert.ElementsMatch(t, expectedRoles, roles)
	})

	t.Run("TestExtractRolesWithBothAccess", func(t *testing.T) {
		claims := jwt.MapClaims{
			"realm_access": map[string]interface{}{
				"roles": []interface{}{"admin", "user"},
			},
			"resource_access": map[string]interface{}{
				"account-console": map[string]interface{}{
					"roles": []interface{}{"pod-reader"},
				},
				"account": map[string]interface{}{
					"roles": []interface{}{"manage-account", "manage-account-links"},
				},
			},
		}
		expectedRoles := []string{"admin", "user", "manage-account", "manage-account-links"}
		roles, _ := ExtractRoles(&claims)
		assert.ElementsMatch(t, expectedRoles, roles)
	})

	t.Run("TestExtractRolesWithNoRoles", func(t *testing.T) {
		claims := jwt.MapClaims{}
		expectedRoles := []string{}
		roles, _ := ExtractRoles(&claims)
		assert.ElementsMatch(t, expectedRoles, roles)
	})
}