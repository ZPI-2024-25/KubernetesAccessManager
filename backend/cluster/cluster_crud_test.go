package cluster

import (
	"context"
	"errors"
	"regexp"
	"strings"
	"testing"
	// "github.com/MicahParks/keyfunc"
	"github.com/ZPI-2024-25/KubernetesAccessManager/models"
	"github.com/golang-jwt/jwt/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/client-go/dynamic"
)

type MockResourceInterface struct {
	dynamic.ResourceInterface
	mock.Mock
	ReturnedValue *unstructured.Unstructured
	ReturnedError error
}

func (m *MockResourceInterface) Get(ctx context.Context, name string, 
	options metav1.GetOptions, subresources ...string) (*unstructured.Unstructured, error) {
	return m.ReturnedValue, m.ReturnedError
}

func (m *MockResourceInterface) Create(ctx context.Context, obj *unstructured.Unstructured, 
	options metav1.CreateOptions, subresources ...string) (*unstructured.Unstructured, error) {
	return m.ReturnedValue, m.ReturnedError
}

func (m *MockResourceInterface) Delete(ctx context.Context, name string, 
	options metav1.DeleteOptions, subresources ...string) error {
	return m.ReturnedError
}

func (m *MockResourceInterface) Update(ctx context.Context, obj *unstructured.Unstructured, 
	options metav1.UpdateOptions, subresources ...string) (*unstructured.Unstructured, error) {
	return m.ReturnedValue, m.ReturnedError
}

func MockResourceDetailsUnstructured(mockResource map[string]interface{}) *models.ResourceDetails {
	var castedResourceDetails interface{} = &unstructured.Unstructured{Object: mockResource}
	return &models.ResourceDetails{ResourceDetails: &castedResourceDetails}
}

func MockResourceDetailsMap(mockResource map[string]interface{}) *models.ResourceDetails {
	var castedResourceDetails interface{} = mockResource
	return &models.ResourceDetails{ResourceDetails: &castedResourceDetails}
}

func MockResourceDetails() *models.ResourceDetails {
	return MockResourceDetailsUnstructured(map[string]interface{}{
		"key": "value",
		"namespace": "validNamespace",
	})}

func MockUnstructured() *unstructured.Unstructured {
	return &unstructured.Unstructured{
		Object: map[string]interface{}{
			"key": "value",
			"namespace": "validNamespace",
		}}}


func TestGetResourceError(t *testing.T) {

	mockGetResourceI := new(mock.Mock)
	mockGetResourceI.On("func1", "Pod", "validNamespace", DefaultNamespace).Return(&MockResourceInterface{}, &models.ModelError{Code: 404, Message: "Not found"})
	getResourceI := func(resourceType string, namespace string, emptyNamespace string) (dynamic.ResourceInterface, *models.ModelError) {
		args := mockGetResourceI.Called(resourceType, namespace, emptyNamespace)
		return args.Get(0).(dynamic.ResourceInterface), args.Get(1).(*models.ModelError)
	}

	t.Run("Test GetResource", func(t *testing.T) {
		_, err := GetResource("Pod", "validNamespace", "validName", getResourceI)
		assert.NotNil(t, err)
		assert.Equal(t, &models.ModelError{Code: 404, Message: "Not found"}, err)
	})
}

func TestGetResourceSuccess(t *testing.T) {

	mockGetResourceI := new(mock.Mock)
	mockGetResourceI.On("func1", "Pod", "validNamespace", DefaultNamespace).
		Return(&MockResourceInterface{ReturnedValue: MockUnstructured()}, nil)
	getResourceI := func(resourceType string, namespace string, emptyNamespace string) (dynamic.ResourceInterface, *models.ModelError) {
		args := mockGetResourceI.Called(resourceType, namespace, emptyNamespace)
		return args.Get(0).(dynamic.ResourceInterface), nil
	}

	t.Run("Test GetResource", func(t *testing.T) {
		result, err := GetResource("Pod", "validNamespace", "validName", getResourceI)
		assert.Nil(t, err)
		expected := map[string]interface{}{
			"key": "value",
			"namespace": "validNamespace",
		}
		resultObj := (*result.ResourceDetails).(*unstructured.Unstructured)
		for key, value := range expected {
			assert.Equal(t, value, resultObj.Object[key])
		}
	})
}

