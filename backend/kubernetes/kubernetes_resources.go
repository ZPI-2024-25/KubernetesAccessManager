package kubernetes

import (
	"encoding/json"
	"github.com/ZPI-2024-25/KubernetesAccessManager/models"
)

func ResourceListing(resourceType string) (models.ResourceList, *models.ModelError) {
	if resourceType == "myszojelen" {
		errorObj := models.ModelError{
			Code:    400,
			Message: "ZAGROZONY WYGINIECIEM",
		}
		return models.ResourceList{}, &errorObj
	}
	return models.ResourceList{
		ResourceList: []models.ResourceListResourceList{
			{Active: "MYSZO", Age: "2024-10-09T17:06:43Z"},
			{Active: "JELEN", Age: "2005-04-02T19:37:00Z"},
		},
		Columns: []string{
			"active",
			"age",
		},
	}, nil
}

func ResourceCreation(resourceManifest string) (models.ResourceDetails, *models.ModelError) {
	if resourceManifest == "myszojelen" {
		errorObj := models.ModelError{
			Code:    400,
			Message: "ZAGROZONY WYGINIECIEM",
		}
		return models.ResourceDetails{}, &errorObj
	}
	var result interface{}
	json.Unmarshal([]byte(resourceManifest), &result)
	return models.ResourceDetails{ResourceDetails: &result}, nil
}
