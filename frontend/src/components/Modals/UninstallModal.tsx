import {Button, Modal} from 'antd';
import {useState} from "react";
import {HelmModalProps} from "../../types";
import {deleteRelease} from "../../api";
import {useNavigate} from "react-router-dom";
import useShowMessage from "../../hooks/useShowMessage.ts";

const UninstallModal = ({open, setOpen, release, removeRelease}: HelmModalProps) => {
    const [confirmLoading, setConfirmLoading] = useState(false);

    const navigate = useNavigate();
    const {showMessage, contextHolder} = useShowMessage();

    const handleOk = async () => {
        if (!release) return;

        setConfirmLoading(true);

        try {
            const result = await deleteRelease(release?.name || "", release?.namespace || "");

            if (result.code === 200) {
                if (removeRelease) {
                    removeRelease(release);
                }
                showMessage({type: 'success', content: 'Uninstalled release.', key: 'uninstall'});
            } else if (result.code === 202) {
                showMessage({type: 'loading', content: 'Uninstall will continue in the background.', key: 'uninstall', afterClose: () => navigate(0)});
            } else {
                showMessage({type: 'error', content: 'Uninstall failed.', key: 'uninstall', afterClose: () => navigate(0)});
            }
        } catch (err) {
            console.error('Error during uninstalling:', err);
            showMessage({type: 'error', content: 'Uninstall error.', key: 'uninstall' , afterClose: () => navigate(0)});
        } finally {
            setConfirmLoading(false);
            setOpen(false);
        }
    };

    return (
        <>
            {contextHolder}
            <Modal
                title={release ? `Uninstall ${release.name} from ${release.namespace}` : 'Uninstall'}
                open={open}
                confirmLoading={confirmLoading}
                onCancel={() => setOpen(false)}
                footer={
                    [
                        <Button key="back" onClick={() => setOpen(false)}>
                            Cancel
                        </Button>,
                        <Button key="submit" type="primary" danger loading={confirmLoading} onClick={handleOk}
                                disabled={!release}>
                            Uninstall
                        </Button>
                    ]
                }
            >
                <p>Do you really want to uninstall the
                    release <strong>{release?.name}</strong> from <strong>{release?.namespace}</strong>?
                </p>
            </Modal>
        </>
    );
};

export default UninstallModal;