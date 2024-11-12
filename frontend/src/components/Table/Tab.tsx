import React, { useEffect, useState } from 'react';
import { Table, Button } from 'antd';
import { ApiResponse, fetchResources } from '../../api';
import { EditOutlined, DeleteOutlined, PlusOutlined } from '@ant-design/icons';

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
    width?: number;
    render?: (text: React.ReactNode, record: DataSourceItem) => React.ReactNode;
}

const Tab: React.FC<TabProps> = ({ resourceLabel }) => {
    const [columns, setColumns] = useState<ColumnType[]>([]);
    const [dataSource, setDataSource] = useState<DataSourceItem[]>([]);

    useEffect(() => {
        if (!resourceLabel) return;
        console.log(resourceLabel);

        const fetchData = async () => {
            const response: ApiResponse = await fetchResources(resourceLabel);

            const dynamicColumns: ColumnType[] = response.columns.map((column) => ({
                title: column,
                dataIndex: column,
                key: column,
                width: 150,
                render: (text: React.ReactNode, record: DataSourceItem): React.ReactNode => {
                    if (column.toLowerCase().includes('age')) {
                        return formatResourceAge(record[column] as string);
                    }
                    return text;
                },
            }));

            dynamicColumns.push({
                dataIndex: "",
                title: 'Actions',
                key: 'actions',
                render: (_, record: DataSourceItem) => (
                    <div>
                        <Button
                            type="link"
                            icon={<EditOutlined />}
                            onClick={() => handleEdit(record)}
                        />
                        <Button
                            type="link"
                            icon={<DeleteOutlined />}
                            onClick={() => handleDelete(record)}
                            danger
                        />
                    </div>
                ),
                width: 100
            });

            const dynamicDataSource: DataSourceItem[] = response.resource_list.map((resource, index) => ({
                key: index,
                ...resource,
            }));

            setColumns(dynamicColumns);
            setDataSource(dynamicDataSource);
        };

        fetchData();
    }, [resourceLabel]);

    const handleAdd = () => {
        console.log("POST");
    };

    const handleEdit = (record: DataSourceItem) => {
        console.log("PUT", record);
    };

    const handleDelete = (record: DataSourceItem) => {
        console.log("DELETE", record);
    };

    return (
        <div>
            <Button
                type="primary"
                icon={<PlusOutlined />}
                onClick={handleAdd}
                style={{ marginBottom: 16 }}
            >
                Add
            </Button>
            <Table
                columns={columns}
                dataSource={dataSource}
                scroll={{ x: 'max-content', y: 55 * 5 }}
            />
        </div>
    );
};

export default Tab;
