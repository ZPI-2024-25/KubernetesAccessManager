import axios from 'axios';
import {HELM_API_URL} from "../../consts/consts.ts";
import {HelmReleaseList} from "../../types";

export async function fetchReleases(namespace: string): Promise<HelmReleaseList> {
    try {
        const namespaceQuery = namespace ? `?namespace=${namespace}` : '';

        const response = await axios.get<HelmReleaseList>(`${HELM_API_URL}/releases${namespaceQuery}`);
        console.log(`GET: ${HELM_API_URL}/releases${namespaceQuery}`)
        console.log(response.data);
        return response.data;
    } catch (error) {
        console.error('Error fetching releases:', error);
        throw error;
    }
}
