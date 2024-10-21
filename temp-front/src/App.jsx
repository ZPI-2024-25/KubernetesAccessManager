import Protected from "./components/Protected";
import Public from "./components/Public";
import Pods from "./components/Pods.jsx";
import useAuth from "./hooks/useAuth";

function App() {
    // eslint-disable-next-line no-unused-vars
  const [isLogin, token, _] = useAuth();
  return isLogin ? <Pods token={token}/> : <Public />;
}

export default App;