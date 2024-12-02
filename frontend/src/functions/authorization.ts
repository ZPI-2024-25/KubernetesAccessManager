import { Permissions, Operation } from '../types/authTypes';

export function allowedOperations(permissions: Permissions, namespace: string, resource: string): string[] {
    const lookupNs = permissions[namespace] ? namespace : "*";
    const lookupRes = permissions[lookupNs] && permissions[lookupNs][resource] ? resource : "*";
    return permissions[lookupNs][lookupRes] || [];
}

export function hasPermission(permissions: Permissions, namespace: string, resource: string, operation: Operation): boolean {
    return allowedOperations(permissions, namespace, resource).includes(operation);
}

export function hasPermissionInAnyNamespace(permissions: Permissions, resource: string, operation: Operation): boolean {
    for (const namespace in permissions) {
        if (permissions[namespace][resource]) {
            if (permissions[namespace][resource].includes(operation)) {
                return true;
            }
        } else if (permissions[namespace]["*"].includes(operation)) {
            return true;
        }
    }
    return false;
}

export function hasAnyPermissionInAnyNamespace(permissions: Permissions, resource: string): boolean {
    for (const namespace in permissions) {
        if (permissions[namespace][resource]) {
            if (permissions[namespace][resource].length > 0) {
                return true;
            }
        } else if (permissions[namespace]["*"].length > 0) {
            return true;
        }
    }
    return false;
}

export function hasPermissionInAnyResource(permissions: Permissions, namespace: string, operation: Operation): boolean {
    for (const resource in permissions[namespace]) {
        if (permissions[namespace][resource].includes(operation)) {
            return true;
        }
    }
    for (const resource in permissions["*"]) {
        if (!permissions[namespace][resource] && permissions["*"][resource].includes(operation)) {
            return true;
        }
    }
    return false;
}

export function allowedNamespaces(permissions: Permissions, operation: Operation, resource: string): string[] {
    const namespaces: string[] = [];
    for (const namespace in permissions) {
        if (permissions[namespace][resource]) {
            if (permissions[namespace][resource].includes(operation)) {
                namespaces.push(namespace);
            }
        } else if (permissions[namespace]["*"].includes(operation)) {
            namespaces.push(namespace);
        }
    }
    return namespaces;
}

export function allowedResources(permissions: Permissions, operation: Operation, namespace: string): string[] {
    const resources: string[] = [];
    for (const resource in permissions[namespace]) {
        if (permissions[namespace][resource].includes(operation)) {
            resources.push(resource);
        }
    }
    for (const resource in permissions["*"]) {
        if (!permissions[namespace][resource] && permissions["*"][resource].includes(operation)) {
            resources.push(resource);
        }
    }
    return resources;
}