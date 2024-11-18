import { BrowserRouter as Router, Routes, Route } from 'react-router-dom';
import {MainPage} from "./pages/MainPage.tsx";
import Menu from "./components/Menu/Menu.tsx";
import EditorPage from "./pages/EditorPage.tsx";
import ResourcePage from "./pages/ResourcePage.tsx";
import AuthCallbackPage from "./pages/AuthCallbackPage.tsx";

function App() {
    return (
        <Router>
            <Routes>
                <Route path="/" element={<Menu />}>
                    <Route index element={<MainPage />} />
                    <Route path=":resourceType" element={<ResourcePage />} />
                </Route>
                <Route path="/auth/callback" element={<AuthCallbackPage />} />
                <Route path="/editor" element={<EditorPage />}/>
            </Routes>
        </Router>
    );
}

export default App;