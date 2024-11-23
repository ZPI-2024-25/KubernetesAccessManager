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