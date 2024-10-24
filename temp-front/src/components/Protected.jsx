import React, { useState, useEffect, useRef } from "react";

const Protected = ({ token, client}) => {
    const isRun = useRef(false);

    useEffect(() => {
        if (isRun.current) return;

        isRun.current = true;

        const config = {
            headers: {
                authorization: `Bearer ${token}`,
            },
        };
    }, []);

    const renderClientInfo = () => {
        if (!client) return null;
        return Object.entries(client).map(([key, value]) => (
            <div key={key}>
                <strong>{key}:</strong> {JSON.stringify(value)}
            </div>
        ));
    };

    return (
        <div>
            <h1>Keycloak Client Information</h1>
            {renderClientInfo()}
        </div>
    );
};

export default Protected;