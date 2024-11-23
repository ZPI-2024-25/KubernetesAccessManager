import {ReactNode, useEffect, useState} from 'react';
import {Button, message, Table} from 'antd';
import {fetchReleases} from '../../api';
import {DeleteOutlined,} from "@ant-design/icons";
import {MdOutlineRestore} from "react-icons/md";
import {HelmDataSourceItem, HelmRelease, HelmReleaseList} from "../../types";
import {helmColumns} from "../../consts/HelmColumns.ts";

const HelmTab = ({showRollbackModal, showUninstallModal, setCurrent}: {
    showRollbackModal: () => void,
    showUninstallModal: () => void,
    setCurrent: (release: HelmRelease) => void
}) => {
    const columns = helmColumns.concat([{
        title: 'Actions',
        dataIndex: "",
        key: 'actions',
        width: 150,
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
    const [dataSource, setDataSource] = useState<HelmDataSourceItem[]>([]);

    useEffect(() => {
        const fetchData = async () => {
            try { //*
                const response: HelmReleaseList = await fetchReleases('');

                const dynamicDataSource: HelmDataSourceItem[] = response.map((resource, index) => ({
                    key: index,
                    ...resource,
                }));
                setDataSource(dynamicDataSource);
            } catch (error) { //*
                console.error('Error fetching releases:', error);
                message.error('Failed to fetch releases.', 2);
            }
        };

        fetchData();
    }, []);

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