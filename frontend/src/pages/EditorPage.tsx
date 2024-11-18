import Editor from "../components/Editor/Editor.tsx";
import {createResource, ResourceDetails} from "../api/createResource.ts";
import {updateResource} from "../api/updateResource.ts";
import {getExampleResourceDefinition} from "../consts/exampleResourceDefinitions.ts";

const EditorPage = () => {
    const resourceType = "Deployment";
    const namespace = "";
    const resourceName = "nginx-deployment";

    return (
        <div style={{display: "flex"}}>
            <Editor name="Create" text={getExampleResourceDefinition(resourceType)} endpoint={(data: ResourceDetails) => createResource(resourceType, namespace, data)}/>
            <Editor name="Update" text={""} endpoint={(data: ResourceDetails) => updateResource(resourceType, namespace, resourceName, data)}/>
        </div>

    );
};

export default EditorPage;