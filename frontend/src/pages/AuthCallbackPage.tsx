import React, { useEffect } from 'react';
import { useNavigate, useSearchParams } from 'react-router-dom';
import { message } from 'antd';

const AuthCallbackPage: React.FC = () => {
    const [searchParams] = useSearchParams();
    const navigate = useNavigate();

    useEffect(() => {
        const handleAuthCallback = async () => {
            const code = searchParams.get('code');
            if (!code) {
                message.error('Brak kodu autoryzacyjnego w URL');
                navigate('/login');
                return;
            }

            try {
                // Wyślij kod do Keycloak i odbierz tokeny
                const response = await fetch('http://localhost:4000/realms/ZPI-realm/protocol/openid-connect/token', {
                    method: 'POST',
                    headers: {
                        'Content-Type': 'application/x-www-form-urlencoded',
                    },
                    body: new URLSearchParams({
                        grant_type: 'authorization_code',
                        code: code,
                        redirect_uri: 'http://localhost:5173/auth/callback',
                        client_id: 'ZPI-client'
                    }),
                });

                if (!response.ok) {
                    throw new Error('Błąd podczas komunikacji z Keycloak');
                }

                const data = await response.json();
                localStorage.setItem('access_token', data.access_token);
                localStorage.setItem('refresh_token', data.refresh_token);

                message.success('Zalogowano pomyślnie');
                navigate('/'); // Przekierowanie na stronę główną
            } catch (error) {
                console.error(error);
                message.error('Nie udało się zalogować');
                navigate('/login');
            }
        };

        handleAuthCallback();
    }, [searchParams, navigate]);

    return <div>Przetwarzanie logowania...</div>;
};

export default AuthCallbackPage;
