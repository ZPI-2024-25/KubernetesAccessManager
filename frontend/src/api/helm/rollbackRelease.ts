import axios from 'axios';
import {HELM_API_URL} from "../../consts/consts.ts";
import {Status, HelmRelease} from "../../types";
import {parseApiError} from "../../functions/apiErrorParser.ts";

export async function rollbackRelease(version: number, releaseName: string, namespace: string): Promise<HelmRelease | Status> {
    try {
        const response = await axios.post<HelmRelease | Status>(`${HELM_API_URL}/releases/${releaseName}/rollback?namespace=${namespace}`, {version});
        console.log(`POST: ${HELM_API_URL}/releases/${releaseName}/rollback?namespace=${namespace}`)
        console.log(response.data);
        return response.data;
    } catch (error) {
        const errorText = parseApiError(error);
        console.error(errorText);
        throw new Error(errorText);
    }
}
