import {Role, RoleMap, RoleOperation} from "../../types";
import styles from "./RoleMapForm.module.css";
import {Button, Collapse, Table, Tag} from "antd";
import {MdCancel} from "react-icons/md";
import {FaSave} from "react-icons/fa";
import {capitalizeFirst} from "../../functions/toUpperCaseFirstLetter.ts";
import {useState} from "react";
import {IoMdClose} from "react-icons/io";
import {ColumnsType} from "antd/es/table";

const RoleMapForm = ({data}: { data: RoleMap }) => {
    // const [form] = Form.useForm();

    const [roleMap, setRoleMap] = useState(data.data.roleMap);
    const [subroleMap, setSubroleMap] = useState(data.data.subroleMap);

    const generateColumns = (role: Role, mapType: string): ColumnsType<RoleOperation> => {
        return [
            {
                title: "",
                dataIndex: "remove",
                key: "remove",
                width: "5%",
                render: (_, record: RoleOperation) => (
                    <Button type="text" danger icon={<IoMdClose/>} onClick={
                        () => mapType === "sub" ? removeOperationSub(role, record) : removeOperation(role, record)
                    }>
                    </Button>
                ),
            },
            {
                title: "Namespace",
                dataIndex: "namespace",
                key: "namespace",
                width: "25%",
                render: (namespace: string) => (!namespace || namespace === "*" ? "All" : namespace),
            },
            {
                title: "Resource",
                dataIndex: "resource",
                key: "resource",
                width: "25%",
                render: (resource: string) => (!resource || resource === "*" ? "All" : resource),
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
            },
        ]
    }


    const newOperation = {
        namespace: "",
        resource: "",
        operations: []
    };

    const handleAddPermission = (role: Role) => {
        const roleIndex = roleMap.findIndex(r => r.name === role.name);
        const updatedRoleMap = [...roleMap];

        const updatedRole = {...role};
        updatedRole.permit = [...(updatedRole.permit || []), newOperation];

        updatedRoleMap[roleIndex] = updatedRole;

        setRoleMap(updatedRoleMap);
    }

    const handleAddDeny = (role: Role) => {
        const roleIndex = roleMap.findIndex(r => r.name === role.name);
        const updatedRoleMap = [...roleMap];

        const updatedRole = {...role};
        updatedRole.deny = [...(updatedRole.deny || []), newOperation];

        updatedRoleMap[roleIndex] = updatedRole;

        setRoleMap(updatedRoleMap);
    }

    const handleAddPermissionSub = (role: Role) => {
        const roleIndex = subroleMap.findIndex(r => r.name === role.name);
        const updatedSubRoleMap = [...subroleMap];

        const updatedRole = {...role};
        updatedRole.permit = [...(updatedRole.permit || []), newOperation];

        updatedSubRoleMap[roleIndex] = updatedRole;

        setSubroleMap(updatedSubRoleMap);
    }

    const handleAddDenySub = (role: Role) => {
        const roleIndex = subroleMap.findIndex(r => r.name === role.name);
        const updatedSubRoleMap = [...subroleMap];

        const updatedRole = {...role};
        updatedRole.deny = [...(updatedRole.deny || []), newOperation];

        updatedSubRoleMap[roleIndex] = updatedRole;

        setSubroleMap(updatedSubRoleMap);
    }

    const removeOperation = (role: Role, operation: RoleOperation) => {
        const roleIndex = roleMap.findIndex(r => r.name === role.name);
        const updatedRoleMap = [...roleMap];

        const updatedRole = {...role};
        updatedRole.permit = updatedRole.permit?.filter(op => op !== operation);
        updatedRole.deny = updatedRole.deny?.filter(op => op !== operation);

        updatedRoleMap[roleIndex] = updatedRole;

        setRoleMap(updatedRoleMap);
    }

    const removeOperationSub = (role: Role, operation: RoleOperation) => {
        const roleIndex = subroleMap.findIndex(r => r.name === role.name);
        const updatedSubRoleMap = [...subroleMap];

        const updatedRole = {...role};
        updatedRole.permit = updatedRole.permit?.filter(op => op !== operation);
        updatedRole.deny = updatedRole.deny?.filter(op => op !== operation);

        updatedSubRoleMap[roleIndex] = updatedRole;

        setSubroleMap(updatedSubRoleMap);
    }

    const renderRoleDetails = (role: Role, mapType: string) => (
        <>
            {role.permit && (
                <>
                    <h4>Permitted Operations:</h4>
                    <Table
                        columns={generateColumns(role, mapType)}
                        dataSource={role.permit}
                        rowKey={(record) => `${record.resource}-${record.namespace}`}
                        pagination={false}
                        size="small"
                    />
                    <Button className={styles.addPermissionButton} type="default"
                            onClick={() => mapType == "sub" ? handleAddPermissionSub(role) : handleAddPermission(role)}>
                        Add Permission
                    </Button>
                </>
            )}
            {role.deny && (
                <>
                    <h4>Denied Operations:</h4>
                    <Table
                        columns={generateColumns(role, mapType)}
                        dataSource={role.deny}
                        rowKey={(record) => `${record.resource}-${record.namespace}`}
                        pagination={false}
                        size="small"
                    />
                    <Button className={styles.addPermissionButton} type="default"
                            onClick={() => mapType == "sub" ? handleAddDenySub(role) : handleAddDeny(role)}>
                        Add Deny
                    </Button>
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
                            style={{cursor: "pointer"}}
                        >
                            {subrole}
                        </Tag>
                    ))}
                </>
            )}
        </>
    );

    // const formBody = (
    //     <Form
    //         form={form}
    //         name="roleForm"
    //         layout="vertical"
    //         onFinish={(values) => {
    //             console.log(values);
    //         }}>
    //         <h2>Roles</h2>
    //         <Collapse accordion>
    //             {roleMap.map((role) => (
    //                 <Collapse.Panel header={role.name} key={role.name}>
    //                     {renderRoleDetails(role)}
    //                 </Collapse.Panel>
    //             ))}
    //         </Collapse>
    //
    //         <h2>Subroles</h2>
    //         <Collapse accordion>
    //             {subroleMap.map((subrole) => (
    //                 <Collapse.Panel header={subrole.name} key={subrole.name}>
    //                     {renderRoleDetails(subrole)}
    //                 </Collapse.Panel>
    //             ))}
    //         </Collapse>
    //     </Form>
    // )

    return (
        <>
            <div className={styles.editButtonContainer}>
                <Button type="default" danger icon={<MdCancel/>}>
                    Cancel
                </Button>
                <Button type="primary" icon={<FaSave/>} onClick={() => {
                    console.log(roleMap)
                    console.log(subroleMap)
                }}>
                    Save
                </Button>
            </div>

            <h2>Roles</h2>
            <Collapse accordion>
                {roleMap.map((role) => (
                    <Collapse.Panel header={role.name} key={role.name}>
                        {renderRoleDetails(role, "role")}
                    </Collapse.Panel>
                ))}
            </Collapse>

            <h2>Subroles</h2>
            <Collapse accordion>
                {subroleMap.map((subrole) => (
                    <Collapse.Panel header={subrole.name} key={subrole.name}>
                        {renderRoleDetails(subrole, "sub")}
                    </Collapse.Panel>
                ))}
            </Collapse>
        </>
    );
};

export default RoleMapForm;