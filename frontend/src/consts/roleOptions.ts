export const resourcesOptions = [
    {value: "*", label: "All"},
    {value: "Pod", label: "Pod"},
    {value: "Deployment", label: "Deployment"},
    {value: "ConfigMap", label: "ConfigMap"},
    {value: "Secret", label: "Secret"},
    {value: "Ingress", label: "Ingress"},
    {value: "PersistentVolumeClaim", label: "PersistentVolumeClaim"},
    {value: "ReplicaSet", label: "ReplicaSet"},
    {value: "StatefulSet", label: "StatefulSet"},
    {value: "DaemonSet", label: "DaemonSet"},
    {value: "Job", label: "Job"},
    {value: "CronJob", label: "CronJob"},
    {value: "Service", label: "Service"},
    {value: "ServiceAccount", label: "ServiceAccount"},
    {value: "Node", label: "Node"},
    {value: "Namespace", label: "Namespace"},
    {value: "CustomResourceDefinition", label: "CustomResourceDefinition"},
    {value: "PersistentVolume", label: "PersistentVolume"},
    {value: "StorageClass", label: "StorageClass"},
    {value: "ClusterRole", label: "ClusterRole"},
    {value: "ClusterRoleBinding", label: "ClusterRoleBinding"},
    {value: "Helm", label: "Helm"},
];

export const operationsOptions = [
    {value: "*", label: "All"},
    {value: "create", label: "Create"},
    {value: "delete", label: "Delete"},
    {value: "read", label: "Read"},
    {value: "update", label: "Update"},
    {value: "list", label: "List"},
];