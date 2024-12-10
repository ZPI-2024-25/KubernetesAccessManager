#!/bin/sh

export KAM_KEYCLOAK_LOGIN_URL="${KAM_KEYCLOAK_LOGIN_URL:=${KAM_KEYCLOAK_URL}/realms/${KAM_KEYCLOAK_REALM}/protocol/openid-connect/auth?client_id=${KAM_KEYCLOAK_CLIENT}&response_type=code&scope=openid}"
export KAM_KEYCLOAK_LOGOUT_URL="${KAM_KEYCLOAK_LOGOUT_URL:=${KAM_KEYCLOAK_URL}/realms/${KAM_KEYCLOAK_REALM}/protocol/openid-connect/logout}"
export KAM_KEYCLOAK_TOKEN_URL="${KAM_KEYCLOAK_TOKEN_URL:=${KAM_KEYCLOAK_URL}/realms/${KAM_KEYCLOAK_REALM}/protocol/openid-connect/token}"
export KAM_ROLMAP_NAME="${KAM_ROLMAP_NAME:=role-map}"
export KAM_ROLEMAP_NAMESPACE="${KAM_ROLEMAP_NAMESPACE:=default}"

for i in $(env | grep KAM_)
do
    key=$(echo $i | cut -d '=' -f 1)
    value=$(echo $i | cut -d '=' -f 2-)
    echo "Replacing ${key} with ${value}"
    find /usr/share/nginx/html -type f \( -name '*.js' -o -name '*.css' \) -exec sed -i "s|${key}|${value}|g" '{}' +
done
