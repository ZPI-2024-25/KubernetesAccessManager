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
	resourceStr      = "resource"
	restartsStr      = "restarts"
	rolesStr         = "roles"
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
	"ReplicaSet":               {nameStr, namespaceStr, desiredStr, currentStr, readyStr, ageStr},
	"Pod":                      {nameStr, namespaceStr, containersStr, restartsStr, controlledByStr, nodeStr, qosStr, ageStr, statusStr},
	"Deployment":               {nameStr, namespaceStr, podsStr, replicasStr, ageStr, conditionsStr},
	"ConfigMap":                {nameStr, namespaceStr, keysStr, ageStr},
	"Secret":                   {nameStr, namespaceStr, labelsStr, keysStr, typeStr, ageStr},
	"Ingress":                  {nameStr, namespaceStr, loadbalancersStr, ageStr}, // loadbalancers untested
	"PersistentVolumeClaim":    {nameStr, namespaceStr, storageClassStr, sizeStr, ageStr, statusStr},
	"StatefulSet":              {nameStr, namespaceStr, podsStr, replicasStr, ageStr},
	"DaemonSet":                {nameStr, namespaceStr, podsStr, nodeSelectorStr, ageStr},
	"Job":                      {nameStr, namespaceStr, completionsStr, ageStr, conditionsStr}, // completions untested, conditions unsure
	"CronJob":                  {nameStr, namespaceStr, scheduleStr, suspendStr, activeStr, lastScheduleStr, ageStr},
	"Service":                  {nameStr, namespaceStr, typeStr, clusterIpStr, portsStr, externalIpStr, selectorStr, ageStr, statusStr}, // status unsure, external ip untested
	"ServiceAccount":           {nameStr, namespaceStr, ageStr},
	"Node":                     {nameStr, taintsStr, rolesStr, versionStr, ageStr, conditionsStr},
	"Namespace":                {nameStr, labelsStr, ageStr, statusStr},
	"CustomResourceDefinition": {resourceStr, groupStr, versionStr, scopeStr, ageStr},
	"PersistentVolume":         {nameStr, storageClassStr, capacityStr, claimStr, ageStr, statusStr},
	"StorageClass":             {nameStr, provisionerStr, reclaimPolicyStr, defaultStr, ageStr},
	"ClusterRole":              {nameStr, ageStr},
	"ClusterRoleBinding":       {nameStr, bindingsStr, ageStr},
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
