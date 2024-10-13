package kubernetes

import "github.com/ZPI-2024-25/KubernetesAccessManager/models"

func ResourceListing(resourceType string) (models.ResourceList, *models.ModelError) {
	return models.ResourceList{
		ResourceList: []models.ResourceListResourceList{
			{Active: "aaa"},
			{Active: "bbb"},
		},
		Columns: []string{
			"aaa",
			"bbb",
		},
	}, nil
}
