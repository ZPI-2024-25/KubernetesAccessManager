import axios from 'axios';
import * as Constants from "../consts/consts.ts";

let interceptorsInitialized = false;

export function initializeAxiosInterceptors() {
    if (interceptorsInitialized) {
        console.warn('Axios interceptors already initialized.');
        return;
    }

    axios.interceptors.request.use(
        (config) => {
            try {
                const token = localStorage.getItem(Constants.ACCESS_TOKEN_STR);
                if (token) {
                    config.headers.Authorization = `Bearer ${token}`;
                }
            } catch (error) {
                console.warn('Cant attach bearer token:', error);
            }
            return config;
        },
        (error) => {
            console.warn('Axios interceptor error:', error);
            return Promise.resolve(error.config || {});
        }
    );

    interceptorsInitialized = true;
    console.info('Axios interceptors initialized.');
}
