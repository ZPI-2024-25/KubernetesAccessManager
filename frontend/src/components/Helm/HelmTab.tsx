import {ReactNode, useEffect, useState} from 'react';
import {Button, message, Table} from 'antd';
import {deleteRelease, fetchReleases} from '../../api';
import {DeleteOutlined,} from "@ant-design/icons";
import {MdOutlineRestore} from "react-icons/md";
import {HelmDataSourceItem, HelmRelease, HelmReleaseList} from "../../types";
import {helmColumns} from "../../consts/HelmColumns.ts";

const HelmTab = ({showModal, setCurrent}: { showModal: () => void, setCurrent: (release: HelmRelease) => void }) => {
    const columns = helmColumns.concat([{
        title: 'Actions',
        dataIndex: "",
        key: 'actions',
        width: 150,
        render: (_: ReactNode, record: HelmDataSourceItem): ReactNode => (
            <div>
                <Button
                    type="link"
                    icon={<MdOutlineRestore />}
                    onClick={() => handleRollback(record)}
                />
                <Button
                    danger
                    type="link"
                    color="danger"
                    icon={<DeleteOutlined />}
                    onClick={() => handleDelete(record)}
                    loading={loadingDelete.includes(record.key as string)}
                />
            </div>
        ),
    }]);
    const [dataSource, setDataSource] = useState<HelmDataSourceItem[]>([]);
    const [loadingDelete, setLoadingDelete] = useState<string[]>([]);
    const [messageApi] = message.useMessage();

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
                messageApi.error('Failed to fetch releases.', 2);
            }
        };

        fetchData();
    }, []);

    const handleRollback = (record: HelmDataSourceItem) => {
        setCurrent(record)
        showModal();
    };

    const handleDelete = async (record: HelmDataSourceItem) => {
        const {name, namespace, key} = record;
        if (!name || !namespace) {
            message.error('Invalid release.', 2);
            return;
        }

        setLoadingDelete((prev) => [...prev, key as string]);

        try {
            const status = await deleteRelease(name, namespace);

            if (status.code === 200) {
                message.success('Release deleted successfully.', 2);

                setDataSource(prevData => prevData.filter(item => item.key !== key));
            } else if (status.code === 202) {
                message.loading('Deletion in progress.', 2).then(async () => {
                    const response: HelmReleaseList = await fetchReleases('');
                    const dynamicDataSource: HelmDataSourceItem[] = response.map((resource, index) => ({
                        key: index,
                        ...resource,
                    }));
                    setDataSource(dynamicDataSource);
                });
            } else {
                message.error('Failed to delete release.', 2);
            }
        } catch (error) {
            console.error('Error during deletion:', error);
            message.error('Rollback error.', 2);
        } finally {
            setLoadingDelete(prev => prev.filter(k => k !== key as string));
        }
    };

    return (
        <Table
            className="ant-table"
            columns={columns}
            dataSource={dataSource}
            scroll={{y: 55 * 5}}
            rowKey="key"
        />
    );
};

export default HelmTab;