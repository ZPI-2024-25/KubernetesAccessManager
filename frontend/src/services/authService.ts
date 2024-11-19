import {jwtDecode} from 'jwt-decode';

export const scheduleTokenRefresh = (
    onRefreshFailed: () => void,
    onRefreshSuccess: () => void,
    refreshTimeoutRef: React.MutableRefObject<NodeJS.Timeout | null>
) => {
    const accessToken = localStorage.getItem('access_token');
    if (!accessToken) return;

    const decoded = decodeToken(accessToken);
    if (!decoded || !decoded.exp) {
        console.warn('Nie można zdekodować tokenu lub brak daty wygaśnięcia');
        return;
    }

    const currentTime = Math.floor(Date.now() / 1000); // Czas w sekundach
    const tokenExpiry = decoded.exp; // Czas wygaśnięcia tokenu
    const timeToExpire = tokenExpiry - currentTime; // Czas pozostały do wygaśnięcia

    // Odśwież token 30 sekund przed jego wygaśnięciem
    const refreshTime = Math.max(timeToExpire - 30, Math.max( Math.ceil(timeToExpire*0.9), 0)); // Nie ustawiaj ujemnych czasów

    console.log(`Token wygasa za ${timeToExpire} sekund. Odświeżenie nastąpi za ${refreshTime} sekund.`);

    if (refreshTimeoutRef.current) {
        clearTimeout(refreshTimeoutRef.current); // Wyczyść poprzedni timeout
    }

    refreshTimeoutRef.current = setTimeout(async () => {
        const success = await refreshToken();
        if (!success) {
            onRefreshFailed(); // Akcja na wypadek błędu odświeżenia
        } else {
            onRefreshSuccess(); // Akcja na wypadek sukcesu
            scheduleTokenRefresh(onRefreshFailed, onRefreshSuccess, refreshTimeoutRef); // Zaplanuj kolejne odświeżenie
        }
    }, refreshTime * 1000); // Ustaw timeout w milisekundach
};

export const decodeToken = (token: string) => {
    try {
        return jwtDecode(token);
    } catch (error) {
        console.error("Błąd podczas dekodowania tokena:", error);
        return null;
    }
};

export const isTokenExpired = (token: string) => {
    const decoded = decodeToken(token);
    if (!decoded || !decoded.exp) {
        return true; // Traktuj brak daty wygaśnięcia jako token nieważny
    }
    const currentTime = Math.floor(Date.now() / 1000); // Czas w sekundach
    return decoded.exp < currentTime;
};

export const refreshToken = async () => {
    const refreshToken = localStorage.getItem('refresh_token');
    if (!refreshToken) {
        console.warn('Brak refresh tokena, użytkownik musi się zalogować ponownie.');
        return false;
    }

    try {
        const response = await fetch('http://localhost:4000/realms/ZPI-realm/protocol/openid-connect/token', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/x-www-form-urlencoded',
            },
            body: new URLSearchParams({
                grant_type: 'refresh_token',
                refresh_token: refreshToken,
                client_id: 'ZPI-client',
            }),
        });

        if (!response.ok) {
            console.error('Błąd podczas odświeżania tokenu:', await response.text());
            return false;
        }

        const data = await response.json();
        if (data.access_token && data.refresh_token) {
            // Zapisz nowe tokeny w localStorage
            localStorage.setItem('access_token', data.access_token);
            localStorage.setItem('refresh_token', data.refresh_token);
            console.log('Tokeny zostały odświeżone.');
            return true;
        }

        return false;
    } catch (error) {
        console.error('Błąd podczas odświeżania tokenu:', error);
        return false;
    }
};
