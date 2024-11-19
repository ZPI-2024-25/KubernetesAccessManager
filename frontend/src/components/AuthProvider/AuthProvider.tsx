import React, { createContext, useContext, useEffect, useRef, useState } from 'react';
import { scheduleTokenRefresh, decodeToken } from '../../services/authService';
import { useNavigate } from 'react-router-dom';
import { message } from 'antd';

type AuthContextType = {
    isLoggedIn: boolean;
    user: { [key: string]: any } | null; // Możesz dostosować typ użytkownika, jeśli masz konkretne dane
};

const AuthContext = createContext<AuthContextType>({
    isLoggedIn: false,
    user: null,
});

export const AuthProvider: React.FC<{ children: React.ReactNode }> = ({ children }) => {
    const navigate = useNavigate();
    const refreshTimeout = useRef<NodeJS.Timeout | null>(null); // Przechowuje timeout do odświeżania
    const [isLoggedIn, setIsLoggedIn] = useState<boolean>(!!localStorage.getItem('access_token'));
    const [user, setUser] = useState<{ [key: string]: any } | null>(null);

    const decodeAndSetUser = (token: string | null) => {
        if (!token) {
            setUser(null);
            return;
        }

        try {
            const decoded = decodeToken(token);
            setUser(decoded); // Przypisujemy odkodowane dane użytkownika do stanu
        } catch (error) {
            console.error('Błąd podczas dekodowania tokena:', error);
            setUser(null);
        }
    };

    useEffect(() => {
        const onRefreshFailed = () => {
            console.warn('Nie udało się odświeżyć tokenu, wylogowywanie...');
            localStorage.removeItem('access_token');
            localStorage.removeItem('refresh_token');
            setIsLoggedIn(false);
            setUser(null);
            message.error('Zaloguj się ponownie');
        };

        const onRefreshSuccess = () => {
            const token = localStorage.getItem('access_token');
            decodeAndSetUser(token); // Dekodujemy i ustawiamy użytkownika po odświeżeniu tokena
            setIsLoggedIn(true);
        };

        // Pierwsze dekodowanie użytkownika na podstawie istniejącego tokena
        decodeAndSetUser(localStorage.getItem('access_token'));

        scheduleTokenRefresh(onRefreshFailed, onRefreshSuccess, refreshTimeout);

        return () => {
            if (refreshTimeout.current) {
                clearTimeout(refreshTimeout.current); // Wyczyść timeout przy demontażu komponentu
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
