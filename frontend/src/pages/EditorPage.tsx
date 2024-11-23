import  { useEffect, useState } from "react";
import { useLocation } from "react-router-dom";
import { message } from "antd";
import Editor from "../components/Editor/Editor.tsx";
import {updateResource} from "../api/k8s/updateResource.ts";
import {getResource} from "../api/k8s/getResource.ts";
import {ResourceDetails} from "../types/ResourceDetails.ts";

const EditorPage = () => {
    const location = useLocation();
    const { resourceType, namespace, resourceName } = location.state || {};
    const [resourceData, setResourceData] = useState<string>("");

    useEffect(() => {
        if (!resourceType || !resourceName) {
            message.error("Missing resource parameters");
            return;
        }

        const fetchData = async () => {
            try {
                const data = await getResource(resourceType, resourceName, namespace);
                setResourceData(data);
            } catch (error) {
                console.error("Failed to fetch resource details:", error);
                message.error("Failed to load resource data");
            }
        };

        fetchData();
    }, [resourceType, namespace, resourceName]);

    return (
        <div style={{ display: "flex" }}>
            <Editor
                name={`Edit ${resourceName}`}
                text={resourceData}
                endpoint={(data: ResourceDetails) =>
                    updateResource(resourceType, namespace, resourceName, data)
                }
            />
        </div>
    );
};

export default EditorPage;