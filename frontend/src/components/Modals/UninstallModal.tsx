import {Button, message, Modal} from 'antd';
import {useState} from "react";
import {HelmModalProps} from "../../types";
import {deleteRelease} from "../../api";
import {useNavigate} from "react-router-dom";

const UninstallModal = ({open, setOpen, release}: HelmModalProps) => {
    const [confirmLoading, setConfirmLoading] = useState(false);
    const [messageApi, contextHolder] = message.useMessage();

    const navigate = useNavigate();

    const showMessage = (
        type: 'success' | 'error' | 'loading',
        content: string,
        duration = 2
    ) => {
        messageApi.open({
            type,
            content,
            duration,
            key: 'delete',
        }).then(() => {
            navigate(0);
        });
    };

    const handleOk = async () => {
        if (!release) return;

        setConfirmLoading(true);

        try {
            const result = await deleteRelease(release?.name || "", release?.namespace || "");

            if (result.code === 200) {
                showMessage('success', 'Uninstalled release.');
            } else if (result.code === 202) {
                showMessage('loading', 'Uninstalling release will continue in the background.');
            } else {
                showMessage('error', 'Uninstall failed.');
            }
        } catch (err) {
            console.error('Error during uninstalling:', err);
            showMessage('error', 'Uninstall error.');
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