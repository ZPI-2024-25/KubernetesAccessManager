import {Button, InputNumber, message, Modal, Space} from 'antd';
import {useEffect, useState} from "react";
import {HelmRelease} from "../../types";
import {rollbackRelease} from "../../api";
import {useNavigate} from "react-router";

const RollbackModal = ({open, setOpen, release}: {
    open: boolean,
    setOpen: (open: boolean) => void,
    release: HelmRelease | undefined
}) => {
    const [confirmLoading, setConfirmLoading] = useState(false);
    const [revision, setRevision] = useState(0);
    const [messageApi, contextHolder] = message.useMessage();

    const navigate = useNavigate();

    useEffect(() => {
        const newRevision = parseInt(release?.revision || "0");
        setRevision(newRevision);
    }, [release]);

    const success = () => {
        messageApi.open({
            type: 'success',
            content: 'Rollback successful.',
            duration: 2,
            key: 'rollback',
        })
    }

    const loading = () => {
        messageApi.open({
            type: 'loading',
            content: 'Rollback will continue in the background.',
            duration: 2,
            key: 'rollback',
        })
    }

    const error = (message: string) => {
        messageApi.open({
            type: 'error',
            content: message,
            duration: 2,
            key: 'rollback',
        })
    }

    const handleOk = async () => {
        setConfirmLoading(true);

        try {
            const result = await rollbackRelease(revision, release?.name || "", release?.namespace || "");

            if ('chart' in result) {
                success();
            } else if ('status' in result && 'message' in result) {
                loading();
            } else {
                error('Rollback failed.');
            }
        } catch (err) {
            console.error('Error during rollback:', err);
            error('Rollback error.');
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
                title={release ? `Rollback ${release.name} from ${release.namespace}` : 'Rollback'}
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
                            Rollback
                        </Button>
                    ]
                }
            >
                <Space>
                    <span><b>Revision:</b></span>
                    <InputNumber min={1} value={revision} onChange={(value) => setRevision(value ? value : 0)} />
                </Space>
            </Modal>
        </>
    );
};

export default RollbackModal;