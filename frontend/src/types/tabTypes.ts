import React, {ReactNode} from "react";
import {HelmRelease} from "./helmTypes.ts";

export interface ResourceDataSourceItem {
    key: string | number;
    [key: string]: unknown;
}

export interface ResourceColumnType {
    title: string;
    dataIndex: string;
    key: string;
    width: number;
    render: (text: React.ReactNode, record: ResourceDataSourceItem) => React.ReactNode;
}

export interface HelmDataSourceItem extends HelmRelease{
    key: string | number;
}

export interface HelmColumnType {
    title: string;
    dataIndex: string;
    key: string;
    width?: number;
    render?: (text: ReactNode, record: HelmDataSourceItem) => ReactNode;
}