import HelmTab from "../components/Helm/HelmTab.tsx";
import {useState} from "react";
import RollbackModal from "../components/Helm/RollbackModal.tsx";
import {HelmRelease} from "../types";

const HelmPage = () => {
    const [openModal, setOpenModal] = useState(false);
    const [currentRelease, setCurrentRelease] = useState<HelmRelease>();

    const showModal = () => {
        setOpenModal(true);
    };

    return (
        <div>
            <HelmTab showModal={showModal} setCurrent={setCurrentRelease} />
            <RollbackModal open={openModal} setOpen={setOpenModal} release={currentRelease}/>
        </div>
    );
};

export default HelmPage;