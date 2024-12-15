import React, {createContext, useContext, useEffect, useRef, useState} from 'react';
import {decodeToken, scheduleTokenRefresh} from '../../services/authService';
import {useNavigate} from 'react-router-dom';
import {message} from 'antd';
import * as Constants from "../../consts/consts.ts";
import { getAuthStatus } from "../../api/auth/authStatus.ts";
import { Permissions } from "../../types/authTypes.ts";

type AuthContextType = {
    isLoggedIn: boolean;
    user: { [key: string]: any } | null;
    permissions: Permissions | null;
    handleLogin: () => void;
    handleLogout: () => void;
    setUserPermissions: (permissions: Permissions) => void;
};

const AuthContext = createContext<AuthContextType>({
    isLoggedIn: false,
    user: null,
    permissions: null,
    handleLogin: () => {},
    handleLogout: () => {},
    setUserPermissions: () => {},
});

export const AuthProvider: React.FC<{ children: React.ReactNode }> = ({ children }) => {
    const navigate = useNavigate();
    const refreshTimeout = useRef<NodeJS.Timeout | null>(null);
    const [isLoggedIn, setIsLoggedIn] = useState<boolean>(!!localStorage.getItem(Constants.ACCESS_TOKEN_STR));
    const [user, setUser] = useState<{ [key: string]: any } | null>(null);
    const [permissions, setUserPermissions] = useState<Permissions | null>(null);

    const decodeAndSetUser = (token: string | null) => {
        if (!token) {
            setUser(null);
            setIsLoggedIn(false);
            return;
        }

        try {
            const decoded = decodeToken(token);
            setUser(decoded);
            setIsLoggedIn(true);
        } catch (error) {
            console.error('Token decode error:', error);
            message.error(`Token decode error: ${error}`);
            setUser(null);
        }
    };

    const handleLogin = () => {
        const redirectUri = `${window.location.origin}/auth/callback`;
        const clientSecret = Constants.KEYCLOAK_CLIENT_SECRET
            ? `&client_secret=${encodeURIComponent(Constants.KEYCLOAK_CLIENT_SECRET)}`
            : '';
        window.location.href = `${Constants.KEYCLOAK_LOGIN_URL}&redirect_uri=${encodeURIComponent(redirectUri)}${clientSecret}`;
    };

    const handleLogout = () => {
        const idToken = localStorage.getItem(Constants.ID_TOKEN_STR);
        window.location.href = `${Constants.KEYCLOAK_LOGOUT_URL}?id_token_hint=${idToken}&post_logout_redirect_uri=${encodeURIComponent(window.location.origin)}`;
        localStorage.removeItem(Constants.ACCESS_TOKEN_STR);
        localStorage.removeItem(Constants.REFRESH_TOKEN_STR);
        localStorage.removeItem(Constants.PERMISSIONS_STR_KEY);
        setUser(null);
        setIsLoggedIn(false);
        setUserPermissions(null);
    };

    useEffect(() => {
        const onRefreshFailed = () => {
            console.warn('Failed to refresh token, logging out...');
            localStorage.removeItem(Constants.ACCESS_TOKEN_STR);
            localStorage.removeItem(Constants.REFRESH_TOKEN_STR);
            localStorage.removeItem(Constants.PERMISSIONS_STR_KEY);
            setIsLoggedIn(false);
            setUser(null);
            setUserPermissions(null);
            message.error('Log in again');
        };

        const onRefreshSuccess = () => {
            const token = localStorage.getItem(Constants.ACCESS_TOKEN_STR);
            decodeAndSetUser(token);
            getAuthStatus().then((permissions: Permissions) => {
                setUserPermissions(permissions);
                localStorage.setItem(Constants.PERMISSIONS_STR_KEY, JSON.stringify(permissions));
            }).catch((error) => {
                console.error('Error fetching user status:', error);
            });
        };

        decodeAndSetUser(localStorage.getItem(Constants.ACCESS_TOKEN_STR));

        setUserPermissions(JSON.parse(localStorage.getItem(Constants.PERMISSIONS_STR_KEY) || 'null'));

        scheduleTokenRefresh(onRefreshFailed, onRefreshSuccess, refreshTimeout);

        return () => {
            if (refreshTimeout.current) {
                clearTimeout(refreshTimeout.current);
            }
        };
    }, [navigate]);

    return (
        <AuthContext.Provider value={{ isLoggedIn, user, permissions, handleLogin, handleLogout, setUserPermissions }}>
            {children}
        </AuthContext.Provider>
    );
};

export const useAuth = () => useContext(AuthContext);
