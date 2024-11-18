export const API_PREFIX = import.meta.env.VITE_API_URL || 'http://localhost:8080/'
export const K8S_API_URL = `${API_PREFIX}api/v1/k8s`;
export const HELM_API_URL = `${API_PREFIX}api/v1/helm`;