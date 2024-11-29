export interface HelmRelease {
    name: string;
    namespace: string;
    chart: string;
    status: string;
    updated: string;
    revision: string;
    app_version: string;
}

export interface HelmReleaseHistory {
    revision: number;
    updated: string;
    status: string;
    chart: string;
    app_version: string;
    description: string;
}