func TestGetResourceErrorFromGet(t *testing.T) {
	mockGetResourceI := new(mock.Mock)
	mockGetResourceI.On("func1", "Pod", "validNamespace", DefaultNamespace).
		Return(&MockResourceInterface{ReturnedValue: MockUnstructured(), ReturnedError: errors.New("error")}, nil)
	getResourceI := func(resourceType string, namespace string, 
		emptyNamespace string) (dynamic.ResourceInterface, *models.ModelError) {
		args := mockGetResourceI.Called(resourceType, namespace, emptyNamespace)
		return args.Get(0).(dynamic.ResourceInterface), nil
	}
	t.Run("Test GetResourceErrorFromGet", func(t *testing.T) {
		res, err := GetResource("Pod", "validNamespace", "validName", getResourceI)

		assert.NotNil(t, err)
		assert.EqualValues(t, err.Code, 500)
		assert.EqualValues(t, err.Message, "Internal server error: error")
		assert.EqualValues(t, res, models.ResourceDetails{})
	})
}

func TestCreateResourceError(t *testing.T) {
	mockGetResourceI := new(mock.Mock)
	mockGetResourceI.On("func1", "Pod", "validNamespace", DefaultNamespace).Return(&MockResourceInterface{}, &models.ModelError{Code: 404, Message: "Not found"})
	getResourceI := func(resourceType string, namespace string, 
		emptyNamespace string) (dynamic.ResourceInterface, *models.ModelError) {
		args := mockGetResourceI.Called(resourceType, namespace, emptyNamespace)
		return args.Get(0).(dynamic.ResourceInterface), args.Get(1).(*models.ModelError)
	}

	t.Run("Test CreateResourceError", func(t *testing.T) {
		_, err := CreateResource("Pod", "validNamespace", models.ResourceDetails{}, getResourceI)
		assert.NotNil(t, err)
		assert.Equal(t, &models.ModelError{Code: 404, Message: "Not found"}, err)
	})
}

func TestCreateResourceSuccess(t *testing.T) {
	mockGetResourceI := new(mock.Mock)
	mockGetResourceI.On("func1", "Pod", "validNamespace", DefaultNamespace).
		Return(&MockResourceInterface{ReturnedValue: MockUnstructured()}, nil)
	getResourceI := func(resourceType string, namespace string, emptyNamespace string) (dynamic.ResourceInterface, *models.ModelError) {
		args := mockGetResourceI.Called(resourceType, namespace, emptyNamespace)
		return args.Get(0).(dynamic.ResourceInterface), nil
	}

	t.Run("Test CreateResource", func(t *testing.T) {
		result, err := CreateResource("Pod", "validNamespace", *MockResourceDetailsMap(map[string]interface{}{
			"key": "value", 
			"namespace": "validNamespace",
			}), getResourceI)
		assert.Nil(t, err)

		expectedResourceDetails := MockResourceDetails()
		expectedObj := (*expectedResourceDetails.ResourceDetails).(*unstructured.Unstructured)
		resultObj := (*result.ResourceDetails).(*unstructured.Unstructured)
		for key, value := range expectedObj.Object {
			assert.Equal(t, value, resultObj.Object[key])
		}
	})
}

