import { BrowserRouter as Router, Routes, Route } from 'react-router-dom';
import {MainPage} from "./pages/MainPage.tsx";
import Menu from "./components/Menu/Menu.tsx";
import EditorPage from "./pages/EditorPage.tsx";
import ResourcePage from "./pages/ResourcePage.tsx";
import CreateResourcePage from "./pages/CreateResourcePage.tsx";
import AuthCallbackPage from "./pages/AuthCallbackPage.tsx";
import {initializeAxiosInterceptors} from './axiosConfig.ts';
import {AuthProvider} from "./components/AuthProvider/AuthProvider.tsx";

initializeAxiosInterceptors();

function App() {

    return (
        <Router>
            <AuthProvider>
                <Routes>
                    <Route path="/" element={<Menu />}>
                        <Route index element={<MainPage />} />
                        <Route path=":resourceType" element={<ResourcePage />} />
                        <Route path="/editor" element={<EditorPage />} />
                        <Route path="/create" element={<CreateResourcePage />} />
                    </Route>
                    <Route path="/auth/callback" element={<AuthCallbackPage />} />
                </Routes>
            </AuthProvider>
        </Router>

    );
}

export default App;