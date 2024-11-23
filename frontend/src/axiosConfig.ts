import axios from 'axios';

let interceptorsInitialized = false;

export function initializeAxiosInterceptors() {
    if (interceptorsInitialized) {
        console.warn('Axios interceptors already initialized.');
        return;
    }

    axios.interceptors.request.use(
        (config) => {
            try {
                const token = localStorage.getItem('access_token');
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
