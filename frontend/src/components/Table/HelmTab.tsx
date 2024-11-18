import {ReactNode, useEffect, useState} from 'react';
import './Tab.css';
import {Button, Table} from 'antd';
import {fetchReleases} from '../../api';
import {formatAge} from "../../functions/formatAge.ts";
import {DeleteOutlined,} from "@ant-design/icons";
import { MdOutlineRestore } from "react-icons/md";
import {HelmReleaseList} from "../../types";

interface DataSourceItem {
    key: string | number;
    name: string;
    namespace: string;
    chart: string;
    status: string;
    updated: string;
    revision: string;
    app_version: string;
}

interface ColumnType {
    title: string;
    dataIndex: string;
    key: string;
    width: number;
    render: (text: ReactNode, record: DataSourceItem) => ReactNode;
}

const Tab = () => {
    const [pageSize, setPageSize] = useState<number>(15);

    const columns: ColumnType[] = [{
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
        render: (_text: ReactNode, record: DataSourceItem): ReactNode => {
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
    }, {
        title: 'Actions',
        dataIndex: "",
        key: 'actions',
        width: 150,
        render: (_, record: DataSourceItem) => (
            <div>
                <Button
                    type="link"
                    icon={<MdOutlineRestore />}
                    onClick={() => handleRollback(record)}
                />
                <Button
                    type="link"
                    icon={<DeleteOutlined />}
                    onClick={() => handleDelete(record)}
                    danger
                />
            </div>
        ),
    }]
    const [dataSource, setDataSource] = useState<DataSourceItem[]>([]);

    useEffect(() => {
        const fetchData = async () => {
            const response: HelmReleaseList = await fetchReleases('');

            const dynamicDataSource: DataSourceItem[] = response.map((resource, index) => ({
                key: index,
                ...resource,
            }));
            setDataSource(dynamicDataSource);
        };

        fetchData();
    }, []);

    const handleRollback = (record: DataSourceItem) => {
        console.log("PUT", record);
    };

    const handleDelete = (record: DataSourceItem) => {
        console.log("DELETE", record);
    };

    return (
        // <div className={styles.fullHeightContainer}>
        //     <div className={styles.tableContainer}>
                <Table
                    className="ant-table"
                    columns={columns}
                    dataSource={dataSource}
                    pagination={{pageSize: pageSize}}
                    scroll={{y: 55 * 5}}
                />
        //     </div>
        // </div>
    );
};

export default Tab;