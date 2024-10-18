import React from "react";
import {FileOutlined, PieChartOutlined, TeamOutlined, UserOutlined} from "@ant-design/icons";
import {MenuItem} from "../types";

function getItem(
    label: React.ReactNode,
    key: React.Key,
    icon?: React.ReactNode,
    children?: MenuItem[],
): MenuItem {
    return {
        key,
        icon,
        children,
        label,
    } as MenuItem;
}
export const items: MenuItem[] = [
    getItem('Nodes', '1', <PieChartOutlined />),
    getItem('Workloads', 'sub1', <UserOutlined />, [
        getItem('Overview', '3'),
        getItem('Pods', '4'),
        getItem('Deployments', '5'),
        getItem('StatefulSets', '6'),
        getItem('ReplicaSets', '7'),
    ]),
    getItem('Config', 'sub2', <TeamOutlined />, [
        getItem('Config Maps', '9'),
        getItem('Secrets', '10'),
        getItem('Resource Quotas', '11'),
        getItem('Limit Ranges', '12'),
        getItem('Horizontal Pod Autoscalers', '13'),
        getItem('Pod Disruption Budgets', '14'),
        getItem('Priority Classes', '15'),
        getItem('Runtime Classes', '16'),
        getItem('Mutating Webhook Configs', '17'),
        getItem('Validating Webhook Configs', '18'),
    ]),
    getItem('Network', 'sub3', <TeamOutlined />, [
        getItem('Services', '19'),
        getItem('Endpoints', '20'),
        getItem('Ingresses', '21'),
        getItem('Ingress Classes', '22'),
        getItem('Network Policies', '23'),
        getItem('Port Forwarding', '24'),
    ]),
    getItem('Storage', 'sub4', <TeamOutlined />, [
        getItem('Persistent Volume Claims', '25'),
        getItem('Persistent Volumes', '26'),
        getItem('Storage Classes', '27'),
    ]),
    getItem('Namespaces', '28', <FileOutlined />),
    getItem('Events', '29', <FileOutlined />),
    getItem('Helm', 'sub5', <FileOutlined />,[
        getItem('Charts','30'),
        getItem('Releases','31'),
    ]),
    getItem('Access Control', 'sub6', <TeamOutlined />, [
        getItem('Service Accounts', '32'),
        getItem('Cluster Roles', '33'),
        getItem('Roles', '34'),
        getItem('Cluster Role Bindings', '35'),
        getItem('Role Bindings', '36'),
        getItem('Pod Security Policies', '37'),
    ]),
    getItem('Custom Resources', 'sub7')[
        getItem('Definitions','38')
        ]
];