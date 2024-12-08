import RoleMapForm from "../components/RoleMap/RoleMapForm.tsx";
import {useLocation} from "react-router-dom";
import styles from "./RolesPage.module.css";

const EditRolesPage = () => {
    const {roleMap} = useLocation().state;

    // const navigate = useNavigate();

    return (
        <div className={styles.container}>
            <RoleMapForm data={roleMap}/>
        </div>
    );
};

export default EditRolesPage;