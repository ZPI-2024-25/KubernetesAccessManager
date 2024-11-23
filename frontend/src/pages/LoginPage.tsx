import "./LoginPage.css";
import { useAuth } from "../components/AuthProvider/AuthProvider";
import { useEffect } from "react";
import { useNavigate } from "react-router-dom";

export const LoginPage = () => {
    const { isLoggedIn } = useAuth();
    const navigate = useNavigate();

    useEffect(() => {
        if (isLoggedIn) {
            navigate("/");
        }
    }, [isLoggedIn, navigate]);

    return (
        <div className="login-page">
            <h1>Ciulu zaloguj się</h1>
        </div>
    );
};
