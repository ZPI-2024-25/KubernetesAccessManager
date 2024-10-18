package cluster

import (
	"context"
	"fmt"
	"github.com/ZPI-2024-25/KubernetesAccessManager/models"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/utils/strings/slices"
	"strconv"
	"strings"
)

var (
	transposedResourceListColumns map[string][]string = transposeResourceListColumns(resourceListColumns)
)

const (
	serviceStr = "Service"
)

func ListResources(resourceType string, namespace string) (models.ResourceList, *models.ModelError) {
	resourceInterface, err := getResourceInterface(resourceType, namespace)
	if err != nil {
		return models.ResourceList{}, err
	}

	resources, listErr := resourceInterface.List(context.TODO(), metav1.ListOptions{})

	if listErr != nil {
		handleKubernetesError(listErr)
	}

	var resourceList models.ResourceList

	resourceList.Columns = GetResourceListColumns(resourceType)

	for _, resource := range resources.Items {
		var resourceDetailsTruncated models.ResourceListResourceList

		metadata := resource.Object["metadata"].(map[string]interface{})
		spec, specExists := resource.Object["spec"].(map[string]interface{})
		status, statusExists := resource.Object["status"].(map[string]interface{})

		extractActive(resource, resourceType, &resourceDetailsTruncated)

		extractAge(resource, resourceType, &resourceDetailsTruncated)

		extractBindings(resource, resourceType, &resourceDetailsTruncated)

		extractCapacity(resource, resourceType, &resourceDetailsTruncated)

		extractClaim(resource, resourceType, &resourceDetailsTruncated)

		extractClusterIp(resource, resourceType, &resourceDetailsTruncated)

		extractCompletions(resource, resourceType, &resourceDetailsTruncated)

		extractConditions(resource, resourceType, &resourceDetailsTruncated)

		extractContainers(resource, resourceType, &resourceDetailsTruncated)

		extractControlledBy(resource, resourceType, &resourceDetailsTruncated)

		extractCurrent(resource, resourceType, &resourceDetailsTruncated)

		extractDefault(resource, resourceType, &resourceDetailsTruncated)

		extractDesired(resource, resourceType, &resourceDetailsTruncated)

		extractExternalIp(resource, resourceType, &resourceDetailsTruncated)

		extractGroup(resource, resourceType, &resourceDetailsTruncated)

		extractKeys(resource, resourceType, &resourceDetailsTruncated)

		extractLabels(resource, resourceType, &resourceDetailsTruncated)

		extractLastSchedule(resource, resourceType, &resourceDetailsTruncated)

		extractLoadbalancers(resource, resourceType, &resourceDetailsTruncated)

		extractName(resource, resourceType, &resourceDetailsTruncated)

		extractNamespace(resource, resourceType, &resourceDetailsTruncated)

		extractNode(resource, resourceType, &resourceDetailsTruncated)

		extractNodeSelector(resource, resourceType, &resourceDetailsTruncated)

		extractPods(resource, resourceType, &resourceDetailsTruncated)

		extractPorts(resource, resourceType, &resourceDetailsTruncated)

		extractProvisioner(resource, resourceType, &resourceDetailsTruncated)

		extractQos(resource, resourceType, &resourceDetailsTruncated)

		extractReady(resource, resourceType, &resourceDetailsTruncated)

		extractReclaimPolicy(resource, resourceType, &resourceDetailsTruncated)

		extractReplicas(resource, resourceType, &resourceDetailsTruncated)

		extractResources(resource, resourceType, &resourceDetailsTruncated)

		extractRestarts(resource, resourceType, &resourceDetailsTruncated)

		extractRoles(resource, resourceType, &resourceDetailsTruncated)

		extractRules(resource, resourceType, &resourceDetailsTruncated)

		extractSchedule(resource, resourceType, &resourceDetailsTruncated)

		extractScope(resource, resourceType, &resourceDetailsTruncated)

		extractSelector(resource, resourceType, &resourceDetailsTruncated)

		extractSize(resource, resourceType, &resourceDetailsTruncated)

		extractStatus(resource, resourceType, &resourceDetailsTruncated)

		extractStorageClass(resource, resourceType, &resourceDetailsTruncated)

		extractSuspend(resource, resourceType, &resourceDetailsTruncated)

		extractTaints(resource, resourceType, &resourceDetailsTruncated)

		extractType(resource, resourceType, &resourceDetailsTruncated)

		extractVersion(resource, resourceType, &resourceDetailsTruncated)

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

			resourceDetailsTruncated.Keys = fmt.Sprintf("%v", keys)
		}

		// Secret: Type
		if resourceType == "Secret" {
			if secretType, found := resource.Object["type"].(string); found {
				resourceDetailsTruncated.Type_ = secretType
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
			resourceDetailsTruncated.Roles = strings.Join(roles, ", ")
		}

		// Namespace, Secret: Labels
		if resourceType == "Namespace" || resourceType == "Secret" {
			if labels, found := metadata["labels"].(map[string]interface{}); found {
				var labelPairs []string
				for key, value := range labels {
					labelPairs = append(labelPairs, fmt.Sprintf("%s=%s", key, value))
				}
				resourceDetailsTruncated.Labels = strings.Join(labelPairs, ", ")
			}
		}

		if specExists {
			if resourceType == "ReplicaSet" {
				desired, found := spec["replicas"].(int64)
				if found {
					resourceDetailsTruncated.Desired = strconv.FormatInt(desired, 10)
				}
			} else {
				if replicas, found := spec["replicas"].(int64); found {
					resourceDetailsTruncated.Replicas = strconv.FormatInt(replicas, 10)
				}
			}

			storageClass, found := spec["storageClassName"].(string)
			if found {
				resourceDetailsTruncated.StorageClass = storageClass
			}

			size, found, _ := unstructured.NestedString(spec, "resources", "requests", "storage")
			if found {
				resourceDetailsTruncated.Size = size
			}

			capacity, found, _ := unstructured.NestedString(spec, "capacity", "storage")
			if found {
				resourceDetailsTruncated.Capacity = capacity
			}

			claim, found, _ := unstructured.NestedString(spec, "claimRef", "name")
			if found {
				resourceDetailsTruncated.Claim = claim
			}

			type_, found := spec["type"].(string)
			if found {
				resourceDetailsTruncated.Type_ = type_
			}

			clusterIp, found := spec["clusterIP"].(string)
			if found {
				resourceDetailsTruncated.ClusterIp = clusterIp
			}

			ports, found := spec["ports"].([]interface{})
			if found {
				var portStrings []string
				for _, port := range ports {
					portMap := port.(map[string]interface{})
					portStrings = append(portStrings, fmt.Sprintf("%d:%d/%s", int(portMap["port"].(int64)), int(portMap["targetPort"].(int64)), portMap["protocol"]))
				}
				resourceDetailsTruncated.Ports = fmt.Sprintf("%v", portStrings)
			}

			if resourceType == "Service" {
				selector, found := spec["selector"].(map[string]interface{})
				if found {
					var selectorOutput []string
					for key, value := range selector {
						selectorOutput = append(selectorOutput, fmt.Sprintf("%s:%s", key, value))
					}
					resourceDetailsTruncated.Selector = strings.Join(selectorOutput, ", ")
				}
			}

			if specExternalIPs, found := spec["externalIPs"].([]interface{}); found {
				var externalIPs []string
				for _, ip := range specExternalIPs {
					if ipStr, ok := ip.(string); ok {
						externalIPs = append(externalIPs, ipStr)
					}
				}
				resourceDetailsTruncated.ExternalIp = strings.Join(externalIPs, ", ")
			}

			//Node: Taints
			if resourceType == "Node" {
				if taints, found := spec["taints"].([]interface{}); found {
					taintsCount := len(taints)
					resourceDetailsTruncated.Taints = strconv.Itoa(taintsCount)
				} else {
					resourceDetailsTruncated.Taints = "0"
				}
			}

			//CRD: Group
			if group, found := spec["group"].(string); found {
				resourceDetailsTruncated.Group = group
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
							resourceDetailsTruncated.Version = versionName
							break
						}
					}
				}
			}

			//CRD: Resource
			if names, found := spec["names"].(map[string]interface{}); found {
				if singular, found := names["singular"].(string); found {
					resourceDetailsTruncated.Resource = cases.Title(language.English, cases.Compact).String(singular)
				}
			}

			//CRD: Scope
			if scope, found := spec["scope"].(string); found {
				resourceDetailsTruncated.Scope = scope
			}

			// CronJob: Schedule
			if schedule, found := spec["schedule"].(string); found {
				resourceDetailsTruncated.Schedule = schedule
			}

			// CronJob: Suspend
			if resourceType == "CronJob" {
				if suspend, found := spec["suspend"].(bool); found {
					resourceDetailsTruncated.Suspend = strconv.FormatBool(suspend)
				}
			}
		}

		if statusExists {
			//Deployment, PersistentVolumeClaim: Pods
			if resourceType == "Deployment" || resourceType == "PersistentVolumeClaim" {
				replicas, found := status["replicas"].(int64)
				unavailableReplicas, found2 := status["unavailableReplicas"].(int64)
				if found && found2 {
					resourceDetailsTruncated.Pods = fmt.Sprintf("%d/%d", replicas-unavailableReplicas, replicas)
				}
			}

			//StatefulSet: Pods
			if resourceType == "StatefulSet" {
				availableReplicas, found := status["availableReplicas"].(int64)
				replicas, found2 := status["replicas"].(int64)
				if found && found2 {
					resourceDetailsTruncated.Pods = fmt.Sprintf("%d/%d", availableReplicas, replicas)
				}
			}

			//DaemonSet: Pods
			if resourceType == "DaemonSet" {
				ready, found := status["numberReady"].(int64)
				desired, found2 := status["desiredNumberScheduled"].(int64)
				if found && found2 {
					resourceDetailsTruncated.Pods = fmt.Sprintf("%d/%d", ready, desired)
				}
			}

			//Replicas: Ready, Current
			if resourceType == "ReplicaSet" {
				if replicas, found := status["replicas"].(int64); found {
					resourceDetailsTruncated.Current = strconv.FormatInt(replicas, 10)
				}

				if ready, found := status["readyReplicas"].(int64); found {
					resourceDetailsTruncated.Ready = strconv.FormatInt(ready, 10)
				} else {
					resourceDetailsTruncated.Ready = "0"
				}
			}

			//LoadBalancers (Ingress)
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
					resourceDetailsTruncated.Loadbalancers = fmt.Sprintf("%v", loadBalancerAddresses)
				}
			}

			// Node: Version
			if nodeInfo, found := status["nodeInfo"].(map[string]interface{}); found {
				if kubeletVersion, versionFound := nodeInfo["kubeletVersion"].(string); versionFound {
					resourceDetailsTruncated.Version = kubeletVersion
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
					resourceDetailsTruncated.Conditions = nodeReady
				}
			}

			// CronJob: Last schedule
			if lastSchedule, found := status["lastScheduleTime"].(string); found {
				resourceDetailsTruncated.LastSchedule = lastSchedule
			}

			// CronJob: Active
			if resourceType == "CronJob" {
				if activeJobs, found := status["active"].([]interface{}); found {
					activeCount := len(activeJobs)
					resourceDetailsTruncated.Active = strconv.Itoa(activeCount)
				} else {
					resourceDetailsTruncated.Active = "0"
				}
			}

			// Job: Conditions
			if resourceType == "Job" {
				if conditions, found := status["conditions"].([]interface{}); found {
					conditionMap, ok := conditions[0].(map[string]interface{})
					if !ok {
						continue
					}

					conditionType, _ := conditionMap["type"].(string)
					resourceDetailsTruncated.Conditions = fmt.Sprintf("%s", conditionType)
				}
			}

		}

		if resourceType == "Job" {
			var completionsDesired int64 = 1
			var completionsSucceeded int64 = 0

			if specExists {
				if completions, found := spec["completions"].(int64); found {
					completionsDesired = completions
				}
			}

			if statusExists {
				if succeeded, found := status["succeeded"].(int64); found {
					completionsSucceeded = succeeded
				}
			}

			resourceDetailsTruncated.Completions = fmt.Sprintf("%d/%d", completionsSucceeded, completionsDesired)
		}

		// ClusterRoleBinding: Bindings
		subjects, subjectsExists := resource.Object["subjects"].([]interface{})
		if subjectsExists {
			var subjectsOutput []string
			for _, subject := range subjects {
				subjectMap := subject.(map[string]interface{})
				subjectName, nameExists := subjectMap["name"].(string)
				if nameExists {
					subjectsOutput = append(subjectsOutput, fmt.Sprintf("%s", subjectName))
				}
			}
			resourceDetailsTruncated.Bindings = strings.Join(subjectsOutput, ", ")
		}

		//StorageClass: Provisioner
		if provisioner, found := resource.Object["provisioner"].(string); found {
			resourceDetailsTruncated.Provisioner = provisioner
		}

		//StorageClass: Reclaim Policy
		if reclaimPolicy, found := resource.Object["reclaimPolicy"].(string); found {
			resourceDetailsTruncated.ReclaimPolicy = reclaimPolicy
		}

		//StorageClass: Default
		if resourceType == "StorageClass" {
			isDefault := "No"
			if annotations, annotationsExists := metadata["annotations"].(map[string]interface{}); annotationsExists {
				if value, found := annotations["storageclass.kubernetes.io/is-default-class"].(string); found && value == "true" {
					isDefault = "Yes"
				} else if value, found := annotations["storageclass.beta.kubernetes.io/is-default-class"].(string); found && value == "true" {
					isDefault = "Yes"
				}
			}
			resourceDetailsTruncated.Default_ = isDefault
		}

		resourceList.ResourceList = append(resourceList.ResourceList, resourceDetailsTruncated)
	}
	return resourceList, nil
}

