export interface RoleOperation {
    resource?: string;
    namespace?: string;
    operations?: string[];
}

export interface Role {
    name: string;
    deny?: RoleOperation[];
    permit?: RoleOperation[];
    subroles?: string[];
}

export interface RoleConfigMap {
    apiVersion: string;
    kind: string;
    metadata: {
        name: string;
        namespace: string;
        [key: string]: unknown;
    };
    data: {
        "role-map": string;
        "subrole-map": string;
    };
}

export interface RoleMap {
    apiVersion: string;
    kind: string;
    metadata: {
        name: string;
        namespace: string;
        [key: string]: unknown;
    };
    data: {
        roleMap: Role[];
        subroleMap: Role[];
    };
}

export interface SimpleRole {
    deny?: RoleOperation[];
    permit?: RoleOperation[];
    subroles?: string[];
}