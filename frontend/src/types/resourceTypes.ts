export interface Resource {
    [key: string]: string;
    name: string,
    namespace: string,
    age: string
}

export interface ResourceList {
    columns: string[];
    resource_list: Resource[];
}

export interface ResourceDetails {
    apiVersion: string;
    kind: string;
    metadata: {
        name: string;
        namespace: string;
        [key: string]: unknown;
    };
    spec?: unknown;
    status?: unknown;
    [key: string]: unknown;
}