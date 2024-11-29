import axios from 'axios';
import {HELM_API_URL} from "../../consts/consts.ts";
import {HelmReleaseHistoryList} from "../../types";

export async function fetchReleaseHistory(releaseName: string, namespace: string): Promise<HelmReleaseHistoryList> {
    try {
        const namespaceQuery = namespace ? `?namespace=${namespace}` : '';

        const response = await axios.get<HelmReleaseHistoryList>(`${HELM_API_URL}/releases/${releaseName}/history${namespaceQuery}`);
        console.log(`GET: ${HELM_API_URL}/releases/${releaseName}/history${namespaceQuery}`)
        console.log(response.data);
        return response.data;
    } catch (error) {
        console.error('Error fetching release history:', error);
        throw error;
    }
}
