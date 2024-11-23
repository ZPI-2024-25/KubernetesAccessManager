import HelmTab from "../components/Helm/HelmTab.tsx";
import {useState} from "react";
import RollbackModal from "../components/Modals/RollbackModal.tsx";
import {HelmRelease} from "../types";
import UninstallModal from "../components/Modals/UninstallModal.tsx";

const HelmPage = () => {
    const [openRollbackModal, setOpenRollbackModal] = useState(false);
    const [openUninstallModal, setOpenUninstallModal] = useState(false);
    const [currentRelease, setCurrentRelease] = useState<HelmRelease>();

    const showRollbackModal = () => {
        setOpenRollbackModal(true);
    };

    const showUninstallModal = () => {
        setOpenUninstallModal(true)
    }

    return (
        <div>
            <HelmTab showRollbackModal={showRollbackModal} showUninstallModal={showUninstallModal}
                     setCurrent={setCurrentRelease}/>
            <RollbackModal open={openRollbackModal} setOpen={setOpenRollbackModal} release={currentRelease}/>
            <UninstallModal open={openUninstallModal} setOpen={setOpenUninstallModal} release={currentRelease}/>
        </div>
    );
};

export default HelmPage;