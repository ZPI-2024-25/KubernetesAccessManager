import axios from "axios";
import {ApiError} from "../types";

export const parseApiError = (error: unknown) => {
    if (axios.isAxiosError(error)) {
        if (error.response) {
            const apiError = error.response.data as ApiError;

            if (apiError.code === 404) {
                return apiError.message.split(':')[0];
            }

            return apiError.message;
        } else if (error.request) {
            return 'No response from server.';
        } else {
            return error.message;
        }
    } else {
        return 'An unexpected error occurred.';
    }
}