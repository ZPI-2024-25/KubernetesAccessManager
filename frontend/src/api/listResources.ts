import axios from 'axios';
import {API_URL} from "../consts/apiConsts.ts";

export interface Resource {
    [key: string]: string;
    active: string;
    age: string;
    bindings: string;
    capacity: string;
    claim: string;
    cluster_ip: string;
    completions: string;
    conditions: string;
    containers: string;
    controlled_by: string;
    cpu: string;
    current: string;
    default: string;
    desired: string;
    disk: string;
}

export interface ApiResponse {
    columns: string[];
    resource_list: Resource[];
}

export async function fetchResources(resourceType: string): Promise<ApiResponse> {
    const response = await axios.get<ApiResponse>(`${API_URL}/${resourceType}`);
    console.log(`GET: ${API_URL}/${resourceType}`)
    console.log(response.data);
    return response.data;
}
