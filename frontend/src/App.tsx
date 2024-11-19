import {BrowserRouter as Router, Route, Routes} from 'react-router-dom';
import {MainPage} from "./pages/MainPage.tsx";
import Menu from "./components/Menu/Menu.tsx";
import EditorPage from "./pages/EditorPage.tsx";
import ResourcePage from "./pages/ResourcePage.tsx";
import AuthCallbackPage from "./pages/AuthCallbackPage.tsx";
import axios from 'axios';
import {AuthProvider} from "./components/AuthProvider/AuthProvider.tsx";

axios.interceptors.request.use(
    (config) => {
        try {
            const token = localStorage.getItem('access_token');
            if (token) {
                config.headers.Authorization = `Bearer ${token}`;
            }
        } catch (error) {
            console.warn('Cant attach bearer token:', error);
        }
        return config;
    },
    (error) => {
        console.warn('Axios interceptor error:', error);
        return Promise.resolve(error.config || {});
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