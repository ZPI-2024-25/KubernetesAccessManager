import axios from 'axios';
import {K8S_API_URL} from "../../consts/consts.ts";
import {ResourceList} from "../../types";
import {parseApiError} from "../../functions/apiErrorParser.ts";

export async function fetchResources(resourceType: string, namespace?: string): Promise<ResourceList> {
    try {
        const response = await axios.get<ResourceList>(`${K8S_API_URL}/${resourceType}?namespace=${namespace}`);
        console.log(`GET: ${K8S_API_URL}/${resourceType}?namespace=${namespace}`)
        console.log(response.data);
        return response.data;
    } catch (error){
        const errorText = parseApiError(error);
        console.error(errorText);
        throw new Error(errorText);
    }
}
