#!/bin/sh

: "${KAM_KEYCLOAK_LOGIN_URL:=}"
: "${KAM_KEYCLOAK_LOGOUT_URL:=}"
: "${KAM_KEYCLOAK_TOKEN_URL:=}"
: "${KAM_ROLMAP_NAME:=}"
: "${KAM_ROLEMAP_NAMESPACE:=}"

for i in $(env | grep KAM_)
do
    key=$(echo $i | cut -d '=' -f 1)
    value=$(echo $i | cut -d '=' -f 2-)
    echo "Replacing ${key} with ${value}"
    find /usr/share/nginx/html -type f \( -name '*.js' -o -name '*.css' \) -exec sed -i "s|${key}|${value}|g" '{}' +
done
