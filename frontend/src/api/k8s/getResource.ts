import axios from "axios";
import * as Constants from "../../consts/consts.ts";
import {ResourceDetails} from "../../types";
import {parseApiError} from "../../functions/apiErrorParser.ts";

export async function getResource(resourceType: string, resourceName: string, namespace: string) : Promise<ResourceDetails> {
    try {
        const response = await axios.get<{resourceDetails: ResourceDetails}>(`${Constants.K8S_API_URL}/${resourceType}/${resourceName}?namespace=${namespace}`, {
            headers: {
                accept: 'application/json',
            },
        });

        console.log("Fetched Resource:", response.data);
        return response.data.resourceDetails;
    } catch (error) {
        const errorText = parseApiError(error);
        console.error(errorText);
        throw new Error(errorText);
    }
}
