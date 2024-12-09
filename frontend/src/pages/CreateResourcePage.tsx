import {useEffect, useState} from 'react';
import Editor from "../components/Editor/Editor.tsx";
import {createResource} from "../api/k8s/createResource.ts";
import {exampleResourceDefinition} from "../consts/exampleResourceDefinitions.ts";
import {ResourceDetails} from "../types";
import {useLocation} from "react-router-dom";
import {Input} from "antd";

const CreateResourcePage = () => {
    const location = useLocation();
    const {resourceType, namespaces} = location.state || {};

    const [namespace, setNamespace] = useState<string>('');

    useEffect(() => {
        console.log(namespaces);
    }, [namespaces]);

    const namespaceSelector = (
        <Input
            placeholder="Namespace"
            onChange={(e) => setNamespace(e.target.value)}
            value={namespace}
        />
    );

    return (
        <div>
            <Editor name="Create" text={exampleResourceDefinition(resourceType)}
                    endpoint={(data: ResourceDetails) => createResource(resourceType, namespace, data)}
                    namespaceSelector={namespaceSelector}/>
        </div>
    );
};

export default CreateResourcePage;