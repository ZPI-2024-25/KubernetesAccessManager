import {useParams} from "react-router";
// import {useState} from "react";
import Tab from "../components/Table/Tab.tsx";

const ResourcePage = () => {
    const { resourceType } = useParams();
    // const [namespace, setNamespace] = useState<string>('');

    return (
        <div>
            { resourceType ? <Tab resourceLabel={resourceType}/> : "Resource not found" }
        </div>
    );
};

export default ResourcePage;