import {Table} from 'antd';
import styles from './Tab.module.css';
import ResourceDetailsDrawer from "../DrawerDetails/ResourceDetailsDrawer.tsx";
import {useState} from "react";
import {ResourceColumnType, ResourceDataSourceItem} from "../../types";

const Tab = ({columns, dataSource, resourceType}: {
    columns: ResourceColumnType[],
    dataSource: ResourceDataSourceItem[],
    resourceType: string
}) => {
    const [selectedRecord, setSelectedRecord] = useState<ResourceDataSourceItem | null>(null);
    const [isDrawerVisible, setDrawerVisible] = useState(false);

    const handleRowClick = (record: ResourceDataSourceItem) => {
        setSelectedRecord(record);
        setDrawerVisible(true);
    };

    const handleCloseDrawer = () => {
        setDrawerVisible(false);
        setSelectedRecord(null);
    };
    return (
        <>
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
                className={styles.tab}
                onRow={(record:ResourceDataSourceItem) => ({
                    onClick: () => handleRowClick(record),
                })}
            />
            <ResourceDetailsDrawer
                visible={isDrawerVisible}
                record={selectedRecord}
                onClose={handleCloseDrawer}
                loading={false}
                resourceType={resourceType}
                />
        </>

    );
};

export default Tab;