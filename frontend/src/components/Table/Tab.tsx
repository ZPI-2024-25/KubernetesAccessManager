import {Select, Table} from 'antd';
import styles from './Tab.module.css';
import ResourceDetailsDrawer from "../DrawerDetails/ResourceDetailsDrawer.tsx";
import {useState} from "react";
import {HelmDataSourceItem, ResourceDataSourceItem} from "../../types";

const Tab = ({columns, dataSource, resourceType}: {
    columns: object[],
    dataSource: ResourceDataSourceItem[] | HelmDataSourceItem[],
    resourceType: string
}) => {
    const [selectedRecord, setSelectedRecord] = useState<object | null>(null);
    const [isDrawerVisible, setDrawerVisible] = useState(false);
    const [selectedNamespace, setSelectedNamespace] = useState<string>('');

    const extractNamespaces = () => {
        const namespaces = new Set<string>();
        dataSource.forEach((record) => namespaces.add(record.namespace ? record.namespace as string : ''));
        return Array.from(namespaces).filter((namespace) => namespace !== '');
    }

    const namespaces = extractNamespaces();

    const filteredDataSource = dataSource.filter(item => item.namespace === selectedNamespace || selectedNamespace === '');

    const handleRowClick = (record: object) => {
        setSelectedRecord(record);
        setDrawerVisible(true);
    };

    const handleCloseDrawer = () => {
        setDrawerVisible(false);
        setSelectedRecord(null);
    };

    return (
        <div className={styles.tabContainer}>
            <div className={styles.controlSection}>
                {namespaces.length > 0 && (
                    <Select
                        className={styles.namespaceSelect}
                        showSearch
                        placeholder="Select namespace"
                        optionFilterProp="label"
                        onChange={(value) => setSelectedNamespace(value)}
                        options={namespaces.map((namespace) => ({
                                value: namespace,
                                label: namespace,
                            })
                        ).concat({value: '', label: 'All namespaces'})
                            .sort((a, b) => a.label.localeCompare(b.label))
                        }
                    />
                )}

            </div>

            <Table
                columns={columns.map((col, index) =>
                    index === columns.length - 1 ? {...col, fixed: 'right'} : col
                )}
                dataSource={filteredDataSource}
                scroll={{x: 'max-content'}}
                pagination={{
                    showSizeChanger: true,
                    pageSizeOptions: ['10', '20', '50'],
                }}
                className={styles.tab}
                onRow={(record: object) => ({
                    onClick: () => handleRowClick(record),
                })}
                rowClassName={(record) =>
                    record === selectedRecord
                        ? styles.selectedRow
                        : styles.rowHover
                }
            />
            <ResourceDetailsDrawer
                visible={isDrawerVisible}
                record={selectedRecord}
                onClose={handleCloseDrawer}
                loading={false}
                resourceType={resourceType}
            />
        </div>
    );
};

export default Tab;
