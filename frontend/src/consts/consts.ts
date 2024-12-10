/*
  While adding a new environment variable to this file, also add this variable to env.sh, .env.production, _helpers.tpl
  and values.yaml files according to the scheme.
  Make sure that the environment variables at the vite level have the prefix VITE_, and those at the docker level KAM_.
  Also make sure that the name of any environment variable is not a prefix of another environment variable.

  I apologize from the bottom of my heart and beg for your forgiveness, but it is the frontend's fault.
 */
export const API_PREFIX = import.meta.env.VITE_API_URL || 'http://localhost:8080'
export const K8S_API_URL = `${API_PREFIX}/api/v1/k8s`;
export const HELM_API_URL = `${API_PREFIX}/api/v1/helm`;
export const AUTH_API_URL = `${API_PREFIX}/api/v1/auth`;
export const KEYCLOAK_URL = import.meta.env.VITE_KEYCLOAK_URL || 'http://localhost:4000'
export const KEYCLOAK_CLIENT = import.meta.env.VITE_KEYCLOAK_CLIENTNAME || 'ZPI-client'
export const KEYCLOAK_REALM = import.meta.env.VITE_KEYCLOAK_REALMNAME || 'ZPI-realm'
export const KEYCLOAK_LOGIN_URL = import.meta.env.KEYCLOAK_LOGIN_URL || `${KEYCLOAK_URL}/realms/${KEYCLOAK_REALM}/protocol/openid-connect/auth?client_id=${KEYCLOAK_CLIENT}&response_type=code&scope=openid`
export const KEYCLOAK_LOGOUT_URL = import.meta.env.KEYCLOAK_LOGOUT_URL || `${KEYCLOAK_URL}/realms/${KEYCLOAK_REALM}/protocol/openid-connect/logout`
export const KEYCLOAK_TOKEN_URL = import.meta.env.KEYCLOAK_TOKEN_URL || `${KEYCLOAK_URL}/realms/${KEYCLOAK_REALM}/protocol/openid-connect/token`

export const ROLEMAP_NAME = import.meta.env.VITE_ROLEMAP_NAME || 'role-map';
export const ROLEMAP_NAMESPACE = import.meta.env.VITE_ROLEMAP_NAMESPACE || 'default';

export const ACCESS_TOKEN_STR = 'access_token';
export const REFRESH_TOKEN_STR = 'refresh_token';
export const ID_TOKEN_STR = 'id_token';
export const PERMISSIONS_STR_KEY = 'user_permissions';