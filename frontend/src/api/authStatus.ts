import axios from "axios";
import { API_URL } from "../consts/apiConsts.ts";

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
        const response = await axios.get<UserStatus>(`${API_URL}/auth/status`);
        console.log('GET: /auth/status');
        console.log('Response data:', response.data);
        return response.data;
    } catch (error) {
        console.error('Error getting user status', error);
        throw error;
    }
}