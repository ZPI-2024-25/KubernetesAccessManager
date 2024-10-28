import { BrowserRouter as Router, Routes, Route } from 'react-router-dom';
import {MainPage} from "./pages/MainPage.tsx";
import {LoginPage} from "./pages/LoginPage.tsx";
import Menu from "./components/Menu/Menu.tsx";

function App() {
    return (
        <Router>
            <Routes>
                <Route path="/" element={<Menu />}>
                    <Route index element={<MainPage />} />
                </Route>
                <Route path="/login" element={<LoginPage />} />
            </Routes>
        </Router>
    );
}

export default App;
