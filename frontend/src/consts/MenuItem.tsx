import React from "react";
import {SiHelm} from "react-icons/si";
import {MdStorage, MdDashboardCustomize, MdAssuredWorkload, MdDriveFileRenameOutline} from "react-icons/md";
import {FaCubes} from "react-icons/fa6";
import {IoServerSharp, IoDocumentText, IoGitNetwork} from "react-icons/io5";

import {MenuItem} from "../types";

function getItem(
    label: React.ReactNode,
    key: React.Key,
    resourcelabel: string,
    icon?: React.ReactNode,
    children?: MenuItem[],
): MenuItem {
    return {
        key,
        icon,
        children,
        label,
        resourceLabel: resourcelabel
    } as MenuItem;
}

export const items: MenuItem[] = [
    getItem('Nodes', '01', 'Node', <IoServerSharp/>),
    getItem('Workloads', 'sub1', 'Workloads', <FaCubes/>, [
        getItem('Pods', '02', 'Pod'),
        getItem('Deployments', '03', 'Deployment'),
        getItem('Daemon Sets', '04', 'DaemonSet'),
        getItem('Stateful Sets', '05', 'StatefulSet'),
        getItem('Replica Sets', '06', 'ReplicaSet'),
        getItem('Jobs', '07', 'Job'),
        getItem('Cron Jobs', '08', 'CronJobs'),
    ]),
    getItem('Config', 'sub2', 'Configs', <IoDocumentText/>, [
        getItem('Config Maps', '09', 'ConfigMap'),
        getItem('Secrets', '10', 'Secret'),
    ]),
    getItem('Network', 'sub3', 'Network', <IoGitNetwork/>, [
        getItem('Services', '11', 'Service'),
        getItem('Ingresses', '12', 'Ingress'),
    ]),
    getItem('Storage', 'sub4', 'Storage', <MdStorage/>, [
        getItem('Persistent Volume Claims', '13', 'PersistentVolumeClaim'),
        getItem('Persistent Volumes', '14', 'PersistentVolume'),
        getItem('Storage Classes', '15', 'StorageClass'),
    ]),
    getItem('Namespaces', '16', 'Namespace', <MdDriveFileRenameOutline/>),
    getItem('Helm', '17', 'Helm', <SiHelm/>),
    getItem('Access Control', 'sub6', 'AccessControl', <MdAssuredWorkload/>, [
        getItem('Service Accounts', '18', 'ServiceAccount'),
        getItem('Cluster Roles', '19', 'ClusterRole'),
        getItem('Cluster Role Bindings', '20', 'ClusterRoleBinding'),
    ]),
    getItem('Custom Resources', 'sub7', 'CustomResource', <MdDashboardCustomize/>, [
        getItem('Definitions', '21', 'CustomResourceDefinition')
    ]),
];
