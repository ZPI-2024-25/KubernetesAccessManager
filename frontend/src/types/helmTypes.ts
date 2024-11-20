import {ReactNode} from "react";

export interface HelmRelease {
    name: string;
    namespace: string;
    chart: string;
    status: string;
    updated: string;
    revision: string;
    app_version: string;
}

export type HelmReleaseList = HelmRelease[];

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

export interface HelmDataSourceItem extends HelmRelease{
    key: string | number;
}

export interface HelmColumnType {
    title: string;
    dataIndex: string;
    key: string;
    width: number;
    render: (text: ReactNode, record: HelmDataSourceItem) => ReactNode;
}