func extractActive(resource unstructured.Unstructured, resourceType string, resourceDetailsTruncated *models.ResourceListResourceList) {

}

func extractAge(resource unstructured.Unstructured, resourceType string, resourceDetailsTruncated *models.ResourceListResourceList) {
	if slices.Contains(transposedResourceListColumns["age"], resourceType) {
		metadata := resource.Object["metadata"].(map[string]interface{})
		if creationTimestampStr, found := metadata["creationTimestamp"].(string); found {
			resourceDetailsTruncated.Age = creationTimestampStr
		} else {
			resourceDetailsTruncated.Age = ""
		}
	}
}

func extractBindings(resource unstructured.Unstructured, resourceType string, resourceDetailsTruncated *models.ResourceListResourceList) {

}

func extractCapacity(resource unstructured.Unstructured, resourceType string, resourceDetailsTruncated *models.ResourceListResourceList) {

}

func extractClaim(resource unstructured.Unstructured, resourceType string, resourceDetailsTruncated *models.ResourceListResourceList) {

}

func extractClusterIp(resource unstructured.Unstructured, resourceType string, resourceDetailsTruncated *models.ResourceListResourceList) {

}

func extractCompletions(resource unstructured.Unstructured, resourceType string, resourceDetailsTruncated *models.ResourceListResourceList) {

}

