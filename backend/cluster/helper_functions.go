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

func GetResourceGroupVersion(resourceType string) (schema.GroupVersionResource, *models.ModelError) {
	if !isResourceTypeAllowed(resourceType) {
		return schema.GroupVersionResource{}, &models.ModelError{Code: 400, Message: fmt.Sprintf("Resource type '%s' not allowed", resourceType)}
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
					return schema.GroupVersionResource{}, &models.ModelError{Code: 500, Message: fmt.Sprintf("%s", err.Error())}
				}

				return schema.GroupVersionResource{
					Group:    groupVersion.Group,
					Version:  groupVersion.Version,
					Resource: apiResource.Name,
				}, nil
			}
		}
	}

	return schema.GroupVersionResource{}, &models.ModelError{Code: 400, Message: fmt.Sprintf("Resource type '%s' not found", resourceType)}
}

func GetResourceListColumns(resourceType string) []string {
	switch resourceType {
	case "ReplicaSet":
		return []string{"Name", "Namespace", "Desired", "Current", "Ready", "Age"} // done
	case "Pod":
		return []string{ // done
			"Name",
			"Namespace",
			"Containers", // changed
			"Restarts",   // possibly changed
			"Controlled By",
			"Node",
			"QoS",
			"Age",
			"Status"}
	case "Deployment":
		return []string{"Name", "Namespace", "Pods", "Replicas", "Age", "Conditions"} // done, conditions not implemented
	case "ConfigMap":
		return []string{"Name", "Namespace", "Keys", "Age"} // done
	case "Secret":
		return []string{"Name", "Namespace", "Labels", "Keys", "Type", "Age"} // done, labels untested
	case "Ingress":
		return []string{"Name", "Namespace", "LoadBalancers", "Rules", "Age"} // done, rules not implemented, loadbalancers untested
	case "PersistentVolumeClaim":
		return []string{"Name", "Namespace", "Storage class", "Size", "Pods", "Age", "Status"} // done, pods not implemented
	case "StatefulSet":
		return []string{"Name", "Namespace", "Pods", "Replicas", "Age"} // done
	case "DaemonSet":
		return []string{"Name", "Namespace", "Pods", "Node Selector", "Age"} // done, node selector not implemented
	case "Job":
		return []string{"Name", "Namespace", "Completions", "Age", "Conditions"} // done, completions untested, conditions unsure
	case "CronJob":
		return []string{"Name", "Namespace", "Schedule", "Suspend", "Active", "Last schedule", "Age"} // done
	case "Service":
		return []string{"Name", "Namespace", "Type", "Cluster IP", "Ports", "External IP", "Selector", "Age", "Status"} // done, status unsure, external ip untested
	case "ServiceAccount":
		return []string{"Name", "Namespace", "Age"} // done
	case "Node":
		return []string{"Name", "Taints", "Roles", "Version", "Age", "Conditions"} // done, conditions simplified
	case "Namespace":
		return []string{"Name", "Labels", "Age", "Status"} // done
	case "CustomResourceDefinition":
		return []string{"Resource", "Group", "Version", "Scope", "Age"} // done
	case "PersistentVolume":
		return []string{"Name", "Storage Class", "Capacity", "Claim", "Age", "Status"} //done
	case "StorageClass":
		return []string{"Name", "Provisioner", "Reclaim Policy", "Default", "Age"} // done
	case "ClusterRole":
		return []string{"Name", "Age"} // done
	case "ClusterRoleBinding":
		return []string{"Name", "Bindings", "Age"} // done
	default:
		return []string{}
	}
}
