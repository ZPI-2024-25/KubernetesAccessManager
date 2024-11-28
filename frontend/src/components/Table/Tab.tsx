import { useState } from 'react';
import { Button, Table } from 'antd';
import { DeleteOutlined, EditOutlined, PlusOutlined } from "@ant-design/icons";
import { useNavigate } from 'react-router-dom';
import { deleteResource } from "../../api/k8s/deleteResource.ts";
import DeleteConfirmModal from "../Modals/DeleteConfirm.tsx";
import {ResourceDataSourceItem} from "../../types";
import {useListResource} from "../../hooks/useListResource.ts";

const Tab = ({ resourcelabel } : {resourcelabel: string}) => {
    const [isModalVisible, setModalVisible] = useState(false);
    const [selectedRecord, setSelectedRecord] = useState<ResourceDataSourceItem | null>(null);
    const navigate = useNavigate();

    const { columns, dataSource, setDataSource } = useListResource(resourcelabel);

    const columnsWithActions = columns.concat({
        dataIndex: "",
        title: 'Actions',
        key: 'actions',
        render: (_, record: ResourceDataSourceItem) => (
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

    const handleAdd = () => {
        const resourceType = resourcelabel;
        navigate(`/create`, {
            state: { resourceType },
        });
    };

    const handleEdit = (record: ResourceDataSourceItem) => {
        const resourceType = resourcelabel;
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

    const showDeleteModal = (record: ResourceDataSourceItem) => {
        setSelectedRecord(record);
        setModalVisible(true);
    };

    const handleDeleteConfirm = async () => {
        if (selectedRecord) {
            try {
                await deleteResource(resourcelabel, selectedRecord.name as string, selectedRecord.namespace as string);
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
                columns={columnsWithActions.map((col, index) =>
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