func extractConditions(resource unstructured.Unstructured, resourceType string, resourceDetailsTruncated *models.ResourceListResourceList) {

}

func extractContainers(resource unstructured.Unstructured, resourceType string, resourceDetailsTruncated *models.ResourceListResourceList) {
	if slices.Contains(transposedResourceListColumns["containers"], resourceType) {
		status, statusExists := resource.Object["status"].(map[string]interface{})
		if statusExists {
			containerStatuses, found := status["containerStatuses"].([]interface{})
			if found {
				containers := len(containerStatuses)
				readyContainers := 0
				for _, containerStatus := range containerStatuses {
					containerStatusMap := containerStatus.(map[string]interface{})
					if containerStatusMap["ready"].(bool) {
						readyContainers++
					}
				}
				resourceDetailsTruncated.Containers = fmt.Sprintf("%d/%d", readyContainers, containers)
			} else {
				resourceDetailsTruncated.Containers = ""
			}
		} else {
			resourceDetailsTruncated.Containers = ""
		}
	}
}

func extractControlledBy(resource unstructured.Unstructured, resourceType string, resourceDetailsTruncated *models.ResourceListResourceList) {
	if slices.Contains(transposedResourceListColumns["controlled_by"], resourceType) {
		metadata := resource.Object["metadata"].(map[string]interface{})
		ownerReferences, found := metadata["ownerReferences"].([]interface{})
		if found {
			var controlledBy []string
			for _, ownerReference := range ownerReferences {
				ownerReferenceMap := ownerReference.(map[string]interface{})
				owner := fmt.Sprintf("%s:%s", ownerReferenceMap["kind"].(string), ownerReferenceMap["name"].(string))
				controlledBy = append(controlledBy, owner)
			}
			resourceDetailsTruncated.ControlledBy = strings.Join(controlledBy, ", ")
		} else {
			resourceDetailsTruncated.ControlledBy = ""
		}
	}
}

