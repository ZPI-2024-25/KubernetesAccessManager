import { BrowserRouter as Router, Routes, Route } from 'react-router-dom';
import {MainPage} from "./pages/MainPage.tsx";
import Menu from "./components/Menu/Menu.tsx";
import EditorPage from "./pages/EditorPage.tsx";
import ResourcePage from "./pages/ResourcePage.tsx";
import CreateResourcePage from "./pages/CreateResourcePage.tsx";

function App() {
    return (
        <Router>
            <Routes>
                <Route path="/" element={<Menu />}>
                    <Route index element={<MainPage />} />
                    <Route path=":resourceType" element={<ResourcePage />}/>
                    <Route path="/editor" element={<EditorPage />}/>
                    <Route path="/create" element={<CreateResourcePage />}/>
                </Route>
                {/*<Route path="/editor" element={<EditorPage />}/>*/}
                {/*<Route path="/create" element={<CreateResourcePage />}/>*/}
            </Routes>
        </Router>
    );
}

export default App;