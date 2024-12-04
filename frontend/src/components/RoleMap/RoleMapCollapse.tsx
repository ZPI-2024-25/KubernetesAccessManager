import {useState, useRef, useEffect} from "react";
import {Collapse, Tag, Table} from "antd";
import type {ColumnsType} from "antd/es/table";
import {RoleMap, Role, RoleOperation} from "../../types";
import styles from "./RoleMapCollapse.module.css";
import {capitalizeFirst} from "../../functions/toUpperCaseFirstLetter.ts";

const RoleMapCollapse = ({data}: { data: RoleMap }) => {
    const {roleMap, subroleMap} = data.data;

    const [activeSubrolePanel, setActiveSubrolePanel] = useState<string | undefined>();
    const subroleRefs = useRef<Record<string, HTMLDivElement | null>>({});

    const generateColumns = (roleOperations: RoleOperation[]): ColumnsType<RoleOperation> => {
        const uniqueOperations = Array.from(
            new Set(roleOperations.flatMap((op) => op.operations ?? ["*"]))
        );

        return [
            {
                title: "Namespace",
                dataIndex: "namespace",
                key: "namespace",
                width: "25%",
                render: (namespace: string) => (!namespace || namespace === "*" ? "All" : namespace),
                sorter: (a, b) => (a.namespace ?? "").localeCompare(b.namespace ?? ""),
            },
            {
                title: "Resource",
                dataIndex: "resource",
                key: "resource",
                width: "25%",
                render: (resource: string) => (!resource || resource === "*" ? "All" : resource),
                sorter: (a, b) => (a.resource ?? "").localeCompare(b.resource ?? ""),
            },
            {
                title: "Operations",
                dataIndex: "operations",
                key: "operations",
                render: (operations: string[]) => (
                    <>
                        {(!operations || operations.length === 0) ? (
                            <Tag color="green">All</Tag>
                        ) : (
                            operations.map((op) => (
                                <Tag color="green" key={op}>
                                    {op === "*" ? "All" : capitalizeFirst(op)}
                                </Tag>
                            ))
                        )}
                    </>
                ),
                filters: uniqueOperations.map((op) => ({
                    text: op === "*" ? "All" : capitalizeFirst(op),
                    value: op,
                })),
                onFilter: (value, record) => {
                    if (value === "*") {
                        return (record.operations ?? []).includes("*");
                    }
                    return (record.operations ?? []).includes(value as string);
                },
            },
        ];
    };

    const handleTagClick = (subrole: string) => {
        setActiveSubrolePanel(subrole);
    };

    useEffect(() => {
        if (activeSubrolePanel) {
            const timer = setTimeout(() => {
                const element = subroleRefs.current[activeSubrolePanel];
                if (element) {
                    element.scrollIntoView({behavior: "smooth", block: "start"});
                }
            }, 200);

            return () => clearTimeout(timer);
        }
    }, [activeSubrolePanel]);

    const renderRoleDetails = (role: Role) => (
        <>
            {role.permit && (
                <>
                    <h4>Permitted Operations:</h4>
                    <Table
                        columns={generateColumns(role.permit)}
                        dataSource={role.permit}
                        rowKey={(record) => `${record.resource}-${record.namespace}`}
                        pagination={false}
                        size="small"
                    />
                </>
            )}
            {role.deny && (
                <>
                    <h4>Denied Operations:</h4>
                    <Table
                        columns={generateColumns(role.deny)}
                        dataSource={role.deny}
                        rowKey={(record) => `${record.resource}-${record.namespace}`}
                        pagination={false}
                        size="small"
                    />
                </>
            )}
            {role.subroles && role.subroles.length > 0 && (
                <>
                    <h4>Subroles:</h4>
                    {role.subroles.map((subrole) => (
                        <Tag
                            color="blue"
                            key={subrole}
                            className={styles.roleTag}
                            onClick={() => handleTagClick(subrole)}
                            style={{cursor: "pointer"}}
                        >
                            {subrole}
                        </Tag>
                    ))}
                </>
            )}
        </>
    );

    return (
        <>
            <h2>Roles</h2>
            <Collapse accordion>
                {roleMap.map((role) => (
                    <Collapse.Panel header={role.name} key={role.name}>
                        {renderRoleDetails(role)}
                    </Collapse.Panel>
                ))}
            </Collapse>

            <h2>Subroles</h2>
            <Collapse
                accordion
                activeKey={activeSubrolePanel}
                onChange={(key) => {
                    if (Array.isArray(key) && key.length > 0) {
                        setActiveSubrolePanel(key[0]);
                    } else {
                        setActiveSubrolePanel(undefined);
                    }
                }}
            >
                {subroleMap.map((subrole) => (
                    <Collapse.Panel
                        header={subrole.name}
                        key={subrole.name}
                        ref={(el) => (subroleRefs.current[subrole.name] = el)}
                    >
                        {renderRoleDetails(subrole)}
                    </Collapse.Panel>
                ))}
            </Collapse>
        </>
    );
};

export default RoleMapCollapse;
