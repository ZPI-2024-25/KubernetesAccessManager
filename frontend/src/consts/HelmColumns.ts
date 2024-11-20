import {ReactNode} from "react";
import {formatAge} from "../functions/formatAge.ts";
import {HelmColumnType, HelmDataSourceItem} from "../types";

export const helmColumns: HelmColumnType[] = [{
    title: 'Name',
    dataIndex: 'name',
    key: 'name',
    width: 150,
    render: (text: ReactNode): ReactNode => {
        return text;
    }
}, {
    title: 'Namespace',
    dataIndex: 'namespace',
    key: 'namespace',
    width: 150,
    render: (text: ReactNode): ReactNode => {
        return text;
    }
}, {
    title: 'Chart',
    dataIndex: 'chart',
    key: 'chart',
    width: 150,
    render: (text: ReactNode): ReactNode => {
        return text;
    }
}, {
    title: 'Status',
    dataIndex: 'status',
    key: 'status',
    width: 150,
    render: (text: ReactNode): ReactNode => {
        return text;
    }
}, {
    title: 'Updated',
    dataIndex: 'updated',
    key: 'updated',
    width: 150,
    render: (_text: ReactNode, record: HelmDataSourceItem): ReactNode => {
        return formatAge(record.updated as string);
    },
}, {
    title: 'Revision',
    dataIndex: 'revision',
    key: 'revision',
    width: 150,
    render: (text: ReactNode): ReactNode => {
        return text;
    }
}, {
    title: 'App Version',
    dataIndex: 'app_version',
    key: 'app_version',
    width: 150,
    render: (text: ReactNode): ReactNode => {
        return text;
    }
},
]