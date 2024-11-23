import React, { createContext, useContext, useEffect, useRef, useState } from 'react';
import { scheduleTokenRefresh, decodeToken } from '../../services/authService';
import { useNavigate } from 'react-router-dom';
import { message } from 'antd';
import { KEYCLOAK_LOGIN_URL, KEYCLOAK_LOGOUT_URL } from "../../consts/apiConsts";

type AuthContextType = {
    isLoggedIn: boolean;
    user: { [key: string]: any } | null;
    handleLogin: () => void;
    handleLogout: () => void;
};

const AuthContext = createContext<AuthContextType>({
    isLoggedIn: false,
    user: null,
    handleLogin: () => {},
    handleLogout: () => {},
});

export const AuthProvider: React.FC<{ children: React.ReactNode }> = ({ children }) => {
    const navigate = useNavigate();
    const refreshTimeout = useRef<NodeJS.Timeout | null>(null);
    const [isLoggedIn, setIsLoggedIn] = useState<boolean>(!!localStorage.getItem('access_token'));
    const [user, setUser] = useState<{ [key: string]: any } | null>(null);

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
        window.location.href = `${KEYCLOAK_LOGIN_URL}&redirect_uri=${encodeURIComponent(redirectUri)}`;
    };

    const handleLogout = () => {
        const logoutUrl = `${KEYCLOAK_LOGOUT_URL}?redirect_uri=${encodeURIComponent(window.location.origin)}`;
        localStorage.removeItem('access_token');
        localStorage.removeItem('refresh_token');
        setIsLoggedIn(false);
        setUser(null);
        window.location.href = logoutUrl;
    };

    useEffect(() => {
        const onRefreshFailed = () => {
            console.warn('Failed to refresh token, logging out...');
            localStorage.removeItem('access_token');
            localStorage.removeItem('refresh_token');
            setIsLoggedIn(false);
            setUser(null);
            message.error('Log in again');
        };

        const onRefreshSuccess = () => {
            const token = localStorage.getItem('access_token');
            decodeAndSetUser(token);
        };

        decodeAndSetUser(localStorage.getItem('access_token'));

        scheduleTokenRefresh(onRefreshFailed, onRefreshSuccess, refreshTimeout);

        return () => {
            if (refreshTimeout.current) {
                clearTimeout(refreshTimeout.current);
            }
        };
    }, [navigate]);

    return (
        <AuthContext.Provider value={{ isLoggedIn, user, handleLogin, handleLogout }}>
            {children}
        </AuthContext.Provider>
    );
};

export const useAuth = () => useContext(AuthContext);
