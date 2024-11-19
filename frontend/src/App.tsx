import {BrowserRouter as Router, Route, Routes} from 'react-router-dom';
import {MainPage} from "./pages/MainPage.tsx";
import Menu from "./components/Menu/Menu.tsx";
import EditorPage from "./pages/EditorPage.tsx";
import ResourcePage from "./pages/ResourcePage.tsx";
import AuthCallbackPage from "./pages/AuthCallbackPage.tsx";
import axios from 'axios';
import {AuthProvider} from "./components/AuthProvider/AuthProvider.tsx";

// Konfiguracja interceptorów
axios.interceptors.request.use((config) => {
    const token = localStorage.getItem('access_token');
    if (token) {
        config.headers.Authorization = `Bearer ${token}`;
    }
    return config;
}, (error) => {
    return Promise.reject(error);
});

axios.interceptors.response.use(
    (response) => response,
    (error) => {
        console.error('Błąd w odpowiedzi:', error);
        return Promise.reject(error);
    }
);

function App() {

    return (
        <Router>
            <AuthProvider>
                <Routes>
                    <Route path="/" element={<Menu/>}>
                        <Route index element={<MainPage/>}/>
                        <Route path=":resourceType" element={<ResourcePage/>}/>
                    </Route>
                    <Route path="/auth/callback" element={<AuthCallbackPage/>}/>
                    <Route path="/editor" element={<EditorPage/>}/>
                </Routes>
            </AuthProvider>
        </Router>
    );
}

export default App;