import {RoleMap} from "../../types";

const RoleMapForm = ({ data }: { data: RoleMap }) => {
    return (
        <div>
            <h1>Role Map Form</h1>
            <pre>{JSON.stringify(data, null, 2)}</pre>
        </div>
    );
};

export default RoleMapForm;