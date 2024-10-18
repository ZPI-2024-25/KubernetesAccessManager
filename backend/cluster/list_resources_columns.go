package cluster

const (
	activeStr        = "active"
	ageStr           = "age"
	bindingsStr      = "bindings"
	capacityStr      = "capacity"
	claimStr         = "claim"
	clusterIpStr     = "cluster_ip"
	completionsStr   = "completions"
	conditionsStr    = "conditions"
	containersStr    = "containers"
	controlledByStr  = "controlled_by"
	currentStr       = "current"
	defaultStr       = "default"
	desiredStr       = "desired"
	externalIpStr    = "external_ip"
	groupStr         = "group"
	keysStr          = "keys"
	labelsStr        = "labels"
	lastScheduleStr  = "last_schedule"
	loadbalancersStr = "loadbalancers"
	nameStr          = "name"
	namespaceStr     = "namespace"
	nodeStr          = "node"
	nodeSelectorStr  = "node_selector"
	podsStr          = "pods"
	portsStr         = "ports"
	provisionerStr   = "provisioner"
	qosStr           = "qos"
	readyStr         = "ready"
	reclaimPolicyStr = "reclaim_policy"
	replicasStr      = "replicas"
	resourcesStr     = "resources"
	restartsStr      = "restarts"
	rolesStr         = "roles"
	rulesStr         = "rules"
	scheduleStr      = "schedule"
	scopeStr         = "scope"
	selectorStr      = "selector"
	sizeStr          = "size"
	statusStr        = "status"
	storageClassStr  = "storage_class"
	suspendStr       = "suspend"
	taintsStr        = "taints"
	typeStr          = "type"
	versionStr       = "version"
)

var resourceListColumns = map[string][]string{
	"ReplicaSet":               {nameStr, namespaceStr, desiredStr, currentStr, readyStr, ageStr},                                        // done
	"Pod":                      {nameStr, namespaceStr, containersStr, restartsStr, controlledByStr, nodeStr, qosStr, ageStr, statusStr}, // done
	"Deployment":               {nameStr, namespaceStr, podsStr, replicasStr, ageStr, conditionsStr},                                     // done, conditions not implemented
	"ConfigMap":                {nameStr, namespaceStr, keysStr, ageStr},                                                                 // done
	"Secret":                   {nameStr, namespaceStr, labelsStr, keysStr, typeStr, ageStr},                                             // done
	"Ingress":                  {nameStr, namespaceStr, loadbalancersStr, rulesStr, ageStr},                                              // done, rules not implemented, loadbalancers untested
	"PersistentVolumeClaim":    {nameStr, namespaceStr, storageClassStr, sizeStr, podsStr, ageStr, statusStr},                            // done, pods not implemented
	"StatefulSet":              {nameStr, namespaceStr, podsStr, replicasStr, ageStr},                                                    // done
	"DaemonSet":                {nameStr, namespaceStr, podsStr, nodeSelectorStr, ageStr},                                                // done, node selector not implemented
	"Job":                      {nameStr, namespaceStr, completionsStr, ageStr, conditionsStr},                                           // done, completions untested, conditions unsure
	"CronJob":                  {nameStr, namespaceStr, scheduleStr, suspendStr, activeStr, lastScheduleStr, ageStr},                     // done
	"Service":                  {nameStr, namespaceStr, typeStr, clusterIpStr, portsStr, externalIpStr, selectorStr, ageStr, statusStr},  // done, status unsure, external ip untested
	"ServiceAccount":           {nameStr, namespaceStr, ageStr},                                                                          // done
	"Node":                     {nameStr, taintsStr, rolesStr, versionStr, ageStr, conditionsStr},                                        // done, conditions simplified
	"Namespace":                {nameStr, labelsStr, ageStr, statusStr},                                                                  // done
	"CustomResourceDefinition": {resourcesStr, groupStr, versionStr, scopeStr, ageStr},                                                   // done
	"PersistentVolume":         {nameStr, storageClassStr, capacityStr, claimStr, ageStr, statusStr},                                     //done
	"StorageClass":             {nameStr, provisionerStr, reclaimPolicyStr, defaultStr, ageStr},                                          // done
	"ClusterRole":              {nameStr, ageStr},                                                                                        // done
	"ClusterRoleBinding":       {nameStr, bindingsStr, ageStr},                                                                           // done
}

func GetResourceListColumns(resourceType string) []string {
	if columns, ok := resourceListColumns[resourceType]; ok {
		return columns
	}
	return []string{}
}

func transposeResourceListColumns(input map[string][]string) map[string][]string {
	result := make(map[string][]string)

	for resourceType, columns := range input {
		for _, column := range columns {
			result[column] = append(result[column], resourceType)
		}
	}

	return result
}
