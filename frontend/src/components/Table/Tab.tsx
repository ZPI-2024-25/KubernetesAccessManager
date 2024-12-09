import {Input, Select, Table} from 'antd';
import styles from './Tab.module.css';
import ResourceDetailsDrawer from "../DrawerDetails/ResourceDetailsDrawer.tsx";
import {useState} from "react";
import {HelmDataSourceItem, ResourceDataSourceItem} from "../../types";
import {extractCRDname} from "../../functions/extractCRDname.ts";

const {Search} = Input;

const Tab = ({columns, dataSource, resourceType}: {
    columns: object[],
    dataSource: ResourceDataSourceItem[] | HelmDataSourceItem[],
    resourceType: string
}) => {
    const [selectedRecord, setSelectedRecord] = useState<object | null>(null);
    const [isDrawerVisible, setDrawerVisible] = useState(false);
    const [query, setQuery] = useState<string>('');
    const [selectedNamespace, setSelectedNamespace] = useState<string>('');

    const extractNamespaces = () => {
        const namespaces = new Set<string>();
        dataSource.forEach((record) => namespaces.add(record.namespace ? record.namespace as string : ''));
        return Array.from(namespaces).filter((namespace) => namespace !== '');
    }

    const namespaces = extractNamespaces();

    const filterDataSource = (dataSource: ResourceDataSourceItem[] | HelmDataSourceItem[], selectedNamespace: string, query: string) => {
        const filteredByNamespace = dataSource.filter(item => item.namespace === selectedNamespace || selectedNamespace === '');

        if (dataSource.length > 0 && 'resource' in dataSource[0]) {
            return filteredByNamespace.filter(item => extractCRDname(item as ResourceDataSourceItem).toLowerCase().includes(query.toLowerCase()));
        } else {
            return filteredByNamespace.filter(item => (item.name ? item.name as string : "").toLowerCase().includes(query.toLowerCase()));
        }
    }

    const filteredDataSource = filterDataSource(dataSource, selectedNamespace, query);

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
                <Search
                    className={styles.searchInput}
                    placeholder="Search"
                    onChange={(e) => setQuery(e.target.value)}
                    value={query}
                />

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
