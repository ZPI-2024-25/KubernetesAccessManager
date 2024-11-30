import  { useEffect, useState } from "react";
import { useLocation } from "react-router-dom";
import { message } from "antd";
import Editor from "../components/Editor/Editor.tsx";
import {updateResource} from "../api/k8s/updateResource.ts";
import {getResource} from "../api/k8s/getResource.ts";
import {ResourceDetails} from "../types";
import {stringifyYaml} from "../functions/jsonYamlFunctions.ts";

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
                setResourceData(stringifyYaml(data));
            } catch (error) {
                if (error instanceof Error) {
                    console.error('Error getting resource:', error);
                    message.error(error.message, 4);
                } else {
                    message.error('An unexpected error occurred.');
                }
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