func extractCurrent(resource unstructured.Unstructured, resourceType string, resourceDetailsTruncated *models.ResourceListResourceList) {

}

func extractDefault(resource unstructured.Unstructured, resourceType string, resourceDetailsTruncated *models.ResourceListResourceList) {

}

func extractDesired(resource unstructured.Unstructured, resourceType string, resourceDetailsTruncated *models.ResourceListResourceList) {

}

func extractExternalIp(resource unstructured.Unstructured, resourceType string, resourceDetailsTruncated *models.ResourceListResourceList) {

}

func extractGroup(resource unstructured.Unstructured, resourceType string, resourceDetailsTruncated *models.ResourceListResourceList) {

}

func extractKeys(resource unstructured.Unstructured, resourceType string, resourceDetailsTruncated *models.ResourceListResourceList) {

}

func extractLabels(resource unstructured.Unstructured, resourceType string, resourceDetailsTruncated *models.ResourceListResourceList) {

}

func extractLastSchedule(resource unstructured.Unstructured, resourceType string, resourceDetailsTruncated *models.ResourceListResourceList) {

}

func extractLoadbalancers(resource unstructured.Unstructured, resourceType string, resourceDetailsTruncated *models.ResourceListResourceList) {

}

func extractName(resource unstructured.Unstructured, resourceType string, resourceDetailsTruncated *models.ResourceListResourceList) {
	if slices.Contains(transposedResourceListColumns["name"], resourceType) {
		metadata := resource.Object["metadata"].(map[string]interface{})
		name, found := metadata["name"].(string)
		if found {
			resourceDetailsTruncated.Name = name
		} else {
			resourceDetailsTruncated.Name = ""
		}

	}
}

