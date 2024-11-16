import Protected from "./components/Protected";
import Public from "./components/Public";
import Pods from "./components/Pods.jsx";
import useAuth from "./hooks/useAuth";

function App() {
    // eslint-disable-next-line no-unused-vars
  const [isLogin, token, client] = useAuth();
  return isLogin ? <Protected token={token} client={client}/> : <Public />;
}

export default App;