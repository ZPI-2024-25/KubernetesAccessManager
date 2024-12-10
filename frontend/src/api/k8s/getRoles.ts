import axios from "axios";
import * as Constants from "../../consts/consts.ts";
import { RoleConfigMap } from "../../types";
import { parseApiError } from "../../functions/apiErrorParser.ts";
import { ROLMAP_NAME, ROLEMAP_NAMESPACE } from "../../consts/consts.ts";

export async function getRoles(): Promise<RoleConfigMap> {
    try {
        const response = await axios.get<{
            resourceDetails: RoleConfigMap
        }>(`${Constants.K8S_API_URL}/ConfigMap/${ROLMAP_NAME}?namespace=${ROLEMAP_NAMESPACE}`);

        console.log("Fetched Roles:", response.data);
        return response.data.resourceDetails;
    } catch (error) {
        const errorText = parseApiError(error);
        console.error(errorText);
        throw new Error(errorText);
    }
}