func extractNamespace(resource unstructured.Unstructured, resourceType string, resourceDetailsTruncated *models.ResourceListResourceList) {
	if slices.Contains(transposedResourceListColumns["namespace"], resourceType) {
		metadata := resource.Object["metadata"].(map[string]interface{})
		namespace, found := metadata["namespace"].(string)
		if found {
			resourceDetailsTruncated.Namespace = namespace
		} else {
			resourceDetailsTruncated.Namespace = ""
		}
	}
}

func extractNode(resource unstructured.Unstructured, resourceType string, resourceDetailsTruncated *models.ResourceListResourceList) {
	if slices.Contains(transposedResourceListColumns["node"], resourceType) {
		spec, specExists := resource.Object["spec"].(map[string]interface{})
		if specExists {
			node, found := spec["nodeName"].(string)
			if found {
				resourceDetailsTruncated.Node = node
			} else {
				resourceDetailsTruncated.Node = ""
			}
		}
	}
}

func extractNodeSelector(resource unstructured.Unstructured, resourceType string, resourceDetailsTruncated *models.ResourceListResourceList) {

}

func extractPods(resource unstructured.Unstructured, resourceType string, resourceDetailsTruncated *models.ResourceListResourceList) {

}

func extractPorts(resource unstructured.Unstructured, resourceType string, resourceDetailsTruncated *models.ResourceListResourceList) {

}

func extractProvisioner(resource unstructured.Unstructured, resourceType string, resourceDetailsTruncated *models.ResourceListResourceList) {

}

