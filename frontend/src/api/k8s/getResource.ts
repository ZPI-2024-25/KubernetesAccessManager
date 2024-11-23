import axios from "axios";
import {K8S_API_URL} from "../../consts/apiConsts.ts";
import {stringifyYaml} from "../../functions/jsonYamlFunctions.ts";

export const getResource = async (resourceType: string, resourceName: string, namespace: string) => {
    try {
        const response = await axios.get(`${K8S_API_URL}/${resourceType}/${resourceName}?namespace=${namespace}`, {
            headers: {
                accept: 'application/json',
            },
        });

        console.log("Fetched Resource:", response.data);
        return stringifyYaml(response.data.resourceDetails);
    } catch (error) {
        console.error("Error fetching resource details:", error);
        throw error;
    }
};
