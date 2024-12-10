const isDefaultEnv = (envValue: string, defaultKeySuffix: string) => {
    return !envValue || envValue.startsWith(`KAM_${defaultKeySuffix}`);
};

export const API_PREFIX = import.meta.env.VITE_API_URL || 'http://localhost:8080'
export const K8S_API_URL = `${API_PREFIX}/api/v1/k8s`;
export const HELM_API_URL = `${API_PREFIX}/api/v1/helm`;
export const AUTH_API_URL = `${API_PREFIX}/api/v1/auth`;
export const KEYCLOAK_URL = import.meta.env.VITE_KEYCLOAK_URL || 'http://localhost:4000'
export const KEYCLOAK_CLIENT = import.meta.env.VITE_KEYCLOAK_CLIENTNAME || 'ZPI-client'
export const KEYCLOAK_REALM = import.meta.env.VITE_KEYCLOAK_REALMNAME || 'ZPI-realm'
export const KEYCLOAK_LOGIN_URL = isDefaultEnv(import.meta.env.VITE_KEYCLOAK_LOGIN_URL, "KEYCLOAK_LOGIN_URL")
    ? `${KEYCLOAK_URL}/realms/${KEYCLOAK_REALM}/protocol/openid-connect/auth?client_id=${KEYCLOAK_CLIENT}&response_type=code&scope=openid`
    : import.meta.env.VITE_KEYCLOAK_LOGIN_URL;

export const KEYCLOAK_LOGOUT_URL = isDefaultEnv(import.meta.env.VITE_KEYCLOAK_LOGOUT_URL, "KEYCLOAK_LOGOUT_URL")
    ? `${KEYCLOAK_URL}/realms/${KEYCLOAK_REALM}/protocol/openid-connect/logout`
    : import.meta.env.VITE_KEYCLOAK_LOGOUT_URL;

export const KEYCLOAK_TOKEN_URL = isDefaultEnv(import.meta.env.VITE_KEYCLOAK_TOKEN_URL, "KEYCLOAK_TOKEN_URL")
    ? `${KEYCLOAK_URL}/realms/${KEYCLOAK_REALM}/protocol/openid-connect/token`
    : import.meta.env.VITE_KEYCLOAK_TOKEN_URL;

export const ROLEMAP_NAME = import.meta.env.VITE_ROLEMAP_NAME || 'role-map';
export const ROLEMAP_NAMESPACE = import.meta.env.VITE_ROLEMAP_NAMESPACE || 'default';

export const ACCESS_TOKEN_STR = 'access_token';
export const REFRESH_TOKEN_STR = 'refresh_token';
export const ID_TOKEN_STR = 'id_token';
export const PERMISSIONS_STR_KEY = 'user_permissions';