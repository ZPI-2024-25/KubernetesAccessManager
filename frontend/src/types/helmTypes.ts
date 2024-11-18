export interface HelmRelease {
    name: string;
    namespace: string;
    chart: string;
    status: string;
    updated: string;
    revision: number;
    app_version: string;
}

export interface HelmReleaseList {
    release_list: HelmRelease[];
}

export interface Status {
    code: number;
    message: string;
    status: string;
}

export interface HelmReleaseHistory {
    revision: number;
    updated: string;
    status: string;
    chart: string;
    app_version: string;
    description: string;
}

export interface HelmReleaseHistoryList {
    history_list: HelmReleaseHistory[];
}

export interface ReleaseNameRollbackBody {
    version: number;
}
