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
	"sort"
	"strconv"
	"strings"
)

var (
	transposedResourceListColumns = transposeResourceListColumns(resourceListColumns)
)

const (
	emptyNamespace    = ""
	serviceString     = "Service"
	deploymentString  = "Deployment"
	statefulSetString = "StatefulSet"
	daemonSetString   = "DaemonSet"
	nodeString        = "Node"
	jobString         = "Job"
	secretString      = "Secret"
)

func ListResources(resourceType string, namespace string, getResourceInterface ResourceInterfaceGetter) (models.ResourceList, *models.ModelError) {
	resourceInterface, err := getResourceInterface(resourceType, namespace, emptyNamespace)
	if err != nil {
		return models.ResourceList{}, err
	}

	resources, listErr := resourceInterface.List(context.TODO(), metav1.ListOptions{})

	if listErr != nil {
		return models.ResourceList{}, handleKubernetesError(listErr)
	}

	var resourceList models.ResourceList

	resourceList.Columns = GetResourceListColumns(resourceType)
	resourceList.ResourceList = []models.ResourceListResourceList{}

	for _, resource := range resources.Items {
		var resourceDetailsTruncated models.ResourceListResourceList

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

		extractResource(resource, resourceType, &resourceDetailsTruncated)

		extractRestarts(resource, resourceType, &resourceDetailsTruncated)

		extractRoles(resource, resourceType, &resourceDetailsTruncated)

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

		resourceList.ResourceList = append(resourceList.ResourceList, resourceDetailsTruncated)
	}
	return resourceList, nil
}

func extractActive(resource unstructured.Unstructured, resourceType string, resourceDetailsTruncated *models.ResourceListResourceList) {
	if slices.Contains(transposedResourceListColumns["active"], resourceType) {
		status, statusExists := resource.Object["status"].(map[string]interface{})
		if statusExists {
			if activeJobs, found := status["active"].([]interface{}); found {
				activeCount := len(activeJobs)
				resourceDetailsTruncated.Active = strconv.Itoa(activeCount)
			} else {
				resourceDetailsTruncated.Active = "0"
			}
		}
	}
}

func extractAge(resource unstructured.Unstructured, resourceType string, resourceDetailsTruncated *models.ResourceListResourceList) {
	if slices.Contains(transposedResourceListColumns["age"], resourceType) {
		metadata := resource.Object["metadata"].(map[string]interface{})
		if creationTimestampStr, found := metadata["creationTimestamp"].(string); found {
			resourceDetailsTruncated.Age = creationTimestampStr
		}
	}
}

func extractBindings(resource unstructured.Unstructured, resourceType string, resourceDetailsTruncated *models.ResourceListResourceList) {
	if slices.Contains(transposedResourceListColumns["bindings"], resourceType) {
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
	}
}

func extractCapacity(resource unstructured.Unstructured, resourceType string, resourceDetailsTruncated *models.ResourceListResourceList) {
	if slices.Contains(transposedResourceListColumns["capacity"], resourceType) {
		spec, specExists := resource.Object["spec"].(map[string]interface{})
		if specExists {
			capacity, found, _ := unstructured.NestedString(spec, "capacity", "storage")
			if found {
				resourceDetailsTruncated.Capacity = capacity
			}
		}
	}
}

func extractClaim(resource unstructured.Unstructured, resourceType string, resourceDetailsTruncated *models.ResourceListResourceList) {
	if slices.Contains(transposedResourceListColumns["claim"], resourceType) {
		spec, specExists := resource.Object["spec"].(map[string]interface{})
		if specExists {
			claim, found, _ := unstructured.NestedString(spec, "claimRef", "name")
			if found {
				resourceDetailsTruncated.Claim = claim
			}
		}
	}
}

func extractClusterIp(resource unstructured.Unstructured, resourceType string, resourceDetailsTruncated *models.ResourceListResourceList) {
	if slices.Contains(transposedResourceListColumns["cluster_ip"], resourceType) {
		spec, specExists := resource.Object["spec"].(map[string]interface{})
		if specExists {
			if clusterIp, found := spec["clusterIP"].(string); found {
				resourceDetailsTruncated.ClusterIp = clusterIp
			}
		}
	}
}

