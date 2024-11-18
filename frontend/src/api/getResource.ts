import axios from "axios";

export const getResource = async (resourceType: string, resourceName: string) => {
    try {
        const response = await axios.get(`/api/v1/k8s/${resourceType}/${resourceName}`, {
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
