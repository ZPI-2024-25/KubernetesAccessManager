import {useParams} from "react-router";
import {useState} from "react";
import {ResourceDataSourceItem} from "../types";
import {useListResource} from "../hooks/useListResource.ts";
import {Button} from "antd";
import {DeleteOutlined, EditOutlined, PlusOutlined} from "@ant-design/icons";
import {useNavigate} from "react-router-dom";
import DeleteModal from "../components/Modals/DeleteModal.tsx";
import Tab from "../components/Table/Tab.tsx";

const ResourcePage = () => {
    const {resourceType} = useParams();
    const [openDeleteModal, setOpenDeleteModal] = useState(false);
    const [selectedRecord, setSelectedRecord] = useState<ResourceDataSourceItem>();

    const navigate = useNavigate();
    const {columns, dataSource, setDataSource, wasSuccessful} = useListResource(typeof resourceType === "string" ? resourceType : "", "");
    const columnsWithActions = columns.concat({
        dataIndex: "",
        title: 'Actions',
        key: 'actions',
        render: (_, record: ResourceDataSourceItem) => (
            <div>
                <Button
                    type="link"
                    icon={<EditOutlined/>}
                    onClick={() => handleEdit(record)}
                />
                <Button
                    type="link"
                    icon={<DeleteOutlined/>}
                    onClick={() => handleDelete(record)}
                    danger
                />
            </div>
        ),
        width: 100
    });

    const handleDelete = (record: ResourceDataSourceItem) => {
        setSelectedRecord(record);
        setOpenDeleteModal(true);
    }

    const handleEdit = (record: ResourceDataSourceItem) => {
        const namespace = record.namespace as string;
        const resourceName = record.name as string;

        navigate(`/editor`, {
            state: {resourceType, namespace, resourceName},
        });
    };

    const removeRecord = (record: ResourceDataSourceItem) => {
        setDataSource(dataSource.filter(item => item.key !== record.key));
    }

    function handleAdd() {
        navigate(`/create`, {
            state: { resourceType },
        });
    }

    return (
        <>
            <div>
                <Button
                    type="primary"
                    icon={<PlusOutlined/>}
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
                {wasSuccessful ? <Tab columns={columnsWithActions} dataSource={dataSource}/> : <Tab columns={[]} dataSource={[]} />}
            </div>
            <DeleteModal open={openDeleteModal} setOpen={setOpenDeleteModal}
                         resourceType={typeof resourceType === "string" ? resourceType : ""}
                         resource={selectedRecord}
                         removeResource={removeRecord}/>
        </>
    );
};

export default ResourcePage;