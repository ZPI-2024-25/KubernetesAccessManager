import { RoleConfigMap } from "../../types";
import axios from "axios";
import * as Constants from "../../consts/consts.ts";
import { parseApiError } from "../../functions/apiErrorParser.ts";
import { ROLMAP_NAME, ROLEMAP_NAMESPACE } from "../../consts/consts.ts";

export async function updateRoles(roleConfigMap: RoleConfigMap): Promise<RoleConfigMap> {
    try {
        const response = await axios.put<{
            resourceDetails: RoleConfigMap
        }>(`${Constants.K8S_API_URL}/ConfigMap/${ROLMAP_NAME}?namespace=${ROLEMAP_NAMESPACE}`, roleConfigMap);
        console.log(`PUT: ${Constants.K8S_API_URL}/ConfigMap/${ROLMAP_NAME}?namespace=${ROLEMAP_NAMESPACE}`);
        console.log('Request data:', roleConfigMap);
        console.log('Response data:', response.data);
        return response.data.resourceDetails;
    } catch (error) {
        const errorText = parseApiError(error);
        console.error(errorText);
        throw new Error(errorText);
    }
}