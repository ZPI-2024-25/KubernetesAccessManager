import "./MainPage.css"
import Tab from "../components/Table/Tab.tsx";
import { useOutletContext } from 'react-router-dom';



export const MainPage = () => {
    const { currentResourceLabel } = useOutletContext<{ currentResourceLabel: string }>();

    return (
        <>
            <Tab resourceLabel={currentResourceLabel} />
        </>
    );
};