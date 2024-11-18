import axios from 'axios';
import {API_URL} from "../consts/apiConsts.ts";

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

export async function fetchResources(resourceType: string): Promise<ApiResponse> {
    const response = await axios.get<ApiResponse>(`${API_URL}/${resourceType}`);
    console.log(`GET: ${API_URL}/${resourceType}`)
    console.log(response.data);
    return response.data;
}
