import React from "react";
import {FileOutlined, PieChartOutlined, TeamOutlined, UserOutlined} from "@ant-design/icons";
import {MenuItem} from "../types";

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
        resourcelabel
    } as MenuItem;
}
export const items: MenuItem[] = [
    getItem('Nodes', '1','Node', <PieChartOutlined />),
    getItem('Workloads', 'sub1','Workloads', <UserOutlined />, [
        getItem('Overview', '3', 'Overwiew'),
        getItem('Pods', '4', 'Pod'),
        getItem('Deployments', '5', 'Deployment'),
        getItem('Stateful Sets', '6', 'StatefulSet'),
        getItem('Replica Sets', '7', 'ReplicaSet'),
    ]),
    getItem('Config', 'sub2','Configs', <TeamOutlined />, [
        getItem('Config Maps', '9', 'ConfigMap'),
        getItem('Secrets', '10', 'Secret'),
        getItem('Resource Quotas', '11', 'ResourceQuota'),
        getItem('Limit Ranges', '12', 'LimitRange'),
        getItem('Horizontal Pod Autoscalers', '13', 'HorizontalPodAutoscaler'),
        getItem('Pod Disruption Budgets', '14', 'PodDisruptionBudget'),
        getItem('Priority Classes', '15', 'PriorityClass'),
        getItem('Runtime Classes', '16', 'RuntimeClass'),
        getItem('Mutating Webhook Configs', '17', 'MutatingWebhookConfig'),
        getItem('Validating Webhook Configs', '18', 'ValidatingWebhookConfig'),
    ]),
    getItem('Network', 'sub3','Network', <TeamOutlined />, [
        getItem('Services', '19', 'Service'),
        getItem('Endpoints', '20', 'Endpoint'),
        getItem('Ingresses', '21', 'Ingress'),
        getItem('Ingress Classes', '22', 'IngressClass'),
        getItem('Network Policies', '23', 'NetworkPolicy'),
        getItem('Port Forwarding', '24', 'PortForwarding'),
    ]),
    getItem('Storage', 'sub4', 'Storage', <TeamOutlined />, [
        getItem('Persistent Volume Claims', '25', 'PersistentVolumeClaim'),
        getItem('Persistent Volumes', '26', 'PersistentVolume'),
        getItem('Storage Classes', '27', 'Storage Class'),
    ]),
    getItem('Namespaces', '28', 'Namespace', <FileOutlined />),
    getItem('Events', '29', 'Event', <FileOutlined />),
    getItem('Helm', 'sub5', 'Helm', <FileOutlined />,[
        getItem('Charts','30', 'Chart'),
        getItem('Releases','31', 'Release'),
    ]),
    getItem('Access Control', 'sub6', 'AccessControl', <TeamOutlined />, [
        getItem('Service Accounts', '32', 'ServiceAccount'),
        getItem('Cluster Roles', '33', 'ClusterRole'),
        getItem('Roles', '34', 'Role'),
        getItem('Cluster Role Bindings', '35', 'ClusterRoleBinding'),
        getItem('Role Bindings', '36', 'RoleBinding'),
        getItem('Pod Security Policies', '37', 'PodSecurityPolicy'),
    ]),
    getItem('Custom Resources', 'sub7', 'CustomResource', <TeamOutlined />, [
        getItem('Definitions','38', 'Definition')
    ]),
];