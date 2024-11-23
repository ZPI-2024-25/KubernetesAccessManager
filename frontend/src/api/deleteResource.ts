import axios from "axios";
import {API_URL} from "../consts/apiConsts.ts";

export const deleteResource = async (resourceType: string, resourceName: string, namespace: string) => {
    try {
        const response = await axios.delete(`${API_URL}/${resourceType}/${resourceName}?namespace=${namespace}`, {
            headers: {
                accept: 'application/json',
            },
        });

        console.log("Deleted Resource:", response.data);

        return response.data.resourceDetails;
    } catch (error) {
        console.error("Error deleting resource:", error);
        throw error;
    }
};
