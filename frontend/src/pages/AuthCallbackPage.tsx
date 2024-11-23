import React, { useEffect, useRef } from 'react';
import { useNavigate, useSearchParams } from 'react-router-dom';
import { message } from 'antd';
import axios from 'axios';
import * as Constants from "../consts/consts.ts";

const AuthCallbackPage: React.FC = () => {
    const [searchParams] = useSearchParams();
    const navigate = useNavigate();
    const hasHandledCallback = useRef(false);

    useEffect(() => {
        const handleAuthCallback = async () => {
            if (hasHandledCallback.current) {
                return;
            }
            hasHandledCallback.current = true;

            const code = searchParams.get('code');
            if (!code) {
                message.error('No authorization code in URL');
                navigate('/login');
                return;
            }

            try {
                console.log('Code:', code);
                console.log('Redirect URI:', `${window.location.origin}/auth/callback`);

                const response = await axios.post(
                    `${Constants.KEYCLOAK_TOKEN_URL}`,
                    new URLSearchParams({
                        grant_type: 'authorization_code',
                        code: code,
                        redirect_uri: `${window.location.origin}/auth/callback`,
                        client_id: `${Constants.KEYCLOAK_CLIENT_ID}`,
                    }),
                    {
                        headers: {
                            'Content-Type': 'application/x-www-form-urlencoded',
                        },
                    }
                );

                const data = response.data;

                if (!data.access_token || !data.refresh_token) {
                    throw new Error('Response does not contain tokens!');
                }

                localStorage.setItem(Constants.ACCESS_TOKEN_STR, data.access_token);
                localStorage.setItem(Constants.REFRESH_TOKEN_STR, data.refresh_token);

                message.success('Logged in successfully');
                navigate('/');
            } catch (error) {
                console.error('Error during login:', error);
                message.error('Cannot log in');
                navigate('/login');
            }
        };

        handleAuthCallback();
    }, [searchParams, navigate]);

    return <div>Logging in progress...</div>;
};

export default AuthCallbackPage;
