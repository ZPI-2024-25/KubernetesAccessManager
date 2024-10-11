import {Route, Router, Routes} from "react-router-dom";
import {MainPage} from "../pages/MainPage.tsx";
import {LogInPage} from "../pages/LogInPage.tsx";

export const AppRouter = () => {
    return(
        <Router>
            <Routes>
                <Route path="/" element={<MainPage />} />
                <Route path="/login" element={<LogInPage />} />
            </Routes>
        </Router>
);
}
