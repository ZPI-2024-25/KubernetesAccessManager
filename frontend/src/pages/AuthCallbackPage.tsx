import React, { useEffect, useRef } from 'react';
import { useNavigate, useSearchParams } from 'react-router-dom';
import { message } from 'antd';

const AuthCallbackPage: React.FC = () => {
    const [searchParams] = useSearchParams();
    const navigate = useNavigate();
    const hasHandledCallback = useRef(false); // Flaga zapobiegająca wielokrotnemu wywołaniu

    useEffect(() => {
        const handleAuthCallback = async () => {
            if (hasHandledCallback.current) {
                return; // Zatrzymaj, jeśli funkcja już została wywołana
            }
            hasHandledCallback.current = true;

            const code = searchParams.get('code');
            if (!code) {
                message.error('Brak kodu autoryzacyjnego w URL');
                navigate('/login');
                return;
            }

            try {
                console.log('Code:', code);
                console.log('Redirect URI:', `${window.location.origin}/auth/callback`);

                const response = await fetch('http://localhost:4000/realms/ZPI-realm/protocol/openid-connect/token', {
                    method: 'POST',
                    headers: {
                        'Content-Type': 'application/x-www-form-urlencoded',
                    },
                    body: new URLSearchParams({
                        grant_type: 'authorization_code',
                        code: code,
                        redirect_uri: `${window.location.origin}/auth/callback`,
                        client_id: 'ZPI-client',
                        // client_secret: 'your-client-secret', // Dodaj jeśli wymagany
                    }),
                });

                if (!response.ok) {
                    const errorText = await response.text();
                    console.error('Błąd odpowiedzi Keycloak:', errorText);
                    throw new Error('Błąd podczas komunikacji z Keycloak');
                }

                const data = await response.json();

                // Walidacja tokenów
                if (!data.access_token || !data.refresh_token) {
                    throw new Error('Odpowiedź nie zawiera wymaganych tokenów');
                }

                // Zapisz tokeny do localStorage
                localStorage.setItem('access_token', data.access_token);
                localStorage.setItem('refresh_token', data.refresh_token);

                message.success('Zalogowano pomyślnie');
                navigate('/'); // Przekierowanie na stronę główną
            } catch (error) {
                console.error('Błąd podczas logowania:', error);
                message.error('Nie udało się zalogować');
                navigate('/login');
            }
        };

        handleAuthCallback();
    }, [searchParams, navigate]);

    return <div>Przetwarzanie logowania...</div>;
};

export default AuthCallbackPage;
