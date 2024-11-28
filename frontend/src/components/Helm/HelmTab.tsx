import {ReactNode} from 'react';
import {Button, Table} from 'antd';
import {DeleteOutlined,} from "@ant-design/icons";
import {MdOutlineRestore} from "react-icons/md";
import {HelmDataSourceItem, HelmRelease} from "../../types";
import {useFetchReleases} from "../../hooks/useFetchReleases.ts";

const HelmTab = ({showRollbackModal, showUninstallModal, setCurrent}: {
    showRollbackModal: () => void,
    showUninstallModal: () => void,
    setCurrent: (release: HelmRelease) => void
}) => {
    const {helmColumns, dataSource} = useFetchReleases();

    const columns = helmColumns.concat([{
        title: 'Actions',
        dataIndex: "",
        key: 'actions',
        width: 60,
        render: (_: ReactNode, record: HelmDataSourceItem): ReactNode => (
            <div>
                <Button
                    type="link"
                    icon={<MdOutlineRestore/>}
                    onClick={() => handleRollback(record)}
                />
                <Button
                    danger
                    type="link"
                    icon={<DeleteOutlined/>}
                    onClick={() => handleDelete(record)}
                />
            </div>
        ),
    }]);

    const handleRollback = (record: HelmDataSourceItem) => {
        setCurrent(record)
        showRollbackModal()
    };

    const handleDelete = async (record: HelmDataSourceItem) => {
        setCurrent(record)
        showUninstallModal()
    };

    return (
        <Table
            columns={columns.map((col, index) =>
                index === columns.length - 1 ? {...col, fixed: 'right'} : col
            )}
            dataSource={dataSource}
            scroll={{x: 'max-content'}}
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
    );
};

export default HelmTab;