func extractQos(resource unstructured.Unstructured, resourceType string, resourceDetailsTruncated *models.ResourceListResourceList) {
	if slices.Contains(transposedResourceListColumns["qos"], resourceType) {
		spec, specExists := resource.Object["spec"].(map[string]interface{})
		qosClass := "BestEffort"

		if specExists {
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
			}
		}

		resourceDetailsTruncated.Qos = qosClass
	}
}

func extractReady(resource unstructured.Unstructured, resourceType string, resourceDetailsTruncated *models.ResourceListResourceList) {

}

func extractReclaimPolicy(resource unstructured.Unstructured, resourceType string, resourceDetailsTruncated *models.ResourceListResourceList) {

}

func extractReplicas(resource unstructured.Unstructured, resourceType string, resourceDetailsTruncated *models.ResourceListResourceList) {

}

func extractResources(resource unstructured.Unstructured, resourceType string, resourceDetailsTruncated *models.ResourceListResourceList) {

}

func extractRestarts(resource unstructured.Unstructured, resourceType string, resourceDetailsTruncated *models.ResourceListResourceList) {
	if slices.Contains(transposedResourceListColumns["restarts"], resourceType) {
		status, statusExists := resource.Object["status"].(map[string]interface{})
		if statusExists {
			containerStatuses, found := status["containerStatuses"].([]interface{})
			if found {
				restarts := 0
				for _, containerStatus := range containerStatuses {
					containerStatusMap := containerStatus.(map[string]interface{})
					restarts += int(containerStatusMap["restartCount"].(int64))
				}
				resourceDetailsTruncated.Restarts = strconv.Itoa(restarts)
			} else {
				resourceDetailsTruncated.Restarts = ""
			}
		} else {
			resourceDetailsTruncated.Restarts = ""
		}
	}
}

func extractRoles(resource unstructured.Unstructured, resourceType string, resourceDetailsTruncated *models.ResourceListResourceList) {

}

func extractRules(resource unstructured.Unstructured, resourceType string, resourceDetailsTruncated *models.ResourceListResourceList) {

}

func extractSchedule(resource unstructured.Unstructured, resourceType string, resourceDetailsTruncated *models.ResourceListResourceList) {

}

func extractScope(resource unstructured.Unstructured, resourceType string, resourceDetailsTruncated *models.ResourceListResourceList) {

}

func extractSelector(resource unstructured.Unstructured, resourceType string, resourceDetailsTruncated *models.ResourceListResourceList) {

}

func extractSize(resource unstructured.Unstructured, resourceType string, resourceDetailsTruncated *models.ResourceListResourceList) {

}

func extractStatus(resource unstructured.Unstructured, resourceType string, resourceDetailsTruncated *models.ResourceListResourceList) {
	if slices.Contains(transposedResourceListColumns["status"], resourceType) {
		status, statusExists := resource.Object["status"].(map[string]interface{})
		if statusExists {
			if resourceType == serviceStr {
				_, found := status["loadBalancer"].(map[string]interface{})
				if found {
					resourceDetailsTruncated.Status = "Active"
				} else {
					resourceDetailsTruncated.Status = "Pending"
				}
			} else {
				phase, found := status["phase"].(string)
				if found {
					resourceDetailsTruncated.Status = phase
				} else {
					resourceDetailsTruncated.Status = ""
				}
			}
		} else {
			resourceDetailsTruncated.Status = ""
		}
	}
}

func extractStorageClass(resource unstructured.Unstructured, resourceType string, resourceDetailsTruncated *models.ResourceListResourceList) {

}

func extractSuspend(resource unstructured.Unstructured, resourceType string, resourceDetailsTruncated *models.ResourceListResourceList) {

}

func extractTaints(resource unstructured.Unstructured, resourceType string, resourceDetailsTruncated *models.ResourceListResourceList) {

}

func extractType(resource unstructured.Unstructured, resourceType string, resourceDetailsTruncated *models.ResourceListResourceList) {

}

func extractVersion(resource unstructured.Unstructured, resourceType string, resourceDetailsTruncated *models.ResourceListResourceList) {

}
