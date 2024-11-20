export const getExampleResourceDefinition = (resourceType: string) => {
    const resourceDefinition = (() => {
        switch (resourceType) {
            case "Pod":
                return POD_YAML;
            case "Deployment":
                return DEPLOYMENT_YAML;
            case "Service":
                return SERVICE_YAML;
            case "ConfigMap":
                return CONFIG_MAP_YAML;
            case "Secret":
                return SECRET_YAML;
            case "Ingress":
                return INGRESS_YAML;
            case "PersistentVolumeClaim":
                return PERSISTENT_VOLUME_CLAIM_YAML;
            case "PersistentVolume":
                return PERSISTENT_VOLUME_YAML;
            case "Namespace":
                return NAMESPACE_YAML;
            case "ServiceAccount":
                return SERVICE_ACCOUNT_YAML;
            case "StatefulSet":
                return STATEFUL_SET_YAML;
            case "ReplicaSet":
                return REPLICA_SET_YAML;
            case "DaemonSet":
                return DAEMON_SET_YAML;
            case "Job":
                return JOB_YAML;
            case "CronJob":
                return CRON_JOB_YAML;
            case "CustomResourceDefinition":
                return CUSTOM_RESOURCE_DEFINITION_YAML;
            case "StorageClass":
                return STORAGE_CLASS_YAML;
            case "ClusterRole":
                return CLUSTER_ROLE_YAML;
            case "ClusterRoleBinding":
                return CLUSTER_ROLE_BINDING_YAML;
            default:
                return "";
        }
    })();

    return resourceDefinition.trim();
};

export const POD_YAML = `
apiVersion: v1
kind: Pod
metadata:
  name: title
  labels:
    role: title
spec:
  containers:
    - name: title
      image: nginx
      imagePullPolicy: IfNotPresent
      ports:
        - name: title
          containerPort: 80
          protocol: TCP
  restartPolicy: Always
`;

export const DEPLOYMENT_YAML = `
apiVersion: apps/v1
kind: Deployment
metadata:
  name: title
  labels:
    app: title
spec:
  replicas: 1
  selector:
    matchLabels:
      app: title
  template:
    metadata:
      name: title
      labels:
        app: title
    spec:
      containers:
        - name: title
          image: nginx
          imagePullPolicy: IfNotPresent
          ports:
            - containerPort: 80
              protocol: TCP
      restartPolicy: Always
`;

export const SERVICE_YAML = `
apiVersion: v1
kind: Service
metadata:
  name: title
spec:
  selector:
    app: title
  ports:
    - protocol: TCP
      port: 80
      targetPort: 8080
  type: NodePort
`;

export const CONFIG_MAP_YAML = `
apiVersion: v1
kind: ConfigMap
metadata:
  name: title
data:
  # kv format
  key: "value"
  $actual-config-file$: "name.properties"

  # file format
  name.properties: |
    field=value1,value2
`;

export const SECRET_YAML = `
apiVersion: v1
kind: Secret
metadata:
  name: title-basic-auth
type: kubernetes.io/basic-auth
stringData:
  username: developer
  password: password
`;

export const INGRESS_YAML = `
apiVersion: networking.k8s.io/v1
kind: Ingress
metadata:
  name: title
  annotations:
    nginx.ingress.kubernetes.io/rewrite-target: /
spec:
  ingressClassName: nginx
  rules:
    - http:
        paths:
          - path: /api
            pathType: Prefix
            backend:
              service:
                name: title
                port:
                  number: 80
`;

export const PERSISTENT_VOLUME_CLAIM_YAML = `
apiVersion: v1
kind: PersistentVolumeClaim
metadata:
  name: title
spec:
  accessModes:
    - ReadWriteOnce
  resources:
    requests:
      storage: 1Gi
`;

export const PERSISTENT_VOLUME_YAML = `
apiVersion: v1
kind: PersistentVolume
metadata:
  name: title
spec:
  capacity:
    storage: 5Gi
  accessModes:
    - ReadWriteOnce
  persistentVolumeReclaimPolicy: Retain
  hostPath:
    path: /data/title-pv
`;

export const NAMESPACE_YAML = `
apiVersion: v1
kind: Namespace
metadata:
  name: title
`;

export const SERVICE_ACCOUNT_YAML = `
apiVersion: v1
kind: ServiceAccount
metadata:
  name: title
`;

