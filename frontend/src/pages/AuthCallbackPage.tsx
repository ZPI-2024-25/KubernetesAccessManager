import React, { useEffect, useRef } from 'react';
import { useNavigate, useSearchParams } from 'react-router-dom';
import { message } from 'antd';
import {KEYCLOAK_CLIENT_ID, KEYCLOAK_TOKEN_URL} from "../consts/apiConsts.ts";

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

                const response = await fetch(`${KEYCLOAK_TOKEN_URL}`, {
                    method: 'POST',
                    headers: {
                        'Content-Type': 'application/x-www-form-urlencoded',
                    },
                    body: new URLSearchParams({
                        grant_type: 'authorization_code',
                        code: code,
                        redirect_uri: `${window.location.origin}/auth/callback`,
                        client_id: `${KEYCLOAK_CLIENT_ID}`,
                    }),
                });

                if (!response.ok) {
                    const errorText = await response.text();
                    console.error('Keycloak response error:', errorText);
                    throw new Error('Error while communication with Keycloak');
                }

                const data = await response.json();

                if (!data.access_token || !data.refresh_token) {
                    throw new Error('Response does not contains token!');
                }

                localStorage.setItem('access_token', data.access_token);
                localStorage.setItem('refresh_token', data.refresh_token);

                message.success('Log in successfully');
                navigate('/');
            } catch (error) {
                console.error('Error during log in:', error);
                message.error('Cannot log in');
                navigate('/login');
            }
        };

        handleAuthCallback();
    }, [searchParams, navigate]);

    return <div>Logging in progress...</div>;
};

export default AuthCallbackPage;
