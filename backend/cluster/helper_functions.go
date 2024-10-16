package cluster

import (
	"fmt"
	"github.com/ZPI-2024-25/KubernetesAccessManager/common"
	"github.com/ZPI-2024-25/KubernetesAccessManager/models"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/discovery"
)

func getAllowedResourceTypes() [20]string {
	return [20]string{
		"Pod",
		"Service",
		"Deployment",
		"ConfigMap",
		"StatefulSet",
		"DaemonSet",
		"Secret",
		"Ingress",
		"PersistentVolumeClaim",
		"ReplicaSet",
		"Node",
		"Namespace",
		"CustomResourceDefinition",
		"PersistentVolume",
		"Job",
		"CronJob",
		"ServiceAccount",
		"StorageClass",
		"ClusterRole",
		"ClusterRoleBinding",
	}
}

func isResourceTypeAllowed(resourceType string) bool {
	allowedResourceTypes := getAllowedResourceTypes()
	for _, allowedResourceType := range allowedResourceTypes {
		if allowedResourceType == resourceType {
			return true
		}
	}
	return false
}

func GetResourceGroupVersion(resourceType string) (output schema.GroupVersionResource, namespaced bool, error *models.ModelError) {
	if !isResourceTypeAllowed(resourceType) {
		return schema.GroupVersionResource{}, false, &models.ModelError{Code: 400, Message: fmt.Sprintf("Resource type '%s' not allowed", resourceType)}
	}

	dynamicClientSingleton, _ := common.GetInstance()
	config := dynamicClientSingleton.GetConfig()

	discoveryClient, _ := discovery.NewDiscoveryClientForConfig(config)

	apiResourceLists, _ := discoveryClient.ServerPreferredResources()

	for _, apiResourceList := range apiResourceLists {
		for _, apiResource := range apiResourceList.APIResources {
			if apiResource.Kind == resourceType {
				groupVersion, err := schema.ParseGroupVersion(apiResourceList.GroupVersion)
				if err != nil {
					return schema.GroupVersionResource{}, false, &models.ModelError{Code: 500, Message: fmt.Sprintf("%s", err.Error())}
				}

				return schema.GroupVersionResource{
					Group:    groupVersion.Group,
					Version:  groupVersion.Version,
					Resource: apiResource.Name,
				}, apiResource.Namespaced, nil
			}
		}
	}

	return schema.GroupVersionResource{}, false, &models.ModelError{Code: 400, Message: fmt.Sprintf("Resource type '%s' not found", resourceType)}
}

func GetResourceListColumns(resourceType string) []string {
	switch resourceType {
	case "ReplicaSet":
		return []string{"name", "namespace", "desired", "current", "ready", "age"} // done
	case "Pod":
		return []string{ // done
			"name",
			"namespace",
			"containers", // changed
			"restarts",   // possibly changed
			"controlled_by",
			"node",
			"qos",
			"age",
			"status"}
	case "Deployment":
		return []string{"name", "namespace", "pods", "replicas", "age", "conditions"} // done, conditions not implemented
	case "ConfigMap":
		return []string{"name", "namespace", "keys", "age"} // done
	case "Secret":
		return []string{"name", "namespace", "labels", "keys", "type", "age"} // done, labels untested
	case "Ingress":
		return []string{"name", "namespace", "loadBalancers", "rules", "age"} // done, rules not implemented, loadbalancers untested
	case "PersistentVolumeClaim":
		return []string{"name", "namespace", "storage_class", "size", "pods", "age", "status"} // done, pods not implemented
	case "StatefulSet":
		return []string{"name", "namespace", "pods", "replicas", "age"} // done
	case "DaemonSet":
		return []string{"name", "namespace", "pods", "node_selector", "age"} // done, node selector not implemented
	case "Job":
		return []string{"name", "namespace", "completions", "age", "conditions"} // done, completions untested, conditions unsure
	case "CronJob":
		return []string{"name", "namespace", "schedule", "suspend", "active", "last_schedule", "age"} // done
	case "Service":
		return []string{"name", "namespace", "type", "cluster_ip", "ports", "external_ip", "selector", "age", "status"} // done, status unsure, external ip untested
	case "ServiceAccount":
		return []string{"name", "namespace", "age"} // done
	case "Node":
		return []string{"name", "taints", "roles", "version", "age", "conditions"} // done, conditions simplified
	case "Namespace":
		return []string{"name", "labels", "age", "status"} // done
	case "CustomResourceDefinition":
		return []string{"resource", "group", "version", "scope", "age"} // done
	case "PersistentVolume":
		return []string{"name", "storage Class", "capacity", "claim", "age", "status"} //done
	case "StorageClass":
		return []string{"name", "provisioner", "reclaim Policy", "default", "age"} // done
	case "ClusterRole":
		return []string{"name", "age"} // done
	case "ClusterRoleBinding":
		return []string{"name", "bindings", "age"} // done
	default:
		return []string{}
	}
}
