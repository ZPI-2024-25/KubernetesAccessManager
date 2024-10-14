package cluster

import (
	"fmt"
	"github.com/ZPI-2024-25/KubernetesAccessManager/common"
	"github.com/ZPI-2024-25/KubernetesAccessManager/models"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/discovery"
)

func getAllowedResourceTypes() [16]string {
	return [16]string{
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
		"CronJob",
		"ServiceAccount",
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
			"CPU",        // don't know what to put here
			"Memory",     // don't know what to put here
			"Restarts",   // possibly changed
			"Controlled By",
			"Node",
			"QoS",
			"Age",
			"Status"}
	case "Deployment":
		return []string{"Name", "Namespace", "Pods", "Replicas", "Age", "Conditions"} // done except conditions
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
	case "Jobs":
		return []string{"Name", "Namespace", "Completions", "Age", "Conditions"}
	case "CronJob":
		return []string{"Name", "Namespace", "Schedule", "Suspend", "Active", "Last schedule", "Age"}
	case "Service":
		return []string{"Name", "Namespace", "Type", "Cluster IP", "Ports", "External IP", "Selector", "Age", "Status"} // done, status unsure, external ip untested
	case "ServiceAccount":
		return []string{"Name", "Namespace", "Age"}
	case "Node":
		return []string{"Name", "CPU", "Memory", "Disk", "Taints", "Roles", "Version", "Age", "Conditions"} // done, (cpu, memory, disk) not implemented, conditions simplified
	case "Namespace":
		return []string{"Name", "Labels", "Age", "Status"} // done
	case "CustomResourceDefinition":
		return []string{"Resource", "Group", "Version", "Scope", "Age"} // done
	case "PersistentVolume":
		return []string{"Name", "Storage Class", "Capacity", "Claim", "Age", "Status"} //done
	case "StorageClass":
		return []string{"Name", "Provisioner", "Reclaim Policy", "Default", "Age"}
	case "ClusterRole":
		return []string{"Name", "Age"}
	case "ClusterRoleBinding":
		return []string{"Name", "Bindings", "Age"}
	default:
		return []string{}
	}
}
