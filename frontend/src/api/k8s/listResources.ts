import axios from 'axios';
import * as Constants from "../../consts/consts.ts";

export interface Resource {
    [key: string]: string;
    name: string,
    namespace: string,
    age: string
}

export interface ApiResponse {
    columns: string[];
    resource_list: Resource[];
}

export async function fetchResources(resourceType: string, namespace?: string): Promise<ApiResponse> {
    const namespaceQuery = namespace ? `?namespace=${namespace}` : '';

    const response = await axios.get<ApiResponse>(`${Constants.K8S_API_URL}/${resourceType}${namespaceQuery}`);
    console.log(`GET: ${Constants.K8S_API_URL}/${resourceType}${namespaceQuery}`)
    console.log(response.data);
    return response.data;
}
