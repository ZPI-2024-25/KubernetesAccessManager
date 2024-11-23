import axios from 'axios';
import * as Constants from "../consts/consts.ts";

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
    const response = await axios.get<ApiResponse>(`${Constants.API_URL}/${resourceType}`);
    console.log(`GET: ${Constants.API_URL}/${resourceType}`)
    console.log(response.data);
    return response.data;
}
