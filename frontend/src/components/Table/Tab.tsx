import {Table} from 'antd';
import styles from './Tab.module.css';
import DrawerDetails from "../DrawerDetails/DrawerDetails.tsx";
import {useState} from "react";

const Tab = ({columns, dataSource}: {
    columns: object[],
    dataSource: object[]
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
                onRow={(record) => ({
                    onClick: () => handleRowClick(record),
                })}
            />
            <DrawerDetails
                visible={isDrawerVisible}
                record={selectedRecord}
                onClose={handleCloseDrawer} loading={false}            />
        </>

    );
};

export default Tab;