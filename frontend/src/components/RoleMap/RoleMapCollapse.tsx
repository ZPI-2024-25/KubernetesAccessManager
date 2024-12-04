import {Collapse, Tag, Table} from "antd";
import type {ColumnsType} from "antd/es/table";
import {RoleMap, Role, RoleOperation} from "../../types";
import styles from "./RoleMapCollapse.module.css";

const RoleMapCollapse = ({data}: { data: RoleMap }) => {
    const {roleMap, subroleMap} = data.data;

    // Define columns for operation details in the table
    const operationColumns: ColumnsType<RoleOperation> = [
        {
            title: "Resource",
            dataIndex: "resource",
            key: "resource",
        },
        {
            title: "Namespace",
            dataIndex: "namespace",
            key: "namespace",
        },
        {
            title: "Operations",
            dataIndex: "operations",
            key: "operations",
            render: (operations: string[]) => (
                <>
                    {operations.map((op) => (
                        <Tag color="green" key={op}>
                            {op}
                        </Tag>
                    ))}
                </>
            ),
        },
    ];

    // Function to render role details inside Collapse.Panel
    const renderRoleDetails = (role: Role) => (
        <div>
            <h4>Permitted Operations:</h4>
            <Table
                columns={operationColumns}
                dataSource={role.permit || []}
                rowKey={(record) => `${record.resource}-${record.namespace}`}
                pagination={false}
                size="small"
            />
            {role.deny && (
                <>
                    <h4>Denied Operations:</h4>
                    <Table
                        columns={operationColumns}
                        dataSource={role.deny}
                        rowKey={(record) => `${record.resource}-${record.namespace}`}
                        pagination={false}
                        size="small"
                    />
                </>
            )}
            {role.subroles && role.subroles.length > 0 && (
                <div>
                    <h4>Subroles:</h4>
                    {role.subroles.map((subrole) => (
                        <Tag color="blue" key={subrole}>
                            {subrole}
                        </Tag>
                    ))}
                </div>
            )}
        </div>
    );

    return (
        <div className={styles.container}>
            <h2>Roles</h2>
            <Collapse accordion>
                {roleMap.map((role) => (
                    <Collapse.Panel header={role.name} key={role.name}>
                        {renderRoleDetails(role)}
                    </Collapse.Panel>
                ))}
            </Collapse>

            <h2>Subroles</h2>
            <Collapse accordion>
                {subroleMap.map((subrole) => (
                    <Collapse.Panel header={subrole.name} key={subrole.name}>
                        {renderRoleDetails(subrole)}
                    </Collapse.Panel>
                ))}
            </Collapse>
        </div>
    );
};

export default RoleMapCollapse;
