import Tab from "../components/Table/Tab.tsx";
import {ReactNode, useState} from "react";
import RollbackModal from "../components/Modals/RollbackModal.tsx";
import {HelmDataSourceItem} from "../types";
import UninstallModal from "../components/Modals/UninstallModal.tsx";
import {useListReleases} from "../hooks/useListReleases.ts";
import {MdOutlineRestore} from "react-icons/md";
import {DeleteOutlined} from "@ant-design/icons";
import {Button} from "antd";

const HelmPage = () => {
    const [openRollbackModal, setOpenRollbackModal] = useState(false);
    const [openUninstallModal, setOpenUninstallModal] = useState(false);
    const [currentRelease, setCurrentRelease] = useState<HelmDataSourceItem>();

    const {helmColumns, dataSource, setDataSource} = useListReleases();
    const columns = helmColumns.concat({
        title: 'Actions',
        dataIndex: "",
        key: 'actions',
        width: 150,
        render: (_: ReactNode, record: HelmDataSourceItem): ReactNode => (
            <div>
                <Button
                    type="link"
                    icon={<MdOutlineRestore/>}
                    onClick={() => handleRollback(record)}
                />
                <Button
                    danger
                    type="link"
                    icon={<DeleteOutlined/>}
                    onClick={() => handleDelete(record)}
                />
            </div>
        ),
    })

    const handleRollback = (record: HelmDataSourceItem) => {
        setCurrentRelease(record);
        setOpenRollbackModal(true);
    }

    const handleDelete = (record: HelmDataSourceItem) => {
        setCurrentRelease(record);
        setOpenUninstallModal(true);
    }

    const removeRelease = (release: HelmDataSourceItem) => {
        setDataSource(dataSource.filter(item => item.key !== release.key));
    }

    return (
        <div>
            <Tab columns={columns} dataSource={dataSource}/>
            <RollbackModal open={openRollbackModal} setOpen={setOpenRollbackModal} release={currentRelease}/>
            <UninstallModal open={openUninstallModal} setOpen={setOpenUninstallModal} release={currentRelease} removeRelease={removeRelease}/>
        </div>
    );
};

export default HelmPage;