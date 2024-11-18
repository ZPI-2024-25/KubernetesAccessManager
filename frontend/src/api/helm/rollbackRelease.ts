import axios from 'axios';
import {HELM_API_URL} from "../../consts/apiConsts.ts";
import {HelmReleaseHistoryList, Status, ReleaseNameRollbackBody} from "../../types";

export async function rollbackRelease(version: number, releaseName: string, namespace: string): Promise<HelmReleaseHistoryList | Status> {
    try {
        const namespaceQuery = namespace ? `?namespace=${namespace}` : '';

        const body: ReleaseNameRollbackBody = {
            version
        }

        const response = await axios.post<HelmReleaseHistoryList | Status>(`${HELM_API_URL}/releases/${releaseName}/rollback${namespaceQuery}`, body);
        console.log(`POST: ${HELM_API_URL}/releases/${releaseName}/rollback${namespaceQuery}`)
        console.log(response.data);
        return response.data;
    } catch (error) {
        console.error('Error rollbacking release:', error);
        throw error;
    }
}
