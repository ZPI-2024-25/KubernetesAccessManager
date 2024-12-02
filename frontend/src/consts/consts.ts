export const API_PREFIX = import.meta.env.VITE_API_URL || 'http://localhost:8080/'
export const K8S_API_URL = `${API_PREFIX}api/v1/k8s`;
export const HELM_API_URL = `${API_PREFIX}api/v1/helm`;
export const AUTH_API_URL = `${API_PREFIX}api/v1/auth`;
export const KEYCLOAK_URL = import.meta.env.VITE_KEYCLOAK_URL || 'http://localhost:4000/'
export const KEYCLOAK_LOGIN_URL = import.meta.env.KEYCLOAK_LOGIN_URL || `${KEYCLOAK_URL}realms/ZPI-realm/protocol/openid-connect/auth?client_id=ZPI-client&response_type=code&scope=openid`
export const KEYCLOAK_LOGOUT_URL = import.meta.env.KEYCLOAK_LOGOUT_URL || `${KEYCLOAK_URL}realms/ZPI-realm/protocol/openid-connect/logout`
export const KEYCLOAK_TOKEN_URL = import.meta.env.KEYCLOAK_TOKEN_URL || `${KEYCLOAK_URL}realms/ZPI-realm/protocol/openid-connect/token`
export const KEYCLOAK_CLIENT_ID = import.meta.env.KEYCLOAK_CLIENT_ID || 'ZPI-client'

export const ACCESS_TOKEN_STR = 'access_token';
export const REFRESH_TOKEN_STR = 'refresh_token';
export const ID_TOKEN_STR = 'id_token';
export const PERMISSIONS_STR_KEY = 'user_permissions';
