import {Role, RoleMap, RoleOperation} from "../../types";
import styles from "./RoleMapForm.module.css";
import {Button, Collapse, Form, Input, Select, Table, Tag} from "antd";
import {MdCancel} from "react-icons/md";
import {FaSave} from "react-icons/fa";
import {capitalizeFirst} from "../../functions/toUpperCaseFirstLetter.ts";
import {useState} from "react";
import {IoMdClose} from "react-icons/io";
import {ColumnsType} from "antd/es/table";

const RoleMapForm = ({data}: { data: RoleMap }) => {
    const [form] = Form.useForm();

    const [roleMap, setRoleMap] = useState(data.data.roleMap);
    const [subroleMap, setSubroleMap] = useState(data.data.subroleMap);

    const generateColumns = (role: Role, mapType: "role" | "sub"): ColumnsType<RoleOperation> => {
        return [
            {
                title: "",
                dataIndex: "remove",
                key: "remove",
                width: "5%",
                render: (_, record: RoleOperation) => (
                    <Button type="text" danger icon={<IoMdClose/>} onClick={
                        () => handleRemoveOperation(role, record, mapType)
                    }>
                    </Button>
                ),
            },
            {
                title: "Namespace",
                dataIndex: "namespace",
                key: "namespace",
                width: "25%",
                render: (namespace: string) => (namespace === "*" ? "All" : namespace),
            },
            {
                title: "Resource",
                dataIndex: "resource",
                key: "resource",
                width: "25%",
                render: (resource: string) => (resource === "*" ? "All" : resource),
            },
            {
                title: "Operations",
                dataIndex: "operations",
                key: "operations",
                render: (operations: string[]) => (
                    <>
                        {operations?.length > 0 ? (
                            operations.map((op) => (
                                <Tag color="green" key={op}>
                                    {op === "*" ? "All" : capitalizeFirst(op)}
                                </Tag>
                            ))
                        ) : null
                        }
                    </>
                ),
            },
        ]
    }


    const newOperation: RoleOperation = {
        namespace: "",
        resource: "",
        operations: []
    };

    function getMapAndSetter(mapType: "role" | "sub"): [Role[], (val: Role[]) => void] {
        if (mapType === "role") {
            return [roleMap, setRoleMap];
        } else {
            return [subroleMap, setSubroleMap];
        }
    }

    function handleAddOperation(role: Role, operationType: "permit" | "deny", mapType: "role" | "sub") {
        const [targetMap, setTargetMap] = getMapAndSetter(mapType);

        const roleIndex = targetMap.findIndex(r => r.name === role.name);
        if (roleIndex === -1) return;

        const updatedMap = [...targetMap];
        const updatedRole = {...role};

        updatedRole[operationType] = [...(updatedRole[operationType] || []), newOperation];

        updatedMap[roleIndex] = updatedRole;
        setTargetMap(updatedMap);
    }

    function handleRemoveOperation(role: Role, operation: RoleOperation, mapType: "role" | "sub") {
        const [targetMap, setTargetMap] = getMapAndSetter(mapType);

        const roleIndex = targetMap.findIndex(r => r.name === role.name);
        if (roleIndex === -1) return;

        const updatedMap = [...targetMap];
        const updatedRole = {...role};

        updatedRole.permit = updatedRole.permit?.filter(op => op !== operation);
        updatedRole.deny = updatedRole.deny?.filter(op => op !== operation);

        updatedMap[roleIndex] = updatedRole;
        setTargetMap(updatedMap);
    }

    const subroles = subroleMap.map((subrole) => {
        return {label: subrole.name, value: subrole.name};
    });

    const renderRoleDetails = (role: Role, mapType: "role" | "sub") => (
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
                    onClick={() => handleAddOperation(role, "permit", mapType)}>
                Add Permission
            </Button>


            <h4>Denied Operations:</h4>
            <Table
                columns={generateColumns(role, mapType)}
                dataSource={role.deny}
                rowKey={(record) => `${record.resource}-${record.namespace}`}
                pagination={false}
                size="small"
            />
            <Button className={styles.addPermissionButton} type="default"
                    onClick={() => handleAddOperation(role, "deny", mapType)}>
                Add Deny
            </Button>

            <>
                <h4>Subroles:</h4>
                <Select
                    mode="tags"
                    allowClear
                    style={{width: "100%"}}
                    placeholder="Select subroles"
                    value={(role.subroles ?? [])}
                    onChange={(value) => {
                        if (mapType === "role") {
                            const roleIndex = roleMap.findIndex(r => r.name === role.name);
                            const updatedRoleMap = [...roleMap];

                            const updatedRole = {...role};
                            updatedRole.subroles = value;

                            updatedRoleMap[roleIndex] = updatedRole;

                            setRoleMap(updatedRoleMap);
                        } else {
                            const roleIndex = subroleMap.findIndex(r => r.name === role.name);
                            const updatedSubRoleMap = [...subroleMap];

                            const updatedRole = {...role};
                            updatedRole.subroles = value;

                            updatedSubRoleMap[roleIndex] = updatedRole;

                            setSubroleMap(updatedSubRoleMap);
                        }
                    }}
                    options={subroles}
                    tagRender={({label, value, closable, onClose}) => {
                        const isValid = subroles.some((opt) => opt.value === value);

                        const onPreventMouseDown = (event: React.MouseEvent) => {
                            event.preventDefault();
                            event.stopPropagation();
                        }

                        return (
                            <Tag
                                color={isValid ? "blue" : "red"}
                                onMouseDown={onPreventMouseDown}
                                closable={closable}
                                // onClose={isValid ? onClose : undefined}
                                onClose={onClose}
                                style={{marginInlineEnd: 4}}
                            >
                                {label}
                            </Tag>
                        );
                    }}
                />
            </>

        </>
    );

    return (
        <Form
            form={form}
            name="roleForm"
            layout="vertical"
            onFinish={(values) => {
                console.log(values);
            }}
            initialValues={{roleMap, subroleMap}}
        >
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
                {roleMap.map((role, index) => (
                    <Collapse.Panel
                        header={
                            <Input
                                placeholder="Role Name"
                                value={role.name}
                                onChange={(e) => {
                                    const roleIndex = index;
                                    const updatedRoleMap = [...roleMap];

                                    const updatedRole = {...role};
                                    updatedRole.name = e.target.value;

                                    updatedRoleMap[roleIndex] = updatedRole;

                                    setRoleMap(updatedRoleMap);
                                }}
                            />
                        }
                        key={index}
                    >
                        {renderRoleDetails(role, "role")}
                    </Collapse.Panel>
                ))}
            </Collapse>

            <Button className={styles.addPermissionButton} type="default" onClick={() => {
                setRoleMap([...roleMap, {name: "", permit: [], deny: []}]);
            }}>
                Add Role
            </Button>

            <h2>Subroles</h2>
            <Collapse accordion>
                {subroleMap.map((subrole, index) => (
                    <Collapse.Panel
                        header={
                            <Input
                                placeholder="Role Name"
                                value={subrole.name}
                                onChange={(e) => {
                                    const roleIndex = index;
                                    const updatedSubRoleMap = [...subroleMap];

                                    const updatedSubrole = {...subrole};
                                    updatedSubrole.name = e.target.value;

                                    updatedSubRoleMap[roleIndex] = updatedSubrole;

                                    setSubroleMap(updatedSubRoleMap);
                                }}
                            />
                        }
                        key={index}
                    >
                        {renderRoleDetails(subrole, "sub")}
                    </Collapse.Panel>
                ))}
            </Collapse>

            <Button className={styles.addPermissionButton} type="default" onClick={() => {
                setSubroleMap([...subroleMap, {name: "", permit: [], deny: []}]);
            }}>
                Add subrole
            </Button>
        </Form>

    );
};

export default RoleMapForm;