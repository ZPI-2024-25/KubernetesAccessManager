import axios from "axios";
import {HELM_API_URL} from "../../consts/apiConsts.ts";
import {Status} from "../../types";

export async function deleteRelease(releaseName: string, namespace: string): Promise<Status> {
    try {
        const namespaceQuery = namespace ? `?namespace=${namespace}` : '';

        const response = await axios.delete<Status>(`${HELM_API_URL}/releases/${releaseName}${namespaceQuery}`);
        console.log(`DELETE: ${HELM_API_URL}/releases/${releaseName}${namespaceQuery}`)
        console.log(response.data);
        return response.data;
    } catch (error) {
        console.error('Error deleting release:', error);
        throw error;
    }
}