import axios from "axios";
import * as Constants from "../../consts/consts.ts";
import {Status} from "../../types";
import {parseApiError} from "../../functions/apiErrorParser.ts";

export async function deleteResource(resourceType: string, resourceName: string, namespace: string): Promise<Status>  {
    try {
        const response = await axios.delete<Status>(`${Constants.K8S_API_URL}/${resourceType}/${resourceName}?namespace=${namespace}`);
        console.log(`DELETE: ${Constants.K8S_API_URL}/${resourceType}/${resourceName}?namespace=${namespace}`);
        console.log("Response data:", response.data);
        return response.data;
    } catch (error) {
        const errorText = parseApiError(error);
        console.error(errorText);
        throw new Error(errorText);
    }
}
