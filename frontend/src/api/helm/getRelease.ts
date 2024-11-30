import axios from "axios";
import {HELM_API_URL} from "../../consts/consts.ts";
import {HelmRelease} from "../../types";
import {parseApiError} from "../../functions/apiErrorParser.ts";

export async function fetchRelease(releaseName: string, namespace: string): Promise<HelmRelease> {
    try {
        const namespaceQuery = namespace ? `?namespace=${namespace}` : '';

        const response = await axios.get<HelmRelease>(`${HELM_API_URL}/releases/${releaseName}${namespaceQuery}`);
        console.log(`GET: ${HELM_API_URL}/releases/${releaseName}${namespaceQuery}`)
        console.log(response.data);
        return response.data;
    } catch (error) {
        const errorText = parseApiError(error);
        console.error(errorText);
        throw new Error(errorText);
    }
}