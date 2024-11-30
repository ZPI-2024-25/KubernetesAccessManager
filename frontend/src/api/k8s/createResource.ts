import axios from "axios";
import * as Constants from "../../consts/consts.ts";
import {ResourceDetails} from "../../types"
import {parseApiError} from "../../functions/apiErrorParser.ts";

export async function createResource(resourceType: string, namespace: string, resourceData: unknown): Promise<ResourceDetails> {
    try {
        const response = await axios.post<ResourceDetails>(`${Constants.K8S_API_URL}/${resourceType}?namespace=${namespace}`, resourceData);
        console.log(`POST: ${Constants.K8S_API_URL}/${resourceType}?namespace=${namespace}`);
        console.log('Request data:', resourceData);
        console.log('Response data:', response.data);
        return response.data;
    } catch (error) {
        const errorText = parseApiError(error);
        console.error(errorText);
        throw new Error(errorText);
    }
}