func TestCreateResourceErrorFromCreate(t *testing.T) {
	mockGetResourceI := new(mock.Mock)
	mockGetResourceI.On("func1", "Pod", "validNamespace", DefaultNamespace).
		Return(&MockResourceInterface{ReturnedValue: MockUnstructured(), ReturnedError: errors.New("error")}, nil)
	getResourceI := func(resourceType string, namespace string, 
		emptyNamespace string) (dynamic.ResourceInterface, *models.ModelError) {
		args := mockGetResourceI.Called(resourceType, namespace, emptyNamespace)
		return args.Get(0).(dynamic.ResourceInterface), nil
	}
	t.Run("Test CreateResourceErrorFromCreate", func(t *testing.T) {
		res, err := CreateResource("Pod", "validNamespace", *MockResourceDetailsMap(map[string]interface{}{
			"key": "value",
			"namespace": "validNamespace",
		}), getResourceI)

		assert.NotNil(t, err)
		assert.EqualValues(t, 500, err.Code)
		assert.EqualValues(t, "Internal server error: error", err.Message)
		assert.EqualValues(t, res, models.ResourceDetails{})
	})
}

func TestDeleteResourceError(t *testing.T) {
	mockGetResourceI := new(mock.Mock)
	mockGetResourceI.On("func1", "Pod", "validNamespace", DefaultNamespace).Return(&MockResourceInterface{}, &models.ModelError{Code: 404, Message: "Not found"})
	getResourceI := func(resourceType string, namespace string, emptyNamespace string) (dynamic.ResourceInterface, *models.ModelError) {
		args := mockGetResourceI.Called(resourceType, namespace, emptyNamespace)
		return args.Get(0).(dynamic.ResourceInterface), args.Get(1).(*models.ModelError)
	}

	t.Run("Test DeleteResourceError", func(t *testing.T) {
		err := DeleteResource("Pod", "validNamespace", "validName", getResourceI)
		assert.NotNil(t, err)
		assert.Equal(t, &models.ModelError{Code: 404, Message: "Not found"}, err)
	})
}

func TestDeleteResourceSuccess(t *testing.T) {
	mockGetResourceI := new(mock.Mock)
	mockGetResourceI.On("func1", "Pod", "validNamespace", DefaultNamespace).Return(&MockResourceInterface{}, nil)
	getResourceI := func(resourceType string, namespace string, emptyNamespace string) (dynamic.ResourceInterface, *models.ModelError) {
		args := mockGetResourceI.Called(resourceType, namespace, emptyNamespace)
		return args.Get(0).(dynamic.ResourceInterface), nil
	}

	t.Run("Test DeleteResource", func(t *testing.T) {
		err := DeleteResource("Pod", "validNamespace", "validName", getResourceI)
		assert.Nil(t, err)
	})
}

func TestDeleteResourceErrorFromDelete(t *testing.T) {
	mockGetResourceI := new(mock.Mock)
	mockGetResourceI.On("func1", "Pod", "validNamespace", DefaultNamespace).Return(&MockResourceInterface{ReturnedError: errors.New("error")}, nil)
	getResourceI := func(resourceType string, namespace string, emptyNamespace string) (dynamic.ResourceInterface, *models.ModelError) {
		args := mockGetResourceI.Called(resourceType, namespace, emptyNamespace)
		return args.Get(0).(dynamic.ResourceInterface), nil
	}
	t.Run("Test DeleteResourceErrorFromDelete", func(t *testing.T) {
		err := DeleteResource("Pod", "validNamespace", "validName", getResourceI)

		assert.NotNil(t, err)
		assert.EqualValues(t, 500, err.Code)
		assert.EqualValues(t, "Internal server error: error", err.Message)
	})
}

func TestUpdateResourceError(t *testing.T) {
	mockGetResourceI := new(mock.Mock)
	mockGetResourceI.On("func1", "Pod", "validNamespace", DefaultNamespace).Return(&MockResourceInterface{}, &models.ModelError{Code: 404, Message: "Not found"})
	getResourceI := func(resourceType string, namespace string, emptyNamespace string) (dynamic.ResourceInterface, *models.ModelError) {
		args := mockGetResourceI.Called(resourceType, namespace, emptyNamespace)
		return args.Get(0).(dynamic.ResourceInterface), args.Get(1).(*models.ModelError)
	}

	t.Run("Test UpdateResourceError", func(t *testing.T) {
		_, err := UpdateResource("Pod", "validNamespace", "validName", *MockResourceDetails(), getResourceI)
		assert.NotNil(t, err)
		assert.Equal(t, &models.ModelError{Code: 404, Message: "Not found"}, err)
	})
}

