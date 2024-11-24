import axios from "axios";
import * as Constants from "../consts/consts.ts";
import {ResourceDetails} from "../types/ResourceDetails.ts";

export async function updateResource(resourceType: string, namespace: string, resourceName: string, resourceData: unknown): Promise<ResourceDetails> {
    try {
        const namespaceQuery = namespace ? `?namespace=${namespace}` : '?namespace=default';

        const response = await axios.put<ResourceDetails>(`${Constants.API_URL}/${resourceType}/${resourceName}${namespaceQuery}`, resourceData);
        console.log(`PUT: ${Constants.API_URL}/${resourceType}/${resourceName}${namespaceQuery}`);
        console.log('Request data:', resourceData);
        console.log('Response data:', response.data);
        return response.data;
    } catch (error) {
        console.error('Error updating resource:', error);
        throw error;
    }
}