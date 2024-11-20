import React, { createContext, useContext, useEffect, useRef, useState } from 'react';
import { scheduleTokenRefresh, decodeToken } from '../../services/authService';
import { useNavigate } from 'react-router-dom';
import { message } from 'antd';

type AuthContextType = {
    isLoggedIn: boolean;
    user: { [key: string]: any } | null;
};

const AuthContext = createContext<AuthContextType>({
    isLoggedIn: false,
    user: null,
});

export const AuthProvider: React.FC<{ children: React.ReactNode }> = ({ children }) => {
    const navigate = useNavigate();
    const refreshTimeout = useRef<NodeJS.Timeout | null>(null);
    const [isLoggedIn, setIsLoggedIn] = useState<boolean>(!!localStorage.getItem('access_token'));
    const [user, setUser] = useState<{ [key: string]: any } | null>(null);

    const decodeAndSetUser = (token: string | null) => {
        if (!token) {
            setUser(null);
            return;
        }

        try {
            const decoded = decodeToken(token);
            setUser(decoded);
            setIsLoggedIn(true);
        } catch (error) {
            console.error('Token decode error:', error);
            setUser(null);
        }
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
            setIsLoggedIn(true);
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
        <AuthContext.Provider value={{ isLoggedIn, user }}>
            {children}
        </AuthContext.Provider>
    );
};

export const useAuth = () => useContext(AuthContext);
