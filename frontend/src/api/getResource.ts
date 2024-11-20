import axios from "axios";
import {API_URL} from "../consts/apiConsts.ts";

export const getResource = async (resourceType: string, resourceName: string, namespace: string) => {
    try {
        const response = await axios.get(`${API_URL}/${resourceType}/${resourceName}?namespace=${namespace}`, {
            headers: {
                accept: 'application/json',
            },
        });

        console.log("Fetched Resource:", response.data);

        return response.data.resourceDetails;
    } catch (error) {
        console.error("Error fetching resource details:", error);
        throw error;
    }
};