func extractCompletions(resource unstructured.Unstructured, resourceType string, resourceDetailsTruncated *models.ResourceListResourceList) {
	if slices.Contains(transposedResourceListColumns["completions"], resourceType) {
		spec, specExists := resource.Object["spec"].(map[string]interface{})
		status, statusExists := resource.Object["status"].(map[string]interface{})

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
}

func extractConditions(resource unstructured.Unstructured, resourceType string, resourceDetailsTruncated *models.ResourceListResourceList) {
	if slices.Contains(transposedResourceListColumns["conditions"], resourceType) {
		status, statusExists := resource.Object["status"].(map[string]interface{})
		if statusExists {
			if resourceType == deploymentString || resourceType == nodeString {
				if conditions, found := status["conditions"].([]interface{}); found {
					var conditionsOutput []string
					for _, condition := range conditions {
						conditionMap, ok := condition.(map[string]interface{})
						if !ok {
							continue
						}
						conditionType, _ := conditionMap["type"].(string)
						conditionStatus, _ := conditionMap["status"].(string)

						if conditionStatus == "True" {
							conditionsOutput = append(conditionsOutput, conditionType)
						}
					}
					resourceDetailsTruncated.Conditions = strings.Join(conditionsOutput, ", ")
				}
			}

			if resourceType == jobString {
				if conditions, found := status["conditions"].([]interface{}); found {
					conditionMap, ok := conditions[0].(map[string]interface{})
					if !ok {
						resourceDetailsTruncated.Conditions = "Unknown"
						return
					}

					conditionType, _ := conditionMap["type"].(string)
					resourceDetailsTruncated.Conditions = fmt.Sprintf("%s", conditionType)
				}
			}
		}
	}
}

func extractContainers(resource unstructured.Unstructured, resourceType string, resourceDetailsTruncated *models.ResourceListResourceList) {
	if slices.Contains(transposedResourceListColumns["containers"], resourceType) {
		status, statusExists := resource.Object["status"].(map[string]interface{})
		if statusExists {
			if containerStatuses, found := status["containerStatuses"].([]interface{}); found {
				containers := len(containerStatuses)
				readyContainers := 0
				for _, containerStatus := range containerStatuses {
					containerStatusMap := containerStatus.(map[string]interface{})
					if containerStatusMap["ready"].(bool) {
						readyContainers++
					}
				}
				resourceDetailsTruncated.Containers = fmt.Sprintf("%d/%d", readyContainers, containers)
			}
		}
	}
}

func extractControlledBy(resource unstructured.Unstructured, resourceType string, resourceDetailsTruncated *models.ResourceListResourceList) {
	if slices.Contains(transposedResourceListColumns["controlled_by"], resourceType) {
		metadata := resource.Object["metadata"].(map[string]interface{})
		if ownerReferences, found := metadata["ownerReferences"].([]interface{}); found {
			var controlledBy []string
			for _, ownerReference := range ownerReferences {
				ownerReferenceMap := ownerReference.(map[string]interface{})
				owner := fmt.Sprintf("%s:%s", ownerReferenceMap["kind"].(string), ownerReferenceMap["name"].(string))
				controlledBy = append(controlledBy, owner)
			}
			resourceDetailsTruncated.ControlledBy = strings.Join(controlledBy, ", ")
		}
	}
}

func extractCurrent(resource unstructured.Unstructured, resourceType string, resourceDetailsTruncated *models.ResourceListResourceList) {
	if slices.Contains(transposedResourceListColumns["current"], resourceType) {
		status, statusExists := resource.Object["status"].(map[string]interface{})
		if statusExists {
			if replicas, found := status["availableReplicas"].(int64); found {
				resourceDetailsTruncated.Current = strconv.FormatInt(replicas, 10)
			} else {
				resourceDetailsTruncated.Current = "0"
			}
		} else {
			resourceDetailsTruncated.Current = "0"
		}
	}
}

func extractDefault(resource unstructured.Unstructured, resourceType string, resourceDetailsTruncated *models.ResourceListResourceList) {
	if slices.Contains(transposedResourceListColumns["default"], resourceType) {
		metadata := resource.Object["metadata"].(map[string]interface{})
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
}

func extractDesired(resource unstructured.Unstructured, resourceType string, resourceDetailsTruncated *models.ResourceListResourceList) {
	if slices.Contains(transposedResourceListColumns["desired"], resourceType) {
		spec, specExists := resource.Object["spec"].(map[string]interface{})
		if specExists {
			if replicas, found := spec["replicas"].(int64); found {
				resourceDetailsTruncated.Desired = strconv.FormatInt(replicas, 10)
			}
		}
	}
}

func extractExternalIp(resource unstructured.Unstructured, resourceType string, resourceDetailsTruncated *models.ResourceListResourceList) {
	if slices.Contains(transposedResourceListColumns["external_ip"], resourceType) {
		resourceDetailsTruncated.ExternalIp = "-"

		serviceType, found, err := unstructured.NestedString(resource.Object, "spec", "type")
		if err != nil || !found {
			return
		}

		switch serviceType {
		case "LoadBalancer":
			ingressList, found, err := unstructured.NestedSlice(resource.Object, "status", "loadBalancer", "ingress")
			if err != nil || !found || len(ingressList) == 0 {
				resourceDetailsTruncated.ExternalIp = "<pending>"
				return
			}

			var externalIPs []string
			for _, ingress := range ingressList {
				ingressMap, ok := ingress.(map[string]interface{})
				if !ok {
					continue
				}
				if ip, found, _ := unstructured.NestedString(ingressMap, "ip"); found {
					externalIPs = append(externalIPs, ip)
				}
			}

			if len(externalIPs) > 0 {
				resourceDetailsTruncated.ExternalIp = strings.Join(externalIPs, ",")
			} else {
				resourceDetailsTruncated.ExternalIp = "<pending>"
			}

		case "NodePort", "ClusterIP":
			externalIPs, found, err := unstructured.NestedStringSlice(resource.Object, "spec", "externalIPs")
			if err != nil || !found || len(externalIPs) == 0 {
				resourceDetailsTruncated.ExternalIp = "-"
			} else {
				resourceDetailsTruncated.ExternalIp = strings.Join(externalIPs, ",")
			}

		default:
			resourceDetailsTruncated.ExternalIp = "<unknown>"
		}
	}
}

func extractGroup(resource unstructured.Unstructured, resourceType string, resourceDetailsTruncated *models.ResourceListResourceList) {
	if slices.Contains(transposedResourceListColumns["group"], resourceType) {
		spec, specExists := resource.Object["spec"].(map[string]interface{})
		if specExists {
			if group, found := spec["group"].(string); found {
				resourceDetailsTruncated.Group = group
			}
		}
	}
}

func extractKeys(resource unstructured.Unstructured, resourceType string, resourceDetailsTruncated *models.ResourceListResourceList) {
	if slices.Contains(transposedResourceListColumns["keys"], resourceType) {
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

		sort.Strings(keys)

		resourceDetailsTruncated.Keys = strings.Join(keys, ", ")
	}
}

func extractLabels(resource unstructured.Unstructured, resourceType string, resourceDetailsTruncated *models.ResourceListResourceList) {
	if slices.Contains(transposedResourceListColumns["labels"], resourceType) {
		metadata := resource.Object["metadata"].(map[string]interface{})
		if labels, found := metadata["labels"].(map[string]interface{}); found {
			var labelPairs []string
			for key, value := range labels {
				labelPairs = append(labelPairs, fmt.Sprintf("%s=%s", key, value))
			}

			sort.Strings(labelPairs)
			resourceDetailsTruncated.Labels = strings.Join(labelPairs, ", ")
		}

	}
}

func extractLastSchedule(resource unstructured.Unstructured, resourceType string, resourceDetailsTruncated *models.ResourceListResourceList) {
	if slices.Contains(transposedResourceListColumns["last_schedule"], resourceType) {
		status, statusExists := resource.Object["status"].(map[string]interface{})
		if statusExists {
			if lastSchedule, found := status["lastScheduleTime"].(string); found {
				resourceDetailsTruncated.LastSchedule = lastSchedule
			}
		}
	}
}

func extractLoadbalancers(resource unstructured.Unstructured, resourceType string, resourceDetailsTruncated *models.ResourceListResourceList) {
	if slices.Contains(transposedResourceListColumns["loadbalancers"], resourceType) {
		status, statusExists := resource.Object["status"].(map[string]interface{})
		if statusExists {
			if lb, lbFound := status["loadBalancer"].(map[string]interface{}); lbFound {
				if ingressList, ingressFound := lb["ingress"].([]interface{}); ingressFound {
					var loadBalancerAddresses []string
					for _, ingress := range ingressList {
						ingressMap := ingress.(map[string]interface{})
						if ip, ipFound := ingressMap["ip"].(string); ipFound {
							loadBalancerAddresses = append(loadBalancerAddresses, ip)
						}
					}
					resourceDetailsTruncated.Loadbalancers = strings.Join(loadBalancerAddresses, ", ")
				}
			}
		}
	}
}

func extractName(resource unstructured.Unstructured, resourceType string, resourceDetailsTruncated *models.ResourceListResourceList) {
	if slices.Contains(transposedResourceListColumns["name"], resourceType) {
		metadata := resource.Object["metadata"].(map[string]interface{})
		if name, found := metadata["name"].(string); found {
			resourceDetailsTruncated.Name = name
		}
	}
}

func extractNamespace(resource unstructured.Unstructured, resourceType string, resourceDetailsTruncated *models.ResourceListResourceList) {
	if slices.Contains(transposedResourceListColumns["namespace"], resourceType) {
		metadata := resource.Object["metadata"].(map[string]interface{})
		if namespace, found := metadata["namespace"].(string); found {
			resourceDetailsTruncated.Namespace = namespace
		}
	}
}

func extractNode(resource unstructured.Unstructured, resourceType string, resourceDetailsTruncated *models.ResourceListResourceList) {
	if slices.Contains(transposedResourceListColumns["node"], resourceType) {
		spec, specExists := resource.Object["spec"].(map[string]interface{})
		if specExists {
			if node, found := spec["nodeName"].(string); found {
				resourceDetailsTruncated.Node = node
			}
		}
	}
}

func extractNodeSelector(resource unstructured.Unstructured, resourceType string, resourceDetailsTruncated *models.ResourceListResourceList) {
	if slices.Contains(transposedResourceListColumns["node_selector"], resourceType) {
		nodeSelector, found, err := unstructured.NestedStringMap(resource.Object, "spec", "template", "spec", "nodeSelector")
		if err != nil || !found || len(nodeSelector) == 0 {
			resourceDetailsTruncated.NodeSelector = "None"
			return
		}

		var selectorStrings []string
		for key, value := range nodeSelector {
			selectorStrings = append(selectorStrings, fmt.Sprintf("%s=%s", key, value))
		}

		sort.Strings(selectorStrings)
		resourceDetailsTruncated.NodeSelector = strings.Join(selectorStrings, ", ")
	}
}

func extractPods(resource unstructured.Unstructured, resourceType string, resourceDetailsTruncated *models.ResourceListResourceList) {
	if slices.Contains(transposedResourceListColumns["pods"], resourceType) {
		status, statusExists := resource.Object["status"].(map[string]interface{})
		if statusExists {
			if resourceType == deploymentString {
				replicas, found := status["replicas"].(int64)
				unavailableReplicas, found2 := status["unavailableReplicas"].(int64)
				if found && found2 {
					resourceDetailsTruncated.Pods = fmt.Sprintf("%d/%d", replicas-unavailableReplicas, replicas)
				} else if found {
					resourceDetailsTruncated.Pods = fmt.Sprintf("%d/%d", replicas, replicas)
				}
			}

			if resourceType == statefulSetString {
				availableReplicas, found := status["availableReplicas"].(int64)
				replicas, found2 := status["replicas"].(int64)
				if found && found2 {
					resourceDetailsTruncated.Pods = fmt.Sprintf("%d/%d", availableReplicas, replicas)
				} else if found2 {
					resourceDetailsTruncated.Pods = fmt.Sprintf("0/%d", replicas)
				}
			}

			if resourceType == daemonSetString {
				ready, found := status["numberReady"].(int64)
				desired, found2 := status["desiredNumberScheduled"].(int64)
				if found && found2 {
					resourceDetailsTruncated.Pods = fmt.Sprintf("%d/%d", ready, desired)
				}
			}
		}
	}
}

func extractPorts(resource unstructured.Unstructured, resourceType string, resourceDetailsTruncated *models.ResourceListResourceList) {
	if slices.Contains(transposedResourceListColumns["ports"], resourceType) {
		spec, specExists := resource.Object["spec"].(map[string]interface{})
		if specExists {
			if ports, found := spec["ports"].([]interface{}); found {
				var portStrings []string
				for _, port := range ports {
					portMap := port.(map[string]interface{})

					portNumber, found1 := portMap["port"].(int64)
					targetPort, found2 := portMap["targetPort"].(int64)
					nodePort, found3 := portMap["nodePort"].(int64)
					protocol, found4 := portMap["protocol"].(string)

					if found1 && found2 && found4 {
						portStrings = append(portStrings, fmt.Sprintf("%d:%d/%s", int(portNumber), int(targetPort), protocol))
					} else if found1 && found3 && found4 {
						portStrings = append(portStrings, fmt.Sprintf("%d:%d/%s", int(portNumber), int(nodePort), protocol))
					} else if found1 && found4 {
						portStrings = append(portStrings, fmt.Sprintf("%d/%s", int(portNumber), protocol))
					} else if found1 {
						portStrings = append(portStrings, fmt.Sprintf("%d", int(portNumber)))
					}
				}
				resourceDetailsTruncated.Ports = strings.Join(portStrings, ", ")
			}
		}
	}
}

func extractProvisioner(resource unstructured.Unstructured, resourceType string, resourceDetailsTruncated *models.ResourceListResourceList) {
	if slices.Contains(transposedResourceListColumns["provisioner"], resourceType) {
		if provisioner, found := resource.Object["provisioner"].(string); found {
			resourceDetailsTruncated.Provisioner = provisioner
		}
	}
}

func extractQos(resource unstructured.Unstructured, resourceType string, resourceDetailsTruncated *models.ResourceListResourceList) {
	if slices.Contains(transposedResourceListColumns["qos"], resourceType) {
		qosClass, found, err := unstructured.NestedString(resource.Object, "status", "qosClass")
		if err != nil || !found {
			resourceDetailsTruncated.Qos = "Unknown"
		} else {
			resourceDetailsTruncated.Qos = qosClass
		}
	}
}

func extractReady(resource unstructured.Unstructured, resourceType string, resourceDetailsTruncated *models.ResourceListResourceList) {
	if slices.Contains(transposedResourceListColumns["ready"], resourceType) {
		status, statusExists := resource.Object["status"].(map[string]interface{})
		if statusExists {
			if ready, found := status["readyReplicas"].(int64); found {
				resourceDetailsTruncated.Ready = strconv.FormatInt(ready, 10)
			} else {
				resourceDetailsTruncated.Ready = "0"
			}
		} else {
			resourceDetailsTruncated.Ready = "0"
		}
	}
}

func extractReclaimPolicy(resource unstructured.Unstructured, resourceType string, resourceDetailsTruncated *models.ResourceListResourceList) {
	if slices.Contains(transposedResourceListColumns["reclaim_policy"], resourceType) {
		if reclaimPolicy, found := resource.Object["reclaimPolicy"].(string); found {
			resourceDetailsTruncated.ReclaimPolicy = reclaimPolicy
		}
	}
}

func extractReplicas(resource unstructured.Unstructured, resourceType string, resourceDetailsTruncated *models.ResourceListResourceList) {
	if slices.Contains(transposedResourceListColumns["replicas"], resourceType) {
		spec, specExists := resource.Object["spec"].(map[string]interface{})
		if specExists {
			if replicas, found := spec["replicas"].(int64); found {
				resourceDetailsTruncated.Replicas = strconv.FormatInt(replicas, 10)
			} else {
				resourceDetailsTruncated.Replicas = "0"
			}
		} else {
			resourceDetailsTruncated.Replicas = "0"
		}
	}
}

func extractResource(resource unstructured.Unstructured, resourceType string, resourceDetailsTruncated *models.ResourceListResourceList) {
	if slices.Contains(transposedResourceListColumns["resource"], resourceType) {
		spec, specExists := resource.Object["spec"].(map[string]interface{})
		if specExists {
			if names, found := spec["names"].(map[string]interface{}); found {
				if singular, found := names["singular"].(string); found {
					resourceDetailsTruncated.Resource = cases.Title(language.English, cases.Compact).String(singular)
				}
			}
		}
	}
}

func extractRestarts(resource unstructured.Unstructured, resourceType string, resourceDetailsTruncated *models.ResourceListResourceList) {
	if slices.Contains(transposedResourceListColumns["restarts"], resourceType) {
		status, statusExists := resource.Object["status"].(map[string]interface{})
		if statusExists {
			if containerStatuses, found := status["containerStatuses"].([]interface{}); found {
				restarts := 0
				for _, containerStatus := range containerStatuses {
					containerStatusMap := containerStatus.(map[string]interface{})
					restarts += int(containerStatusMap["restartCount"].(int64))
				}
				resourceDetailsTruncated.Restarts = strconv.Itoa(restarts)
			}
		}
	}
}

func extractRoles(resource unstructured.Unstructured, resourceType string, resourceDetailsTruncated *models.ResourceListResourceList) {
	if slices.Contains(transposedResourceListColumns["roles"], resourceType) {
		metadata := resource.Object["metadata"].(map[string]interface{})
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

			sort.Strings(roles)
			resourceDetailsTruncated.Roles = strings.Join(roles, ", ")
		}
	}
}

func extractSchedule(resource unstructured.Unstructured, resourceType string, resourceDetailsTruncated *models.ResourceListResourceList) {
	if slices.Contains(transposedResourceListColumns["schedule"], resourceType) {
		spec, specExists := resource.Object["spec"].(map[string]interface{})
		if specExists {
			if schedule, found := spec["schedule"].(string); found {
				resourceDetailsTruncated.Schedule = schedule
			}
		}
	}
}

func extractScope(resource unstructured.Unstructured, resourceType string, resourceDetailsTruncated *models.ResourceListResourceList) {
	if slices.Contains(transposedResourceListColumns["scope"], resourceType) {
		spec, specExists := resource.Object["spec"].(map[string]interface{})
		if specExists {
			if scope, found := spec["scope"].(string); found {
				resourceDetailsTruncated.Scope = scope
			}
		}
	}
}

func extractSelector(resource unstructured.Unstructured, resourceType string, resourceDetailsTruncated *models.ResourceListResourceList) {
	if slices.Contains(transposedResourceListColumns["selector"], resourceType) {
		spec, specExists := resource.Object["spec"].(map[string]interface{})
		if specExists {
			if selector, found := spec["selector"].(map[string]interface{}); found {
				var selectorOutput []string
				for key, value := range selector {
					selectorOutput = append(selectorOutput, fmt.Sprintf("%s:%s", key, value))
				}

				sort.Strings(selectorOutput)
				resourceDetailsTruncated.Selector = strings.Join(selectorOutput, ", ")
			}
		}
	}
}

func extractSize(resource unstructured.Unstructured, resourceType string, resourceDetailsTruncated *models.ResourceListResourceList) {
	if slices.Contains(transposedResourceListColumns["size"], resourceType) {
		spec, specExists := resource.Object["spec"].(map[string]interface{})
		if specExists {
			size, found, _ := unstructured.NestedString(spec, "resources", "requests", "storage")
			if found {
				resourceDetailsTruncated.Size = size
			}
		}
	}
}

func extractStatus(resource unstructured.Unstructured, resourceType string, resourceDetailsTruncated *models.ResourceListResourceList) {
	if slices.Contains(transposedResourceListColumns["status"], resourceType) {
		status, statusExists := resource.Object["status"].(map[string]interface{})
		if statusExists {
			if resourceType == serviceString {
				serviceType, found, err := unstructured.NestedString(resource.Object, "spec", "type")
				if err != nil || !found {
					resourceDetailsTruncated.Status = "Unknown"
					return
				}

				if serviceType == "LoadBalancer" {
					ingresses, found, err := unstructured.NestedSlice(resource.Object, "status", "loadBalancer", "ingress")
					if err != nil {
						resourceDetailsTruncated.Status = "Unknown"
						return
					}

					if found && len(ingresses) > 0 {
						resourceDetailsTruncated.Status = "Active"
					} else {
						resourceDetailsTruncated.Status = "Pending"
					}
				} else {
					resourceDetailsTruncated.Status = "Active"
				}
			} else {
				if phase, found := status["phase"].(string); found {
					resourceDetailsTruncated.Status = phase
				}
			}
		}
	}
}

func extractStorageClass(resource unstructured.Unstructured, resourceType string, resourceDetailsTruncated *models.ResourceListResourceList) {
	if slices.Contains(transposedResourceListColumns["storage_class"], resourceType) {
		spec, specExists := resource.Object["spec"].(map[string]interface{})
		if specExists {
			if storageClass, found := spec["storageClassName"].(string); found {
				resourceDetailsTruncated.StorageClass = storageClass
			}
		}
	}
}

func extractSuspend(resource unstructured.Unstructured, resourceType string, resourceDetailsTruncated *models.ResourceListResourceList) {
	if slices.Contains(transposedResourceListColumns["suspend"], resourceType) {
		spec, specExists := resource.Object["spec"].(map[string]interface{})
		if specExists {
			if suspend, found := spec["suspend"].(bool); found {
				resourceDetailsTruncated.Suspend = strconv.FormatBool(suspend)
			}
		}
	}
}

func extractTaints(resource unstructured.Unstructured, resourceType string, resourceDetailsTruncated *models.ResourceListResourceList) {
	if slices.Contains(transposedResourceListColumns["taints"], resourceType) {
		spec, specExists := resource.Object["spec"].(map[string]interface{})
		if specExists {
			if taints, found := spec["taints"].([]interface{}); found {
				taintsCount := len(taints)
				resourceDetailsTruncated.Taints = strconv.Itoa(taintsCount)
			} else {
				resourceDetailsTruncated.Taints = "0"
			}
		}
	}
}

func extractType(resource unstructured.Unstructured, resourceType string, resourceDetailsTruncated *models.ResourceListResourceList) {
	if slices.Contains(transposedResourceListColumns["type"], resourceType) {
		spec, specExists := resource.Object["spec"].(map[string]interface{})
		if resourceType == secretString {
			if secretType, found := resource.Object["type"].(string); found {
				resourceDetailsTruncated.Type_ = secretType
			}
		} else {
			if specExists {
				if type_, found := spec["type"].(string); found {
					resourceDetailsTruncated.Type_ = type_
				}
			}
		}
	}
}

func extractVersion(resource unstructured.Unstructured, resourceType string, resourceDetailsTruncated *models.ResourceListResourceList) {
	if slices.Contains(transposedResourceListColumns["version"], resourceType) {
		spec, specExists := resource.Object["spec"].(map[string]interface{})
		status, statusExists := resource.Object["status"].(map[string]interface{})

		if resourceType == nodeString {
			if statusExists {
				if nodeInfo, found := status["nodeInfo"].(map[string]interface{}); found {
					if kubeletVersion, versionFound := nodeInfo["kubeletVersion"].(string); versionFound {
						resourceDetailsTruncated.Version = kubeletVersion
					}
				}
			}
		} else {
			if specExists {
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
			}
		}

	}
}
