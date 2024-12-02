import axios from "axios";
import {HELM_API_URL} from "../../consts/consts.ts";
import {Status} from "../../types";
import {parseApiError} from "../../functions/apiErrorParser.ts";

export async function deleteRelease(releaseName: string, namespace: string): Promise<Status> {
    try {
        const namespaceQuery = namespace ? `?namespace=${namespace}` : '';

        const response = await axios.delete<Status>(`${HELM_API_URL}/releases/${releaseName}${namespaceQuery}`);
        console.log(`DELETE: ${HELM_API_URL}/releases/${releaseName}${namespaceQuery}`)
        console.log(response.data);
        return response.data;
    } catch (error) {
        const errorText = parseApiError(error);
        console.error(errorText);
        throw new Error(errorText);
    }
}