func TestUpdateResourceSuccess(t *testing.T) {
	expected := map[string]interface{}{
			"key": "value",
			"namespace": "validNamespace",
			"metadata": map[string]interface{}{"name": "validName"},
		}

	mockGetResourceI := new(mock.Mock)
	mockGetResourceI.On("func1", "Pod", "validNamespace", DefaultNamespace).
		Return(&MockResourceInterface{ReturnedValue: &unstructured.Unstructured{Object: expected}}, nil)

	getResourceI := func(resourceType string, namespace string, emptyNamespace string) (dynamic.ResourceInterface, *models.ModelError) {
		args := mockGetResourceI.Called(resourceType, namespace, emptyNamespace)
		return args.Get(0).(dynamic.ResourceInterface), nil
	}

	t.Run("Test UpdateResource", func(t *testing.T) {
		result, err := UpdateResource("Pod", "validNamespace", "validName", *MockResourceDetailsMap(expected), getResourceI)
		assert.Nil(t, err)
		resultObj := (*result.ResourceDetails).(*unstructured.Unstructured)
		for key, value := range expected {
			assert.EqualValues(t, value, resultObj.Object[key])
		}
	})
}

func TestUpdateResourceErrorFromUpdate (t *testing.T) {
	dummy := map[string]interface{}{
			"key": "value",
			"namespace": "validNamespace",
			"metadata": map[string]interface{}{"name": "validName"},
		}

	mockGetResourceI := new(mock.Mock)
	mockGetResourceI.On("func1", "Pod", "validNamespace", DefaultNamespace).
		Return(&MockResourceInterface{ReturnedValue: MockUnstructured(), ReturnedError: errors.New("error")}, nil)
	getResourceI := func(resourceType string, namespace string, emptyNamespace string) (dynamic.ResourceInterface, *models.ModelError) {
		args := mockGetResourceI.Called(resourceType, namespace, emptyNamespace)
		return args.Get(0).(dynamic.ResourceInterface), nil
	}
	t.Run("Test UpdateResourceErrorFromUpdate", func(t *testing.T) {
		res, err := UpdateResource("Pod", "validNamespace", "validName", *MockResourceDetailsMap(dummy), getResourceI)

		assert.NotNil(t, err)
		assert.EqualValues(t, 500, err.Code)
		assert.EqualValues(t, "Internal server error: error", err.Message)
		assert.EqualValues(t, res, models.ResourceDetails{})
	})
}

func TestUpdateResourceWrongName (t *testing.T) {
	expected := map[string]interface{}{
			"key": "value",
			"namespace": "validNamespace",
			"metadata": map[string]interface{}{"name": "differentNameThanExpected"},
		}

	mockGetResourceI := new(mock.Mock)
	mockGetResourceI.On("func1", "Pod", "validNamespace", DefaultNamespace).
		Return(&MockResourceInterface{ReturnedValue: &unstructured.Unstructured{Object: expected}}, nil)

	getResourceI := func(resourceType string, namespace string, emptyNamespace string) (dynamic.ResourceInterface, *models.ModelError) {
		args := mockGetResourceI.Called(resourceType, namespace, emptyNamespace)
		return args.Get(0).(dynamic.ResourceInterface), nil
	}

	t.Run("Test UpdateResource", func(t *testing.T) {
		_, err := UpdateResource("Pod", "validNamespace", "validName", *MockResourceDetailsMap(expected), getResourceI)

		assert.NotNil(t, err)
		assert.EqualValues(t, 400, err.Code)
		assert.EqualValues(t, "Invalid Input: Different resource names", err.Message)
	})
}

