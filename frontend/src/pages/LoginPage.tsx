import "./LoginPage.css";
import {useAuth} from "../components/AuthProvider/AuthProvider";
import {useEffect} from "react";
import {useNavigate} from "react-router-dom";
import {Button} from "antd";

export const LoginPage = () => {
    const {isLoggedIn, handleLogin} = useAuth();
    const navigate = useNavigate();

    useEffect(() => {
        if (isLoggedIn) {
            navigate("/");
        }
    }, [isLoggedIn, navigate]);

    return (
        <div className="login-page">
            <div className="login-card">
                <div className="login-card-header">
                    Login
                </div>

                <div className="login-card-content">
                    <p>This application uses OIDC for authentication.</p>
                    <p>Click the button below to login.</p>

                    <Button className="login-button" type="link" onClick={handleLogin}>
                        Login
                    </Button>
                </div>
            </div>
        </div>
    );
};
