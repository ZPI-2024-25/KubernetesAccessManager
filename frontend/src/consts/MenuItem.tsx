import React from "react";
import { FileOutlined, PieChartOutlined, TeamOutlined, UserOutlined } from "@ant-design/icons";
import { MenuItem } from "../types";

function getItem(
    label: React.ReactNode,
    key: React.Key,
    resourcelabel:string,
    icon?: React.ReactNode,
    children?: MenuItem[],
): MenuItem {
    return {
        key,
        icon,
        children,
        label,
        resourcelabel: resourcelabel
    } as MenuItem;
}

export const items: MenuItem[] = [
    getItem('Cluster', '01', 'Cluster'),
    getItem('Applications', '02', 'Applications'),
    getItem('Nodes', '03', 'Node', <PieChartOutlined />),
    getItem('Workloads', 'sub1', 'Workloads', <UserOutlined />, [
        getItem('Overview', '04', 'Overview'),
        getItem('Pods', '05', 'Pod'),
        getItem('Deployments', '06', 'Deployment'),
        getItem('Daemon Sets', '07', 'DaemonSet'),
        getItem('Stateful Sets', '08', 'StatefulSet'),
        getItem('Replica Sets', '09', 'ReplicaSet'),
        getItem('Replication Controllers', '10', 'ReplicationController'),
        getItem('Jobs', '11', 'Job'),
        getItem('Cron Jobs', '12', 'CronJobs'),
    ]),
    getItem('Config', 'sub2', 'Configs', <TeamOutlined />, [
        getItem('Config Maps', '13', 'ConfigMap'),
        getItem('Secrets', '14', 'Secret'),
        getItem('Resource Quotas', '15', 'ResourceQuota'),
        getItem('Limit Ranges', '16', 'LimitRange'),
        getItem('Horizontal Pod Autoscalers', '17', 'HorizontalPodAutoscaler'),
        getItem('Vertical Pod Autoscalers', '18', 'VerticalPodAutoscaler'),
        getItem('Pod Disruption Budgets', '19', 'PodDisruptionBudget'),
        getItem('Priority Classes', '20', 'PriorityClass'),
        getItem('Runtime Classes', '21', 'RuntimeClass'),
        getItem('Leases', '22', 'Leases'),
        getItem('Mutating Webhook Configs', '23', 'MutatingWebhookConfig'),
        getItem('Validating Webhook Configs', '24', 'ValidatingWebhookConfig'),
    ]),
    getItem('Network', 'sub3', 'Network', <TeamOutlined />, [
        getItem('Services', '25', 'Service'),
        getItem('Endpoints', '26', 'Endpoint'),
        getItem('Ingresses', '27', 'Ingress'),
        getItem('Ingress Classes', '28', 'IngressClass'),
        getItem('Network Policies', '29', 'NetworkPolicy'),
        getItem('Port Forwarding', '30', 'PortForwarding'),
    ]),
    getItem('Storage', 'sub4', 'Storage', <TeamOutlined />, [
        getItem('Persistent Volume Claims', '31', 'PersistentVolumeClaim'),
        getItem('Persistent Volumes', '32', 'PersistentVolume'),
        getItem('Storage Classes', '33', 'StorageClass'),
    ]),
    getItem('Namespaces', '34', 'Namespace', <FileOutlined />),
    getItem('Events', '35', 'Event', <FileOutlined />),
    getItem('Helm', 'sub5', 'Helm', <FileOutlined />, [
        getItem('Charts', '36', 'Chart'),
        getItem('Releases', '37', 'Release'),
    ]),
    getItem('Access Control', 'sub6', 'AccessControl', <TeamOutlined />, [
        getItem('Service Accounts', '38', 'ServiceAccount'),
        getItem('Cluster Roles', '39', 'ClusterRole'),
        getItem('Roles', '40', 'Role'),
        getItem('Cluster Role Bindings', '41', 'ClusterRoleBinding'),
        getItem('Role Bindings', '42', 'RoleBinding'),
        getItem('Pod Security Policies', '43', 'PodSecurityPolicy'),
    ]),
    getItem('Custom Resources', 'sub7', 'CustomResource', <TeamOutlined />, [
        getItem('Definitions', '44', 'CustomResourceDefinition')
    ]),
];
