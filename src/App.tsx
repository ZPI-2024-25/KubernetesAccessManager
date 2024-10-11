import { BrowserRouter as Router, Routes, Route, Link } from 'react-router-dom';
import {MainPage} from "./pages/MainPage.tsx";
import {LogInPage} from "./pages/LogInPage.tsx";

function App() {
    return (
        <Router>
            <Routes>
                <Route path="/" element={<MainPage />} />
                <Route path="/login" element={<LogInPage />} />
            </Routes>
        </Router>
    );
}

export default App;
