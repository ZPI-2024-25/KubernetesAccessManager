import axios from 'axios';
import { jwtDecode } from 'jwt-decode';
import * as Constants from "../consts/consts.ts";
import {MutableRefObject} from "react";

export const scheduleTokenRefresh = (
    onRefreshFailed: () => void,
    onRefreshSuccess: () => void,
    refreshTimeoutRef: MutableRefObject<NodeJS.Timeout | null>
) => {
    const accessToken = localStorage.getItem(Constants.ACCESS_TOKEN_STR);
    if (!accessToken) return;

    const decoded = decodeToken(accessToken);
    if (!decoded || !decoded.exp) {
        console.warn('Cannot decode token or no expire date in token.');
        return;
    }

    const currentTime = Math.floor(Date.now() / 1000);
    const tokenExpiry = decoded.exp;
    const timeToExpire = tokenExpiry - currentTime;

    const refreshTime = Math.max(
        timeToExpire - 30,
        Math.max(Math.ceil(timeToExpire * 0.9) - 1, 0)
    );

    console.log(`Token expire in ${timeToExpire} seconds. Refresh will occur in ${refreshTime} seconds.`);

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
    const refreshToken = localStorage.getItem(Constants.REFRESH_TOKEN_STR);
    if (!refreshToken) {
        console.warn('No refresh token, user needs to login once again.');
        return false;
    }

    try {
        const response = await axios.post(
            `${Constants.KEYCLOAK_TOKEN_URL}`,
            new URLSearchParams({
                grant_type: Constants.REFRESH_TOKEN_STR,
                refresh_token: refreshToken,
                client_id: `${Constants.KEYCLOAK_CLIENT}`,
            }),
            {
                headers: {
                    'Content-Type': 'application/x-www-form-urlencoded',
                },
            }
        );

        if (response.status !== 200) {
            console.error('Error during token refreshment:', response.data);
            return false;
        }

        const data = response.data;
        if (data.access_token && data.refresh_token) {
            localStorage.setItem(Constants.ACCESS_TOKEN_STR, data.access_token);
            localStorage.setItem(Constants.REFRESH_TOKEN_STR, data.refresh_token);
            console.log('Token refreshed.');
            return true;
        }

        return false;
    } catch (error) {
        console.error('Error during token refreshment:', error);
        return false;
    }
};
