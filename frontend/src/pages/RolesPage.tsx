import {useEffect, useState} from "react";
import {getRoles} from "../api/k8s/getRoles.ts";
import {RoleMap} from "../types";
import {convertRoleConfigMapToRoleMap} from "../functions/roleMapConversions.ts";
import RoleMapCollapse from "../components/RoleMap/RoleMapCollapse.tsx";
import {Button, message} from "antd";
import {FaEdit} from "react-icons/fa";
import styles from "./RolesPage.module.css";
import {useNavigate} from "react-router-dom";
import {useAuth} from "../components/AuthProvider/AuthProvider.tsx";
import {hasPermission} from "../functions/authorization.ts";
import {ROLEMAP_NAMESPACE} from "../consts/consts.ts";

const RolesPage = () => {
    const [roleMap, setRoleMap] = useState<RoleMap>();

    const navigate = useNavigate();

    const {permissions} = useAuth();

    useEffect(() => {
        const func = async () => {
            try {
                const response = await getRoles();

                const rolemap = convertRoleConfigMapToRoleMap(response);

                setRoleMap(rolemap);
            } catch (error) {
                if (error instanceof Error) {
                    console.error('Error fetching releases:', error);
                    message.error(error.message, 4, () => {
                        navigate('/');
                    });
                } else {
                    message.error('An unexpected error occurred.', () => {
                        navigate('/');
                    });
                }
            }
        }
        func();
    }, []);

    return (
        <div className={styles.container}>
            {
                ((permissions !== null && hasPermission(permissions, ROLEMAP_NAMESPACE, 'ConfigMap', 'u') && hasPermission(permissions, ROLEMAP_NAMESPACE, 'ConfigMap', 'r')) ?
                        (
                            <div className={styles.editButtonContainer}>
                                <Button type="primary" icon={<FaEdit/>} onClick={() => {
                                    navigate('edit', {
                                        state: {roleMap}
                                    })
                                }}>
                                    Edit Roles
                                </Button>
                            </div>
                        ) :
                        null
                )
            }
            {roleMap && <RoleMapCollapse data={roleMap}/>}
        </div>
    );
};

export default RolesPage;