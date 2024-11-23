import React, { useEffect, useState } from 'react';
import { Button, Table } from 'antd';
import { ApiResponse, fetchResources } from '../../api';
import { formatAge } from "../../functions/formatAge.ts";
import { DeleteOutlined, EditOutlined, PlusOutlined } from "@ant-design/icons";
import { useNavigate } from 'react-router-dom';
import { deleteResource } from "../../api/deleteResource";
import DeleteConfirmModal from "../Confirm/DeleteConfirm.tsx";

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
    const [isModalVisible, setModalVisible] = useState(false);
    const [selectedRecord, setSelectedRecord] = useState<DataSourceItem | null>(null);
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
                    if (column.toLowerCase().includes('age')) {
                        return formatAge(record[column] as string);
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
                            onClick={() => showDeleteModal(record)}
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
        navigate('/create');
    };

    const handleEdit = (record: DataSourceItem) => {
        const resourceType = resourceLabel;
        const namespace = record.namespace as string;
        const resourceName = record.name as string;

        navigate(`/editor`, {
            state: { resourceType, namespace, resourceName },
        });
    };
    // For details in future

    // const handleDetails = (record: DataSourceItem) => {
    //     console.log('Details', record);
    // };

    const showDeleteModal = (record: DataSourceItem) => {
        setSelectedRecord(record);
        setModalVisible(true);
    };

    const handleDeleteConfirm = async () => {
        if (selectedRecord) {
            try {
                await deleteResource(resourceLabel, selectedRecord.name as string, selectedRecord.namespace as string);
                setDataSource((prev) => prev.filter(item => item.key !== selectedRecord.key));
                setModalVisible(false);
                setSelectedRecord(null);
            } catch (error) {
                console.error("DELETE eror:", error);
            }
        }
    };

    const handleCancel = () => {
        setModalVisible(false);
        setSelectedRecord(null);
    };

    return (
        <div>
            <Button
                type="primary"
                icon={<PlusOutlined />}
                onClick={handleAdd}
                size="large"
                style={{
                    position: "fixed",
                    bottom: "16px",
                    right: "16px",
                    zIndex: 1000,
                    borderRadius: "50px",
                    padding: "0 16px",
                }}
            >
                Add
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
                    marginTop: '64px',
                    maxHeight: 'calc(100vh - 80px)',
                    overflowY: 'auto',
                }}
            />
            {selectedRecord && (
                <DeleteConfirmModal
                    visible={isModalVisible}
                    resourceName={selectedRecord.name as string}
                    namespace={selectedRecord.namespace as string}
                    onConfirm={handleDeleteConfirm}
                    onCancel={handleCancel}
                />
            )}
        </div>
    );
};

export default Tab;
