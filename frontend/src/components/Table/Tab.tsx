import {Table} from 'antd';
import styles from './Tab.module.css';
import ResourceDetailsDrawer from "../DrawerDetails/ResourceDetailsDrawer.tsx";
import {useState} from "react";

const Tab = ({columns, dataSource, resourceType}: {
    columns: object[],
    dataSource: object[],
    resourceType: string
}) => {
    const [selectedRecord, setSelectedRecord] = useState<object | null>(null);
    const [isDrawerVisible, setDrawerVisible] = useState(false);

    const handleRowClick = (record: object) => {
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
                onRow={(record:object) => ({
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