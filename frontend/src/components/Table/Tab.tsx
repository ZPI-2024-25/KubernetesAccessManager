import React, { useEffect, useState } from 'react';
import {Button, Table} from 'antd';
import { ApiResponse, fetchResources } from '../../api';
import {formatAge} from "../../functions/formatAge.ts";
import {DeleteOutlined, EditOutlined, FolderOutlined, PlusOutlined} from "@ant-design/icons";
import { useNavigate } from 'react-router-dom';


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
    const navigate = useNavigate();

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
                        return formatAge(record[column] as string);
                    }
                    return (text);
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
                        <Button
                            type="link"
                            icon={<FolderOutlined />}
                            onClick={() => handleDetails(record)}
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
    const handleAdd = () => {console.log("POST");
        navigate('/create');
    };

    const handleEdit = (record: DataSourceItem) => {
        const resourceType = resourceLabel;
        const namespace = record.namespace as string;
        const resourceName = record.name as string;
        console.log(resourceType,namespace,resourceName)

        navigate(`/editor`, {
            state: { resourceType, namespace, resourceName },
        });
    };
    const handleDetails = (record: DataSourceItem) => {
        console.log('Details', record);

    };


    const handleDelete = (record: DataSourceItem) => {
        console.log("DELETE", record);
        const resourceType = resourceLabel;
        const namespace = record.namespace as string;
        const resourceName = record.name as string;

        navigate(`/delete`, {
            state: { resourceType, namespace, resourceName },
        });
    };

    return (
        <div>
            <Button
                type="primary"
                icon={<PlusOutlined />}
                onClick={handleAdd}
                style={{ marginBottom: 16 }}
            >
            </Button>
            <Table
                columns={columns.map((col, index) =>
                    index === columns.length - 1 ? { ...col, fixed: 'right' } : col
                )}
                dataSource={dataSource}
                scroll={{ x: 'max-content' }}
                pagination={{
                    showSizeChanger: true,
                    pageSizeOptions: ['10', '20', '50'],
                }}
                style={{
                    maxHeight: 'calc(100vh - 80px)',
                    overflowY: 'auto',
                }}            />
        </div>
    );
};

export default Tab;