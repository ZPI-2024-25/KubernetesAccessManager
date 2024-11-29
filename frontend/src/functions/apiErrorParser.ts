import axios from "axios";
import {ApiError} from "../types";

export const parseApiError = (error: unknown) => {
    if (axios.isAxiosError(error)) {
        if (error.response) {
            const apiError = error.response.data as ApiError;

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