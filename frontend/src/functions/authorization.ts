import { UserStatus } from '../types/authTypes';

export function hasPermission(userStatus: UserStatus, namespace: string, resource: string, operation: "c" | "r" | "u" | "d" | "l"): boolean {
    const permissions = userStatus.permissions;
    const lookupNs = permissions[namespace] ? namespace : "*";
    const lookupRes = permissions[lookupNs] && permissions[lookupNs][resource] ? resource : "*";

    return permissions[lookupNs] && permissions[lookupNs][lookupRes] && permissions[lookupNs][lookupRes].includes(operation);
}

export function hasPermissionInAnyNamespace(userStatus: UserStatus, resource: string, operation: "c" | "r" | "u" | "d" | "l"): boolean {
    const permissions = userStatus.permissions;
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

export function hasAnyPermissionInAnyNamespace(userStatus: UserStatus, resource: string): boolean {
    const permissions = userStatus.permissions;
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

export function hasPermissionInAnyResource(userStatus: UserStatus, namespace: string, operation: "c" | "r" | "u" | "d" | "l"): boolean {
    const permissions = userStatus.permissions;
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

export function allowedNamespaces(userStatus: UserStatus, operation: "c" | "r" | "u" | "d" | "l", resource: string): string[] {
    const permissions = userStatus.permissions;
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

export function allowedResources(userStatus: UserStatus, operation: "c" | "r" | "u" | "d" | "l", namespace: string): string[] {
    const permissions = userStatus.permissions;
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

export function allowedOperations(userStatus: UserStatus, namespace: string, resource: string): string[] {
    const permissions = userStatus.permissions;
    const lookupNs = permissions[namespace] ? namespace : "*";
    const lookupRes = permissions[lookupNs] && permissions[lookupNs][resource] ? resource : "*";
    return permissions[lookupNs][lookupRes] || [];
}