export const STATEFUL_SET_YAML = `
apiVersion: apps/v1
kind: StatefulSet
metadata:
  name: title
spec:
  serviceName: "title-service"
  replicas: 3
  selector:
    matchLabels:
      app: title
  template:
    metadata:
      labels:
        app: title
    spec:
      containers:
        - name: title
          image: nginx:1.14.2
          ports:
            - containerPort: 80
          volumeMounts:
            - name: data
              mountPath: /var/lib/data
  volumeClaimTemplates:
    - metadata:
        name: data
      spec:
        accessModes: ["ReadWriteOnce"]
        resources:
          requests:
            storage: 1Gi
`;

export const REPLICA_SET_YAML = `
apiVersion: apps/v1
kind: ReplicaSet
metadata:
  name: title
  labels:
    app: title
    tier: title
spec:
  replicas: 3
  selector:
    matchLabels:
      tier: title
  template:
    metadata:
      labels:
        tier: title
    spec:
      containers:
      - name: title
        image: nginx
`;

export const DAEMON_SET_YAML = `
apiVersion: apps/v1
kind: DaemonSet
metadata:
  name: title
spec:
  selector:
    matchLabels:
      app: title
  template:
    metadata:
      labels:
        app: title
    spec:
      containers:
        - name: busybox
          image: busybox
          args:
            - /bin/sh
            - -c
            - 'while true; do ping -c 4 8.8.8.8; sleep 60; done'
`;

export const JOB_YAML = `
apiVersion: batch/v1
kind: Job
metadata:
  name: title
spec:
  template:
    spec:
      containers:
        - name: title
          image: python:latest
          command: [ "python", "-c" ]
          args: [ "print('Hello from the Kubernetes job')" ]
      restartPolicy: Never
  backoffLimit: 4
`;

export const CRON_JOB_YAML = `
apiVersion: batch/v1
kind: CronJob
metadata:
  name: title
spec:
  schedule: "* * * * *" #	Run every minute
  jobTemplate:
    spec:
      template:
        spec:
          containers:
            - name: title
              image: busybox:latest
              imagePullPolicy: IfNotPresent
              command:
                - /bin/sh
                - -c
                - date; echo Hello!
          restartPolicy: OnFailure
`;

export const CUSTOM_RESOURCE_DEFINITION_YAML = `
apiVersion: apiextensions.k8s.io/v1
kind: CustomResourceDefinition
metadata:
  name: shirts.stable.example.com
spec:
  group: stable.example.com
  scope: Namespaced
  names:
    plural: shirts
    singular: shirt
    kind: Shirt
  versions:
  - name: v1
    served: true
    storage: true
    schema:
      openAPIV3Schema:
        type: object
        properties:
          spec:
            type: object
            properties:
              color:
                type: string
              size:
                type: string
    selectableFields:
    - jsonPath: .spec.color
    - jsonPath: .spec.size
    additionalPrinterColumns:
    - jsonPath: .spec.color
      name: Color
      type: string
    - jsonPath: .spec.size
      name: Size
      type: string
`;

export const STORAGE_CLASS_YAML = `
apiVersion: storage.k8s.io/v1
kind: StorageClass
metadata:
  name: title
  annotations:
    storageclass.kubernetes.io/is-default-class: "false"
provisioner: csi-driver.example-vendor.example
reclaimPolicy: Retain
allowVolumeExpansion: true
mountOptions:
  - discard
volumeBindingMode: WaitForFirstConsumer
parameters:
  guaranteedReadWriteLatency: "true"
`;

export const CLUSTER_ROLE_YAML = `
apiVersion: rbac.authorization.k8s.io/v1
kind: ClusterRole
metadata:
  name: secret-reader
rules:
  - apiGroups: [ "" ]
    resources: [ "secrets" ]
    verbs: [ "get", "watch", "list" ]
`;

export const CLUSTER_ROLE_BINDING_YAML = `
apiVersion: rbac.authorization.k8s.io/v1
kind: ClusterRoleBinding
metadata:
  name: read-secrets-global
subjects:
- kind: Group
  name: manager
  apiGroup: rbac.authorization.k8s.io
roleRef:
  kind: ClusterRole
  name: secret-reader
  apiGroup: rbac.authorization.k8s.io
`;