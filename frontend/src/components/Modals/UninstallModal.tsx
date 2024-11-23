import {Button, message, Modal} from 'antd';
import {useState} from "react";
import {HelmRelease} from "../../types";
import {deleteRelease} from "../../api";
import {useNavigate} from "react-router";

const UninstallModal = ({open, setOpen, release}: {
    open: boolean,
    setOpen: (open: boolean) => void,
    release: HelmRelease | undefined
}) => {
    const [confirmLoading, setConfirmLoading] = useState(false);
    const [messageApi, contextHolder] = message.useMessage();

    const navigate = useNavigate();

    const success = () => {
        messageApi.open({
            type: 'success',
            content: 'Uninstalled release.',
            duration: 2,
            key: 'delete',
        })
    }

    const loading = () => {
        messageApi.open({
            type: 'loading',
            content: 'Uninstalling release will continue in the background.',
            duration: 2,
            key: 'delete',
        })
    }

    const error = (message: string) => {
        messageApi.open({
            type: 'error',
            content: message,
            duration: 2,
            key: 'delete',
        })
    }

    const handleOk = async () => {
        setConfirmLoading(true);

        try {
            const result = await deleteRelease(release?.name || "", release?.namespace || "");

            if (result.code === 200) {
                success();
            } else if (result.code === 202) {
                loading();
            } else {
                error('Uninstall failed.');
            }
        } catch (err) {
            console.error('Error during uninstalling:', err);
            error('Uninstall error.');
        } finally {
            setConfirmLoading(false);
            setOpen(false);
            navigate(0)
        }
    };

    return (
        <>
            {contextHolder}
            <Modal
                title={release ? `Uninstall ${release.name} from ${release.namespace}` : 'Uninstall'}
                open={open}
                onOk={release ? handleOk : undefined}
                confirmLoading={confirmLoading}
                onCancel={() => setOpen(false)}
                footer={
                    [
                        <Button key="back" onClick={() => setOpen(false)}>
                            Cancel
                        </Button>,
                        <Button key="submit" type="primary" danger loading={confirmLoading} onClick={handleOk}>
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