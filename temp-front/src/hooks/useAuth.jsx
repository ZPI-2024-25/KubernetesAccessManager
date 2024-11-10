import React, { useState, useEffect, useRef } from "react";
import Keycloak from "keycloak-js";

const client = new Keycloak({
    url: "http://localhost:4000/",
    realm: "ZPI-realm",
    clientId: "ZPI-client",
});

const useAuth = () => {
    const isRun = useRef(false);
    const [token, setToken] = useState(null);
    const [client2, setclient2] = useState(null);
    const [isLogin, setLogin] = useState(false);

    useEffect(() => {
        if (isRun.current) return;

        isRun.current = true;
        client
            .init({
                onLoad: "login-required",
            })
            .then((res) => {
                setLogin(res);
                setToken(client.token);
                setclient2(client)
            });
    }, []);

    return [isLogin, token, client2];
};

export default useAuth;