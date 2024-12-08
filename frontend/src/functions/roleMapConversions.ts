import { Role, RoleConfigMap, RoleMap } from "../types";
import { parseYaml, stringifyYaml } from "./jsonYamlFunctions.ts";

/**
 * Converts a RoleConfigMap object to a RoleMap object.
 * Assumes that 'role-map' and 'subrole-map' contain YAML, where top-level keys are role names,
 * and values are objects with the fields 'name', 'deny', 'permit', and 'subroles'.
 *
 * @param roleConfigMap - An object of type RoleConfigMap
 * @returns An object of type RoleMap
 */
export const convertRoleConfigMapToRoleMap = (roleConfigMap: RoleConfigMap): RoleMap => {
    const roleMapYaml = roleConfigMap.data["role-map"];
    const subroleMapYaml = roleConfigMap.data["subrole-map"];

    console.log(subroleMapYaml);

    const roleMapObj = roleMapYaml ? parseYaml<Record<string, Role>>(roleMapYaml) : {};
    const subroleMapObj = subroleMapYaml ? parseYaml<Record<string, Role>>(subroleMapYaml) : {};

    const roleMapArray = roleMapObj ? Object.values(roleMapObj) : [];
    const subroleMapArray = subroleMapObj ? Object.values(subroleMapObj) : [];

    return {
        apiVersion: roleConfigMap.apiVersion,
        kind: roleConfigMap.kind,
        metadata: {
            ...roleConfigMap.metadata
        },
        data: {
            roleMap: roleMapArray,
            subroleMap: subroleMapArray
        }
    };
};

/**
 * Converts a RoleMap object to a RoleConfigMap object.
 * Creates objects from the roleMap and subroleMap arrays, where the key is the 'name' field of the role,
 * and the value is the role object. Then converts them to YAML.
 *
 * @param roleMap - An object of type RoleMap
 * @returns An object of type RoleConfigMap
 */
export const convertRoleMapToRoleConfigMap = (roleMap: RoleMap): RoleConfigMap => {
    const { roleMap: roles, subroleMap: subroles } = roleMap.data;

    const roleMapObj: Record<string, Role> = {};
    roles.forEach(r => {
        if (r.name) {
            roleMapObj[r.name] = r;
        }
    });

    const subroleMapObj: Record<string, Role> = {};
    subroles.forEach(sr => {
        if (sr.name) {
            subroleMapObj[sr.name] = sr;
        }
    });

    const roleMapYaml = stringifyYaml(roleMapObj);
    const subroleMapYaml = stringifyYaml(subroleMapObj);

    return {
        apiVersion: roleMap.apiVersion,
        kind: roleMap.kind,
        metadata: {
            ...roleMap.metadata
        },
        data: {
            "role-map": roleMapYaml,
            "subrole-map": subroleMapYaml
        }
    };
};
