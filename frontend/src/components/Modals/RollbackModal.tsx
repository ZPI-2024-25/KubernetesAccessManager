import {Button, InputNumber, message, Modal, Space} from 'antd';
import {useEffect, useState} from "react";
import {HelmModalProps} from "../../types";
import {rollbackRelease} from "../../api";
import {useNavigate} from "react-router-dom";

const RollbackModal = ({open, setOpen, release}: HelmModalProps) => {
    const [confirmLoading, setConfirmLoading] = useState(false);
    const [revision, setRevision] = useState(0);
    const [messageApi, contextHolder] = message.useMessage();

    const navigate = useNavigate();

    useEffect(() => {
        const newRevision = parseInt(release?.revision || "0");
        setRevision(newRevision);
    }, [release]);

    const showMessage = (
        type: 'success' | 'error' | 'loading',
        content: string,
        duration = 2
    ) => {
        messageApi.open({
            type,
            content,
            duration,
            key: 'rollback',
        }).then(() => {
            navigate(0);
        });
    };

    const handleOk = async () => {
        if (!release) return;

        setConfirmLoading(true);

        try {
            const result = await rollbackRelease(revision, release?.name || "", release?.namespace || "");

            if ('chart' in result) {
                showMessage('success', 'Rollback successful.');
            } else if ('status' in result && 'message' in result) {
                showMessage('loading', 'Rollback will continue in the background.');
            } else {
                showMessage('error', 'Rollback failed.');
            }
        } catch (err) {
            console.error('Error during rollback:', err);
            showMessage('error', 'Rollback error.');
        } finally {
            setConfirmLoading(false);
            setOpen(false);
        }
    };

    return (
        <>
            {contextHolder}
            <Modal
                title={release ? `Rollback ${release.name} from ${release.namespace}` : 'Rollback'}
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
                            Rollback
                        </Button>
                    ]
                }
            >
                <Space>
                    <span><b>Revision:</b></span>
                    <InputNumber min={1} value={revision} onChange={(value) => setRevision(value ? value : 0)}/>
                </Space>
            </Modal>
        </>
    );
};

export default RollbackModal;