func TestMain(t *testing.T) {
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

func TestJsonToken(t *testing.T) {
	t.Run("TestJsonToken", func(t *testing.T) {
		// Extract the "sub" field from the JSON token

		tokenStr := `eyJhbGciOiJSUzI1NiIsInR5cCIgOiAiSldUIiwia2lkIiA6ICI3SERLbTBsSHJLY18ybHc0eFo1S0NBR0JObndCTDJsOUlucFJ5VVU4ZHBjIn0.eyJleHAiOjE3MzAxMjM0NjgsImlhdCI6MTczMDEyMzE2OCwiYXV0aF90aW1lIjoxNzMwMTIzMTAwLCJqdGkiOiJkZjMwNjk5OC00NWQzLTRhNGQtOTE4ZS1hM2Q5YTIwMzc5MzgiLCJpc3MiOiJodHRwOi8vbG9jYWxob3N0OjgwODEvcmVhbG1zL2FjY2Vzcy1tYW5hZ2VyIiwiYXVkIjoiYWNjb3VudCIsInN1YiI6ImRkOTY3NDIxLWEwNGUtNGMzYS1hNzRjLTU3ZTQ4M2RhZDFhOCIsInR5cCI6IkJlYXJlciIsImF6cCI6ImFjY291bnQtY29uc29sZSIsInNpZCI6IjQ1ODZmMzU0LTQ1ZjYtNGUzNi1hODdiLWExYWE3ZjVjZDg3MyIsImFjciI6IjAiLCJyZXNvdXJjZV9hY2Nlc3MiOnsiYWNjb3VudC1jb25zb2xlIjp7InJvbGVzIjpbInBvZC1yZWFkZXIiXX0sImFjY291bnQiOnsicm9sZXMiOlsibWFuYWdlLWFjY291bnQiLCJtYW5hZ2UtYWNjb3VudC1saW5rcyJdfX0sInNjb3BlIjoib3BlbmlkIHByb2ZpbGUgZW1haWwiLCJlbWFpbF92ZXJpZmllZCI6ZmFsc2UsIm5hbWUiOiJNYXJlayBGaXVrIiwicHJlZmVycmVkX3VzZXJuYW1lIjoibWFyZWZlazFAZ21haWwuY29tIiwiZ2l2ZW5fbmFtZSI6Ik1hcmVrIiwiZmFtaWx5X25hbWUiOiJGaXVrIiwiZW1haWwiOiJtYXJlZmVrMUBnbWFpbC5jb20ifQ.zVe1FBnNkx7OlYveHZVG9vNJqwEJTtua5rDFekFJ9sNFAXK7e-xahcuEoOAy4_YTAjfGtgvQMHq2hy61_30Xe1cp6okmH0YnXZ-w4WXaxKdB7tHNcpduFiQSeCFBp4COImTEyuvOqv4PjLjLu5N0wkyfXClhoTIjvn932e_QEpeAjCeG5nDTePk3SqDbVYKo3cK0Ymzap7U4-H1OmM_YGPoYTGzC1Qri2rspPtfoaFP3Uv3jYUmGA1dl8_b90QDRalOq8AZxzrnTJbm1VfHH0tbEfUZqQV8ok_Wjf7PQ27M8dajkXcYDNneFoCVlaFwrfXJcJDdFfOvTS4ryRy1ZyA`
		claims := jwt.MapClaims{}
		_, _ = jwt.ParseWithClaims(tokenStr, &claims, nil)

		var roles []string
		if resourceAccess, ok := claims["resource_access"].(map[string]interface{}); ok {
			for _, resource := range resourceAccess {
				if resourceMap, ok := resource.(map[string]interface{}); ok {
					if resourceRoles, ok := resourceMap["roles"].([]interface{}); ok {
						for _, role := range resourceRoles {
							if roleStr, ok := role.(string); ok {
								roles = append(roles, roleStr)
							}
						}
					}
				}
			}
		} else {
			t.Fatalf("Failed to retrieve resource_access from claims")
		}

		// Expected roles for assertion
		expectedRoles := []string{"pod-reader", "manage-account", "manage-account-links"}

		assert.Equal(t, expectedRoles, roles)


	})
}








