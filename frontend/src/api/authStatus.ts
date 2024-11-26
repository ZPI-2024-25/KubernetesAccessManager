import axios from "axios";
import * as Constants from "../consts/consts.ts";

export interface UserStatus {
    permissions: {
        [namespace: string]: {
            [resource: string]: string[];
        };
    };
    user: {
        exp: number;
        preferred_username: string;
        email: string;
    };
}

export async function getAuthStatus(): Promise<UserStatus> {
    try {
        const response = await axios.get<UserStatus>(`${Constants.API_URL_AUTH}/status`);
        console.log('GET: /auth/status');
        console.log('Response data:', response.data);
        return response.data;
    } catch (error) {
        console.error('Error getting user status', error);
        throw error;
    }
}