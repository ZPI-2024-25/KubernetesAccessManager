import {Table} from 'antd';
import styles from './Tab.module.css';

const Tab = ({columns, dataSource}: {
    columns: object[],
    dataSource: object[]
}) => {
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
            className={styles.tab}
        />
    );
};

export default Tab;