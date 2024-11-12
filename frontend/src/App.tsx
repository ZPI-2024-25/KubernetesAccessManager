import { BrowserRouter as Router, Routes, Route } from 'react-router-dom';
import LeftMenu from "./components/Menu/Menu.tsx";
import {MainPage} from "./pages/MainPage.tsx";
import {AddResourcePage} from "./pages/AddResourcePage.tsx";

const App: React.FC = () => {
    return (
        <Router>
            <Routes>
                <Route path="/kam/*" element={<LeftMenu />}>
                    <Route index element={<MainPage />} />
                    <Route path="add" element={<AddResourcePage />} />
                </Route>
            </Routes>
        </Router>
    );
};

export default App;
