import React from 'react';
import Editor from "../components/Editor/Editor.tsx";
import {createResource} from "../api/createResource.ts";
import {getExampleResourceDefinition} from "../consts/exampleResourceDefinitions.ts";
import {ResourceDetails} from "../types/ResourceDetails.ts";

const CreateResourcePage: React.FC = () => {
    const resourceType = "Deployment";
    const namespace = "";

    return (
        <div style={{display: "flex"}}>
            <Editor name="Create" text={getExampleResourceDefinition(resourceType)}
                    endpoint={(data: ResourceDetails) => createResource(resourceType, namespace, data)}/>
        </div>
    );
};

export default CreateResourcePage;