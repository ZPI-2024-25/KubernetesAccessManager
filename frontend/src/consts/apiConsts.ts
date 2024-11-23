export const API_PREFIX = import.meta.env.VITE_API_URL || 'http://localhost:8080/'
export const API_URL = `${API_PREFIX}api/v1/k8s`;
export const KEYCLOAK_URL = import.meta.env.KEYCLOAK_URL || 'http://localhost:4000/'
export const KEYCLOAK_LOGIN_URL = import.meta.env.KEYCLOAK_LOGIN_URL || `${KEYCLOAK_URL}realms/ZPI-realm/protocol/openid-connect/auth?client_id=ZPI-client&response_type=code`
export const KEYCLOAK_LOGOUT_URL = import.meta.env.KEYCLOAK_LOGOUT_URL || `${KEYCLOAK_URL}realms/ZPI-realm/protocol/openid-connect/logout`
export const KEYCLOAK_TOKEN_URL = import.meta.env.KEYCLOAK_TOKEN_URL || `${KEYCLOAK_URL}realms/ZPI-realm/protocol/openid-connect/token`
export const KEYCLOAK_CLIENT_ID = import.meta.env.KEYCLOAK_CLIENT_ID || 'ZPI-client'
