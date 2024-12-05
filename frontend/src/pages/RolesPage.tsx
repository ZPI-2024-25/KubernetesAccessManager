import {useEffect, useState} from "react";
import {getRoles} from "../api/k8s/getRoles.ts";
import {RoleMap} from "../types";
import {convertRoleConfigMapToRoleMap} from "../functions/roleMapConversions.ts";
import RoleMapCollapse from "../components/RoleMap/RoleMapCollapse.tsx";
import {Button, message} from "antd";
import {FaEdit} from "react-icons/fa";
import styles from "./RolesPage.module.css";
import {useNavigate} from "react-router-dom";

const RolesPage = () => {
    const [roleMap, setRoleMap] = useState<RoleMap>();

    const navigate = useNavigate();

    useEffect(() => {
        const func = async () => {
            try {
                const response = await getRoles();

                console.log(response);
                const rolemap = convertRoleConfigMapToRoleMap(response);
                console.log(rolemap);
                setRoleMap(rolemap);
            } catch (error) {
                if (error instanceof Error) {
                    console.error('Error fetching releases:', error);
                    message.error(error.message, 4);
                } else {
                    message.error('An unexpected error occurred.');
                }
            }
        }
        func();
    }, []);

    return (
        <div className={styles.container}>
            <div className={styles.editButtonContainer}>
                <Button type="primary" icon={<FaEdit/>} onClick={() => {
                    navigate('edit', {
                        state: {roleMap}
                    })
                }}>
                    Edit Roles
                </Button>
            </div>
            {roleMap && <RoleMapCollapse data={roleMap}/>}
        </div>
    );
};

export default RolesPage;