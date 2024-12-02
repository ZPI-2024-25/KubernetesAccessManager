import React from 'react';
import Editor from "../components/Editor/Editor.tsx";
import {createResource} from "../api/k8s/createResource.ts";
import {exampleResourceDefinition} from "../consts/exampleResourceDefinitions.ts";
import {ResourceDetails} from "../types";
import {useLocation} from "react-router-dom";

const CreateResourcePage: React.FC = () => {
    const location = useLocation();
    const { resourceType } = location.state || {};
    const namespace = "";

    return (
        <div style={{display: "flex"}}>
            <Editor name="Create" text={exampleResourceDefinition(resourceType)}
                    endpoint={(data: ResourceDetails) => createResource(resourceType, namespace, data)}/>
        </div>
    );
};

export default CreateResourcePage;