import axios from "axios";
import {API_URL} from "../consts/apiConsts.ts";
import {ResourceDetails} from "../types/ResourceDetails.ts";



export async function createResource(resourceType: string, namespace: string, resourceData: unknown): Promise<ResourceDetails> {
    try {
        const namespaceQuery = namespace ? `?namespace=${namespace}` : '';

        const response = await axios.post<ResourceDetails>(`${API_URL}/${resourceType}${namespaceQuery}`, resourceData);
        console.log(`POST: ${API_URL}/${resourceType}?namespace=${namespace}`);
        console.log('Request data:', resourceData);
        console.log('Response data:', response.data);
        return response.data;
    } catch (error) {
        console.error('Error creating resource:', error);
        throw error;
    }
}