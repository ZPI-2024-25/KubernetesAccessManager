import axios from "axios";
import * as Constants from "../../consts/consts.ts";
import {ResourceDetails} from "../../types";
import {parseApiError} from "../../functions/apiErrorParser.ts";

export async function updateResource(resourceType: string, namespace: string, resourceName: string, resourceData: unknown): Promise<ResourceDetails> {
    try {
        const response = await axios.put<ResourceDetails>(`${Constants.K8S_API_URL}/${resourceType}/${resourceName}?namespace=${namespace}`, resourceData);
        console.log(`PUT: ${Constants.K8S_API_URL}/${resourceType}/${resourceName}?namespace=${namespace}`);
        console.log('Request data:', resourceData);
        console.log('Response data:', response.data);
        return response.data;
    } catch (error) {
        const errorText = parseApiError(error);
        console.error(errorText);
        throw new Error(errorText);
    }
}