package cluster

import (
	"context"
	"fmt"
	"github.com/ZPI-2024-25/KubernetesAccessManager/common"
	"github.com/ZPI-2024-25/KubernetesAccessManager/models"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"strconv"
	"strings"
	"time"
)

func ListResources(resourceType string, namespace string) (models.ResourceList, *models.ModelError) {
	gvr, httpErr := GetResourceGroupVersion(resourceType)
	if httpErr != nil {
		return models.ResourceList{}, httpErr
	}

	singleton, err := common.GetInstance()
	if err != nil {
		return models.ResourceList{}, &models.ModelError{Code: 500, Message: fmt.Sprintf("Failed to get client instance: %s", err)}
	}
	dynamicClient := singleton.GetClientSet()

	var resources *unstructured.UnstructuredList
	if namespace == "" {
		resources, err = dynamicClient.Resource(gvr).List(context.TODO(), metav1.ListOptions{})
	} else {
		resources, err = dynamicClient.Resource(gvr).Namespace(namespace).List(context.TODO(), metav1.ListOptions{})
	}

	if err != nil {
		if errors.IsNotFound(err) {
			return models.ResourceList{}, &models.ModelError{Code: 404, Message: fmt.Sprintf("Resource not found: %s", err)}
		} else {
			return models.ResourceList{}, &models.ModelError{Code: 500, Message: fmt.Sprintf("Error: %s", err)}
		}
	}

	var resourceList models.ResourceList

	resourceList.Columns = GetResourceListColumns(resourceType)

	for _, resource := range resources.Items {
		var resourceDetails models.ResourceListResourceList

		metadata := resource.Object["metadata"].(map[string]interface{})
		spec, specExists := resource.Object["spec"].(map[string]interface{})
		status, statusExists := resource.Object["status"].(map[string]interface{})

		//Age
		if creationTimestampStr, found := metadata["creationTimestamp"].(string); found {
			creationTimestamp, err := time.Parse(time.RFC3339, creationTimestampStr)
			if err != nil {
				fmt.Println("Failed to parse creationTimestamp:", err)
			} else {
				ageDuration := time.Since(creationTimestamp)

				if ageDuration.Hours() < 24 {
					resourceDetails.Age = fmt.Sprintf("%.f h", ageDuration.Hours())
				} else {
					resourceDetails.Age = fmt.Sprintf("%.f d", ageDuration.Hours()/24)
				}
			}
		}

		//ControlledBy
		ownerReferences, found := metadata["ownerReferences"].([]interface{})
		if found {
			for _, ownerReference := range ownerReferences {
				ownerReferenceMap := ownerReference.(map[string]interface{})
				if ownerReferenceMap["kind"].(string) == "ReplicaSet" {
					resourceDetails.ControlledBy = ownerReferenceMap["name"].(string)
				}
			}
		}

		//Secret: Labels
		if resourceType == "Secret" {
			if labels, found := metadata["labels"].(string); found {
				resourceDetails.Labels = labels
			}
		}

		if resourceType != "CustomResourceDefinition" {
			if name, found := metadata["name"].(string); found {
				resourceDetails.Name = name
			}
		}

		if namespace, found := metadata["namespace"].(string); found {
			resourceDetails.Namespace = namespace
		}

		// ConfigMap Keys
		if resourceType == "ConfigMap" || resourceType == "Secret" {
			data, dataExists := resource.Object["data"].(map[string]interface{})
			binaryData, binaryDataExists := resource.Object["binaryData"].(map[string]interface{})

			var keys []string

			if dataExists {
				for key := range data {
					keys = append(keys, key)
				}
			}

			if binaryDataExists {
				for key := range binaryData {
					keys = append(keys, key)
				}
			}

			resourceDetails.Keys = fmt.Sprintf("%v", keys)
		}

		// Secret: Type
		if resourceType == "Secret" {
			if secretType, found := resource.Object["type"].(string); found {
				resourceDetails.Type_ = secretType
			}
		}

		// Node: Roles
		if labels, labelsFound := metadata["labels"].(map[string]interface{}); labelsFound {
			var roles []string
			for labelKey := range labels {
				if strings.HasPrefix(labelKey, "node-role.kubernetes.io/") {
					role := strings.TrimPrefix(labelKey, "node-role.kubernetes.io/")
					if role == "" {
						role = "master"
					}
					roles = append(roles, role)
				}
			}
			resourceDetails.Roles = strings.Join(roles, ", ")
		}

		// Namespace: Labels
		if resourceType == "Namespace" {
			if labels, found := metadata["labels"].(map[string]interface{}); found {
				var labelPairs []string
				for key, value := range labels {
					labelPairs = append(labelPairs, fmt.Sprintf("%s=%s", key, value))
				}
				resourceDetails.Labels = strings.Join(labelPairs, ", ")
			}
		}

		//QOS
		qosClass := "BestEffort"

		if specExists {
			if resourceType == "ReplicaSet" {
				desired, found := spec["replicas"].(int64)
				if found {
					resourceDetails.Desired = strconv.FormatInt(desired, 10)
				}
			} else {
				if replicas, found := spec["replicas"].(int64); found {
					resourceDetails.Replicas = strconv.FormatInt(replicas, 10)
				}
			}

			node, found := spec["nodeName"].(string)
			if found {
				resourceDetails.Node = node
			}

			//QOS
			containers, found := spec["containers"].([]interface{})
			if found {
				for _, container := range containers {
					containerMap := container.(map[string]interface{})
					resources, resExists := containerMap["resources"].(map[string]interface{})
					if resExists {
						requests, reqExists := resources["requests"].(map[string]interface{})
						limits, limExists := resources["limits"].(map[string]interface{})

						if limExists && reqExists {
							cpuRequest, cpuReqExists := requests["cpu"]
							cpuLimit, cpuLimExists := limits["cpu"]
							memoryRequest, memReqExists := requests["memory"]
							memoryLimit, memLimExists := limits["memory"]

							if cpuReqExists && cpuLimExists && memReqExists && memLimExists &&
								cpuRequest == cpuLimit && memoryRequest == memoryLimit {
								qosClass = "Guaranteed"
							} else {
								qosClass = "Burstable"
							}
						} else if reqExists {
							qosClass = "Burstable"
						}
					}
				}

				resourceDetails.Qos = qosClass
			}

			storageClass, found := spec["storageClassName"].(string)
			if found {
				resourceDetails.StorageClass = storageClass
			}

			size, found, _ := unstructured.NestedString(spec, "resources", "requests", "storage")
			if found {
				resourceDetails.Size = size
			}

			capacity, found, _ := unstructured.NestedString(spec, "capacity", "storage")
			if found {
				resourceDetails.Capacity = capacity
			}

			claim, found, _ := unstructured.NestedString(spec, "claimRef", "name")
			if found {
				resourceDetails.Claim = claim
			}

			type_, found := spec["type"].(string)
			if found {
				resourceDetails.Type_ = type_
			}

			clusterIp, found := spec["clusterIP"].(string)
			if found {
				resourceDetails.ClusterIp = clusterIp
			}

			ports, found := spec["ports"].([]interface{})
			if found {
				var portStrings []string
				for _, port := range ports {
					portMap := port.(map[string]interface{})
					portStrings = append(portStrings, fmt.Sprintf("%d:%d/%s", int(portMap["port"].(int64)), int(portMap["targetPort"].(int64)), portMap["protocol"]))
				}
				resourceDetails.Ports = fmt.Sprintf("%v", portStrings)
			}

			if resourceType == "Service" {
				selector, found := spec["selector"].(map[string]interface{})
				if found {
					var selectorOutput []string
					for key, value := range selector {
						selectorOutput = append(selectorOutput, fmt.Sprintf("%s:%s", key, value))
					}
					resourceDetails.Selector = strings.Join(selectorOutput, ", ")
				}
			}

			if specExternalIPs, found := spec["externalIPs"].([]interface{}); found {
				var externalIPs []string
				for _, ip := range specExternalIPs {
					if ipStr, ok := ip.(string); ok {
						externalIPs = append(externalIPs, ipStr)
					}
				}
				resourceDetails.ExternalIp = strings.Join(externalIPs, ", ")
			}

			//Node: Taints
			if resourceType == "Node" {
				if taints, found := spec["taints"].([]interface{}); found {
					taintsCount := len(taints)
					resourceDetails.Taints = strconv.Itoa(taintsCount)
				} else {
					resourceDetails.Taints = "0"
				}
			}

			//CRD: Group
			if group, found := spec["group"].(string); found {
				resourceDetails.Group = group
			}

			//CRD: Version
			if versions, found := spec["versions"].([]interface{}); found {
				for _, version := range versions {
					versionMap, ok := version.(map[string]interface{})
					if !ok {
						continue
					}
					if storage, found := versionMap["storage"].(bool); found && storage {
						if versionName, found := versionMap["name"].(string); found {
							resourceDetails.Version = versionName
							break
						}
					}
				}
			}

			//CRD: Resource
			if names, found := spec["names"].(map[string]interface{}); found {
				if singular, found := names["singular"].(string); found {
					resourceDetails.Resource = cases.Title(language.English, cases.Compact).String(singular)
				}
			}

			//CRD: Scope
			if scope, found := spec["scope"].(string); found {
				resourceDetails.Scope = scope
			}
		}

		if statusExists {
			//Deployment, PersistentVolumeClaim: Pods
			if resourceType == "Deployment" || resourceType == "PersistentVolumeClaim" {
				replicas, found := status["replicas"].(int64)
				unavailableReplicas, found2 := status["unavailableReplicas"].(int64)
				if found && found2 {
					resourceDetails.Pods = fmt.Sprintf("%d/%d", replicas-unavailableReplicas, replicas)
				}
			}

			//StatefulSet: Pods
			if resourceType == "StatefulSet" {
				availableReplicas, found := status["availableReplicas"].(int64)
				replicas, found2 := status["replicas"].(int64)
				if found && found2 {
					resourceDetails.Pods = fmt.Sprintf("%d/%d", availableReplicas, replicas)
				}
			}

			//DaemonSet: Pods
			if resourceType == "DaemonSet" {
				ready, found := status["numberReady"].(int64)
				desired, found2 := status["desiredNumberScheduled"].(int64)
				if found && found2 {
					resourceDetails.Pods = fmt.Sprintf("%d/%d", ready, desired)
				}
			}

			//Replicas: Ready, Current
			if resourceType == "ReplicaSet" {
				if replicas, found := status["replicas"].(int64); found {
					resourceDetails.Current = strconv.FormatInt(replicas, 10)
				}

				if ready, found := status["readyReplicas"].(int64); found {
					resourceDetails.Ready = strconv.FormatInt(ready, 10)
				} else {
					resourceDetails.Ready = "0"
				}
			}

			// Pods: Containers, Restarts
			containerStatuses, found := status["containerStatuses"].([]interface{})
			if found {
				containers := len(containerStatuses)
				readyContainers := 0
				restarts := 0
				for _, containerStatus := range containerStatuses {
					containerStatusMap := containerStatus.(map[string]interface{})
					if containerStatusMap["ready"].(bool) {
						readyContainers++
					}
					restarts += int(containerStatusMap["restartCount"].(int64))
				}
				resourceDetails.Containers = fmt.Sprintf("%d/%d", readyContainers, containers)
				resourceDetails.Restarts = strconv.Itoa(restarts)
			}

			//Status
			phase, found := status["phase"].(string)
			if found {
				resourceDetails.Status = phase
			}

			//LoadBalancers (Ingress), Status (Service)
			if lb, lbFound := status["loadBalancer"].(map[string]interface{}); lbFound {
				if ingressList, ingressFound := lb["ingress"].([]interface{}); ingressFound {
					var loadBalancerAddresses []string
					for _, ingress := range ingressList {
						ingressMap := ingress.(map[string]interface{})
						if ip, ipFound := ingressMap["ip"].(string); ipFound {
							loadBalancerAddresses = append(loadBalancerAddresses, ip)
						}
						if hostname, hostnameFound := ingressMap["hostname"].(string); hostnameFound {
							loadBalancerAddresses = append(loadBalancerAddresses, hostname)
						}
					}
					resourceDetails.Loadbalancers = fmt.Sprintf("%v", loadBalancerAddresses)
				}

				if resourceType == "Service" {
					resourceDetails.Status = "Active"
				}
			} else {
				if resourceType == "Service" {
					resourceDetails.Status = "Pending"
				}
			}

			// Node: Version
			if nodeInfo, found := status["nodeInfo"].(map[string]interface{}); found {
				if kubeletVersion, versionFound := nodeInfo["kubeletVersion"].(string); versionFound {
					resourceDetails.Version = kubeletVersion
				}
			}

			// Node: Conditions
			if resourceType == "Node" {
				if conditions, found := status["conditions"].([]interface{}); found {
					nodeReady := "Unknown"
					for _, condition := range conditions {
						conditionMap, ok := condition.(map[string]interface{})
						if !ok {
							continue
						}
						conditionType, _ := conditionMap["type"].(string)
						conditionStatus, _ := conditionMap["status"].(string)
						if conditionType == "Ready" {
							if conditionStatus == "True" {
								nodeReady = "Ready"
							} else {
								nodeReady = "NotReady"
							}
							break
						}
					}
					resourceDetails.Conditions = nodeReady
				}
			}

		}

		resourceList.ResourceList = append(resourceList.ResourceList, resourceDetails)
	}

	//metadata, _, err := unstructured.NestedMap(resource.Object, "metadata")
	//if err != nil {
	//	return models.Resource{}, &models.ModelError{Code: 500, Message: fmt.Sprintf("Error: %s", err)}
	//}
	//
	//spec, _, err := unstructured.NestedMap(resource.Object, "spec")
	//if err != nil {
	//	return models.Resource{}, &models.ModelError{Code: 500, Message: fmt.Sprintf("Error: %s", err)}
	//}
	//
	//status, _, err := unstructured.NestedMap(resource.Object, "status")
	//if err != nil {
	//	return models.Resource{}, &models.ModelError{Code: 500, Message: fmt.Sprintf("Error: %s", err)}
	//}
	//
	//var metadataSwagger interface{} = metadata
	//var specSwagger interface{} = spec
	//var statusSwagger interface{} = status
	//
	//return models.Resource{
	//	ApiVersion: resource.GetAPIVersion(),
	//	Kind:       resource.GetKind(),
	//	Metadata:   &metadataSwagger,
	//	Spec:       &specSwagger,
	//	Status:     &statusSwagger,
	//}, nil

	return resourceList, nil
}
