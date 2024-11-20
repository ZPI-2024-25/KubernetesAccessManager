import {jwtDecode} from 'jwt-decode';
import {KEYCLOAK_CLIENT_ID, KEYCLOAK_TOKEN_URL} from "../consts/apiConsts.ts";

export const scheduleTokenRefresh = (
    onRefreshFailed: () => void,
    onRefreshSuccess: () => void,
    refreshTimeoutRef: React.MutableRefObject<NodeJS.Timeout | null>
) => {
    const accessToken = localStorage.getItem('access_token');
    if (!accessToken) return;

    const decoded = decodeToken(accessToken);
    if (!decoded || !decoded.exp) {
        console.warn('Cannot decode token or no expire date in token.');
        return;
    }

    const currentTime = Math.floor(Date.now() / 1000);
    const tokenExpiry = decoded.exp;
    const timeToExpire = tokenExpiry - currentTime;

    const refreshTime = Math.max(timeToExpire - 30, Math.max( Math.ceil(timeToExpire*0.9) - 1, 0));

    console.log(`Token expire in ${timeToExpire} seconds. Refresh will occur in  ${refreshTime} seconds.`);

    if (refreshTimeoutRef.current) {
        clearTimeout(refreshTimeoutRef.current);
    }

    refreshTimeoutRef.current = setTimeout(async () => {
        const success = await refreshToken();
        if (!success) {
            onRefreshFailed();
        } else {
            onRefreshSuccess();
            scheduleTokenRefresh(onRefreshFailed, onRefreshSuccess, refreshTimeoutRef);
        }
    }, refreshTime * 1000);
};

export const decodeToken = (token: string) => {
    try {
        return jwtDecode(token);
    } catch (error) {
        console.error("Error during token decode:", error);
        return null;
    }
};

export const refreshToken = async () => {
    const refreshToken = localStorage.getItem('refresh_token');
    if (!refreshToken) {
        console.warn('No refresh token, user need to login once again.');
        return false;
    }

    try {
        const response = await fetch(`${KEYCLOAK_TOKEN_URL}`, {
            method: 'POST',
            headers: {
                'Content-Type': 'application/x-www-form-urlencoded',
            },
            body: new URLSearchParams({
                grant_type: 'refresh_token',
                refresh_token: refreshToken,
                client_id: `${KEYCLOAK_CLIENT_ID}`,
            }),
        });

        if (!response.ok) {
            console.error('Error during token refreshment:', await response.text());
            return false;
        }

        const data = await response.json();
        if (data.access_token && data.refresh_token) {
            localStorage.setItem('access_token', data.access_token);
            localStorage.setItem('refresh_token', data.refresh_token);
            console.log('Token refreshed.');
            return true;
        }

        return false;
    } catch (error) {
        console.error('Error during token refreshment:', error);
        return false;
    }
};
