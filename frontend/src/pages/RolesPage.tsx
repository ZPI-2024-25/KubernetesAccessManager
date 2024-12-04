import {useEffect, useRef} from "react";
import {getRoles} from "../api/k8s/getRoles.ts";
import {RoleMap} from "../types";
import {convertRoleConfigMapToRoleMap} from "../functions/roleMapConversions.ts";

const RolesPage = () => {
    const roleMap = useRef<RoleMap>();

    useEffect(() => {
        const func = async () => {
            try {
                const response = await getRoles();

                console.log(response);
                const rolemap = convertRoleConfigMapToRoleMap(response);
                console.log(rolemap);
                roleMap.current = rolemap;
            } catch (error) {
                console.error(error);
            }
        }
        func();
    }, []);

    return (
        <div>
            Roles 222
        </div>
    );
};

export default RolesPage;