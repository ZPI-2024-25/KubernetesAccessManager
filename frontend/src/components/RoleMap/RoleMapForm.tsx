import {Role, RoleMap, RoleOperation} from "../../types";
import styles from "./RoleMapForm.module.css";
import {Button, Collapse, Form, Input, Modal, Select, Table, Tag} from "antd";
import {MdCancel} from "react-icons/md";
import {FaSave} from "react-icons/fa";
import {useEffect, useState} from "react";
import {IoMdClose} from "react-icons/io";
import {ColumnsType} from "antd/es/table";
import SubroleSelect from "./SubroleSelect.tsx";
import {operationsOptions, resourcesOptions} from "../../consts/roleOptions.ts";
import {convertRoleMapToRoleConfigMap} from "../../functions/roleMapConversions.ts";
import {updateRoles} from "../../api/k8s/updateRoles.ts";
import {useNavigate} from "react-router-dom";

const RoleMapForm = ({data}: { data: RoleMap }) => {
    const [form] = Form.useForm();
    const navigate = useNavigate();

    const [isModified, setIsModified] = useState(false);

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
                render: (namespace: string, record: RoleOperation) => (
                    <Input
                        value={namespace}
                        placeholder="Namespace"
                        onChange={(e) => handleUpdateOperationField(role, record, mapType, 'namespace', e.target.value)}
                    />
                ),
            },
            {
                title: "Resource",
                dataIndex: "resource",
                key: "resource",
                width: "25%",
                render: (resource: string, record: RoleOperation) => (
                    <Select
                        style={{width: "100%"}}
                        value={resource}
                        options={resourcesOptions}
                        onChange={(val) => handleUpdateOperationField(role, record, mapType, 'resource', val)}
                    />
                ),
            },
            {
                title: "Operations",
                dataIndex: "operations",
                key: "operations",
                render: (operations: string[], record: RoleOperation) => (
                    <Select
                        mode="tags"
                        style={{width: "100%"}}
                        value={operations}
                        options={operationsOptions}
                        onChange={(vals: string[]) => {
                            if (vals.includes("*")) {
                                vals = ["*"];
                            } else {
                                vals = vals.filter(v => v !== "*");
                            }

                            const allowedValues = operationsOptions.map(o => o.value);
                            vals = vals.filter(v => allowedValues.includes(v));

                            handleUpdateOperationField(role, record, mapType, 'operations', vals);
                        }}
                        tagRender={({label, closable, onClose}) => {
                            const onPreventMouseDown = (event: React.MouseEvent) => {
                                event.preventDefault();
                                event.stopPropagation();
                            }

                            return (
                                <Tag
                                    color="green"
                                    onMouseDown={onPreventMouseDown}
                                    closable={closable}
                                    onClose={onClose}
                                    style={{marginInlineEnd: 4}}
                                >
                                    {label}
                                </Tag>
                            );
                        }}
                    />
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

    const handleSubroleNameChange = (index: number, newName: string, oldName: string) => {
        const updatedSubRoleMap = [...subroleMap];
        const updatedSubrole = {...updatedSubRoleMap[index]};
        updatedSubrole.name = newName;
        updatedSubRoleMap[index] = updatedSubrole;

        const updatedRoleMap = roleMap.map((r) => {
            if (!r.subroles) return r;
            const updatedSubroles = r.subroles.map((sr) => sr === oldName ? newName : sr);
            return {...r, subroles: updatedSubroles};
        });

        const updatedSubRoleMap2 = updatedSubRoleMap.map((s) => {
            if (!s.subroles) return s;
            const updatedSubroles = s.subroles.map((sr) => sr === oldName ? newName : sr);
            return {...s, subroles: updatedSubroles};
        });

        setRoleMap(updatedRoleMap);
        setSubroleMap(updatedSubRoleMap2);
    };

    function handleUpdateOperationField(
        role: Role,
        operation: RoleOperation,
        mapType: "role" | "sub",
        field: keyof RoleOperation,
        value: string | string[]
    ) {
        const [targetMap, setTargetMap] = getMapAndSetter(mapType);

        const roleIndex = targetMap.findIndex(r => r.name === role.name);
        if (roleIndex === -1) return;

        const updatedMap = [...targetMap];
        const updatedRole = {...role};

        const isPermit = updatedRole.permit?.includes(operation);
        const isDeny = updatedRole.deny?.includes(operation);

        if (isPermit && updatedRole.permit) {
            const opIndex = updatedRole.permit.indexOf(operation);
            const updatedOperation = {...operation, [field]: value};
            updatedRole.permit = [
                ...updatedRole.permit.slice(0, opIndex),
                updatedOperation,
                ...updatedRole.permit.slice(opIndex + 1)
            ];
        }

        if (isDeny && updatedRole.deny) {
            const opIndex = updatedRole.deny.indexOf(operation);
            const updatedOperation = {...operation, [field]: value};
            updatedRole.deny = [
                ...updatedRole.deny.slice(0, opIndex),
                updatedOperation,
                ...updatedRole.deny.slice(opIndex + 1)
            ];
        }

        updatedMap[roleIndex] = updatedRole;
        setTargetMap(updatedMap);
    }


    const renderOperationsTable = (role: Role, operationType: "permit" | "deny", mapType: "role" | "sub") => {
        return (
            <>
                <Table
                    columns={generateColumns(role, mapType)}
                    dataSource={role[operationType]}
                    rowKey={(record) => `${record.resource}-${record.namespace}`}
                    pagination={false}
                    size="small"
                />
                <Button className={styles.addPermissionButton} type="default"
                        onClick={() => handleAddOperation(role, operationType, mapType)}>
                    Add {operationType === "permit" ? "Permission" : "Deny"}
                </Button>
            </>
        )
    }

    const subroles = subroleMap.map((subrole) => {
        return {label: subrole.name, value: subrole.name};
    });

    const renderRoleDetails = (role: Role, mapType: "role" | "sub") => (
        <>

            <h4>Permitted Operations:</h4>
            {renderOperationsTable(role, "permit", mapType)}

            <h4>Denied Operations:</h4>
            {renderOperationsTable(role, "deny", mapType)}

            <>
                <h4>Subroles:</h4>
                <SubroleSelect
                    role={role}
                    mapType={mapType}
                    roleMap={roleMap}
                    setRoleMap={setRoleMap}
                    subroleMap={subroleMap}
                    setSubroleMap={setSubroleMap}
                    subroles={subroles}
                />
            </>

        </>
    );

    const rolesSection = (
        <>
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
        </>
    )

    const subrolesSection = (
        <>
            <h2>Subroles</h2>
            <Collapse accordion>
                {subroleMap.map((subrole, index) => (
                    <Collapse.Panel
                        header={
                            <Input
                                placeholder="Subrole Name"
                                value={subrole.name}
                                onChange={(e) => {
                                    const oldName = subrole.name;
                                    const newName = e.target.value;

                                    handleSubroleNameChange(index, newName, oldName);
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
        </>
    )

    useEffect(() => {
        const handleBeforeUnload = (e: BeforeUnloadEvent) => {
            if (isModified) {
                e.preventDefault();
            }
        };

        window.addEventListener('beforeunload', handleBeforeUnload);
        return () => {
            window.removeEventListener('beforeunload', handleBeforeUnload);
        };
    }, [isModified]);

    const onCancel = () => {
        console.log(isModified);

        if (!isModified) {
            navigate('/Roles');
            return;
        }

        Modal.confirm({
            title: 'Are you sure?',
            content: 'You have unsaved changes. Are you sure you want to discard them?',
            okText: 'Yes, discard',
            cancelText: 'No',
            onOk: () => {
                setIsModified(false);
                navigate('/Roles');
            }
        });
    };

    const onFinish = () => {
       const updatedData = {
           ...data,
            data: {
                roleMap: roleMap,
                subroleMap: subroleMap
            }
       }

        const covertedData = convertRoleMapToRoleConfigMap(updatedData);
        console.log(covertedData);

        updateRoles(covertedData).then(() => {
            setIsModified(false);

            navigate('/Roles');
        }).catch((error) => {
            console.error(error);
        });
    }

    return (
        <Form
            form={form}
            name="roleForm"
            layout="vertical"
            onFinish={() => onFinish()}
            onValuesChange={() => setIsModified(true)}
        >
            <div className={styles.editButtonContainer}>
                <Button type="default" danger icon={<MdCancel/>} onClick={onCancel}>
                    Cancel
                </Button>
                <Button type="primary" icon={<FaSave/>} htmlType="submit">
                    Save
                </Button>
            </div>

            {rolesSection}

            {subrolesSection}
        </Form>

    );
};

export default RoleMapForm;