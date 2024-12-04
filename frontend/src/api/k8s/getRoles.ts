import axios from "axios";
import * as Constants from "../../consts/consts.ts";
import {RoleConfigMap} from "../../types";
import {parseApiError} from "../../functions/apiErrorParser.ts";

export async function getRoles() : Promise<RoleConfigMap> {
    try {
        // TODO: Proper configMap name and namespace
        const response = await axios.get<{resourceDetails: RoleConfigMap}>(`${Constants.K8S_API_URL}/ConfigMap/role-map?namespace=default`);

        console.log("Fetched Roles:", response.data);
        return response.data.resourceDetails;
    } catch (error) {
        const errorText = parseApiError(error);
        console.error(errorText);
        throw new Error(errorText);
    }
}
