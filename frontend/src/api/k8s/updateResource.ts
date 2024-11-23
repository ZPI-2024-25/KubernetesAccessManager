import axios from "axios";
import {K8S_API_URL} from "../../consts/apiConsts.ts";
import {ResourceDetails} from "../types/ResourceDetails.ts";

export async function updateResource(resourceType: string, namespace: string, resourceName: string, resourceData: unknown): Promise<ResourceDetails> {
    try {
        const namespaceQuery = namespace ? `?namespace=${namespace}` : '';

        const response = await axios.put<ResourceDetails>(`${K8S_API_URL}/${resourceType}/${resourceName}${namespaceQuery}`, resourceData);
        console.log(`PUT: ${K8S_API_URL}/${resourceType}/${resourceName}${namespaceQuery}`);
        console.log('Request data:', resourceData);
        console.log('Response data:', response.data);
        return response.data;
    } catch (error) {
        console.error('Error updating resource:', error);
        throw error;
    }
}