import React, { createContext, useContext, useEffect, useRef, useState } from 'react';
import { scheduleTokenRefresh } from '../../services/authService';
import { useNavigate } from 'react-router-dom';
import {message} from "antd";

const AuthContext = createContext({
    isLoggedIn: false,
});

export const AuthProvider: React.FC<{ children: React.ReactNode }> = ({ children }) => {
    const navigate = useNavigate();
    const refreshTimeout = useRef<NodeJS.Timeout | null>(null); // Przechowuje timeout do odświeżania
    const [isLoggedIn, setIsLoggedIn] = useState<boolean>(!!localStorage.getItem('access_token'));

    useEffect(() => {
        const onRefreshFailed = () => {
            console.warn('Nie udało się odświeżyć tokenu, wylogowywanie...');
            localStorage.removeItem('access_token');
            localStorage.removeItem('refresh_token');
            setIsLoggedIn(false);
            message.error('Zaloguj się ponownie')
        };

        const onRefreshSuccess = () => {
            setIsLoggedIn(true);
        };

        scheduleTokenRefresh(onRefreshFailed, onRefreshSuccess, refreshTimeout);

        return () => {
            if (refreshTimeout.current) {
                clearTimeout(refreshTimeout.current); // Wyczyść timeout przy demontażu komponentu
            }
        };
    }, [navigate]);

    return (
        <AuthContext.Provider value={{ isLoggedIn }}>
            {children}
        </AuthContext.Provider>
    );
};

export const useAuth = () => useContext(AuthContext);
