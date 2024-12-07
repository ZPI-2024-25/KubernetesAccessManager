import {Role, RoleMap, RoleOperation} from "../../types";
import styles from "./RoleMapForm.module.css";
import {Button, Collapse, Form, Input, Modal} from "antd";
import {MdCancel} from "react-icons/md";
import {FaSave} from "react-icons/fa";
import {useState} from "react";
import SubroleSelect from "./SubroleSelect.tsx";
import {convertRoleMapToRoleConfigMap} from "../../functions/roleMapConversions.ts";
import {updateRoles} from "../../api/k8s/updateRoles.ts";
import {useNavigate} from "react-router-dom";
import RoleOperationsTable from "./RoleOperationsTable.tsx";
import {getAuthStatus} from "../../api";
import {Permissions} from "../../types/authTypes.ts";
import * as Constants from "../../consts/consts.ts";
import {useAuth} from "../AuthProvider/AuthProvider.tsx";

const RoleMapForm = ({data}: { data: RoleMap }) => {
    const [form] = Form.useForm();
    const navigate = useNavigate();

    const {setUserPermissions} = useAuth();

    const [roleMap, setRoleMap] = useState(data.data.roleMap);
    const [subroleMap, setSubroleMap] = useState(data.data.subroleMap);

    const newOperation: RoleOperation = {
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
                {
                    (operationType === "permit" && role.permit && role.permit.length > 0) || (operationType === "deny" && role.deny && role.deny.length > 0) ? (
                        <RoleOperationsTable role={role} mapType={mapType}
                                             handleUpdateOperationField={handleUpdateOperationField}
                                             handleRemoveOperation={handleRemoveOperation}
                                             operationType={operationType}/>
                    ) : null
                }
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
                            <div className={styles.roleHeader}>
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

                                <Button
                                    type="text"
                                    danger
                                    icon={<MdCancel style={{fontSize: "120%"}}/>}
                                    onClick={() => {
                                        const updatedRoleMap = [...roleMap];
                                        updatedRoleMap.splice(index, 1);
                                        setRoleMap(updatedRoleMap);
                                    }}/>
                            </div>
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
                            <div className={styles.roleHeader}>
                                <Input
                                    placeholder="Subrole Name"
                                    value={subrole.name}
                                    onChange={(e) => {
                                        const oldName = subrole.name;
                                        const newName = e.target.value;

                                        handleSubroleNameChange(index, newName, oldName);
                                    }}
                                />

                                <Button
                                    type="text"
                                    danger
                                    icon={<MdCancel style={{fontSize: "120%"}}/>}
                                    onClick={() => {
                                        const updatedSubroleMap = [...subroleMap];
                                        updatedSubroleMap.splice(index, 1);
                                        setSubroleMap(updatedSubroleMap);
                                    }}/>
                            </div>
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

    const onCancel = () => {
        Modal.confirm({
            title: 'Are you sure?',
            content: 'Changes you made will be lost.',
            okText: 'Yes',
            cancelText: 'No',
            onOk: () => {
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
            getAuthStatus().then((permissions: Permissions) => {
                setUserPermissions(permissions);
                localStorage.setItem(Constants.PERMISSIONS_STR_KEY, JSON.stringify(permissions));
            }).catch((error) => {
                console.error('Error fetching user status:', error);
            });

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