import React, { useEffect, useState } from 'react';
import { Table } from 'antd';
import { ApiResponse, fetchResources } from '../../api';

const formatResourceAge = (createdAt: string): string => {
    const createdDate = new Date(createdAt);
    const now = new Date();
    const diffInSeconds = Math.floor((now.getTime() - createdDate.getTime()) / 1000);

    if (diffInSeconds < 60) {
        return `${diffInSeconds} seconds ago`;
    }

    const diffInMinutes = Math.floor(diffInSeconds / 60);
    if (diffInMinutes < 60) {
        return `${diffInMinutes} minutes ago`;
    }

    const diffInHours = Math.floor(diffInMinutes / 60);
    if (diffInHours < 24) {
        return `${diffInHours} hours ago`;
    }

    const diffInDays = Math.floor(diffInHours / 24);
    if (diffInDays < 30) {
        return `${diffInDays} days ago`;
    }

    const diffInMonths = Math.floor(diffInDays / 30);
    if (diffInMonths < 12) {
        return `${diffInMonths} months ago`;
    }

    const diffInYears = Math.floor(diffInMonths / 12);
    return `${diffInYears} years ago`;
};

interface TabProps {
    resourceLabel: string;
}

interface DataSourceItem {
    key: string | number;
    [key: string]: unknown;
}

interface ColumnType {
    title: string;
    dataIndex: string;
    key: string;
    width: number;
    render: (text: React.ReactNode, record: DataSourceItem) => React.ReactNode;
}

const Tab: React.FC<TabProps> = ({ resourceLabel }) => {
    const [columns, setColumns] = useState<ColumnType[]>([]);
    const [dataSource, setDataSource] = useState<DataSourceItem[]>([]);

    useEffect(() => {
        if (!resourceLabel) return;

        const fetchData = async () => {
            const response: ApiResponse = await fetchResources(resourceLabel);

            const dynamicColumns: ColumnType[] = response.columns.map((column) => ({
                title: column,
                dataIndex: column,
                key: column,
                width: 150,
                render: (text: React.ReactNode, record: DataSourceItem): React.ReactNode => {
                    // Explicitly assert the type of record[column] as a string
                    if (column.toLowerCase().includes('age')) {
                        return formatResourceAge(record[column] as string);
                    }
                    return text;
                },
            }));

            const dynamicDataSource: DataSourceItem[] = response.resource_list.map((resource, index) => ({
                key: index,
                ...resource,
            }));

            setColumns(dynamicColumns);
            setDataSource(dynamicDataSource);
        };

        fetchData();
    }, [resourceLabel]);

    return (
        <div>
            <Table
                columns={columns}
                dataSource={dataSource}
                scroll={{ x: 'max-content', y: 55 * 5 }}
            />
        </div>
    );
};

export default Tab;
