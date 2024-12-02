import {Button, InputNumber, Modal, Space} from 'antd';
import {useEffect, useState} from "react";
import {HelmModalProps} from "../../types";
import {rollbackRelease} from "../../api";
import {useNavigate} from "react-router-dom";
import useShowMessage from "../../hooks/useShowMessage.ts";

const RollbackModal = ({open, setOpen, release}: HelmModalProps) => {
    const [confirmLoading, setConfirmLoading] = useState(false);
    const [revision, setRevision] = useState(0);

    const navigate = useNavigate();
    const {showMessage, contextHolder} = useShowMessage();

    useEffect(() => {
        const newRevision = parseInt(release?.revision || "0");
        setRevision(newRevision);
    }, [release]);

    const handleOk = async () => {
        if (!release) return;

        setConfirmLoading(true);

        try {
            const result = await rollbackRelease(revision, release?.name || "", release?.namespace || "");

            if ('chart' in result) {
                showMessage({
                    type: 'success',
                    content: 'Rollback successful.',
                    key: 'rollback',
                    afterClose: () => navigate(0)
                });
            } else if ('status' in result && 'message' in result) {
                showMessage({
                    type: 'loading',
                    content: 'Rollback will continue in the background.',
                    key: 'rollback',
                    afterClose: () => navigate(0)
                });
            } else {
                showMessage({
                    type: 'error',
                    content: 'Rollback failed.',
                    key: 'rollback',
                    afterClose: () => navigate(0)
                });
            }
        } catch (err) {
            if (err instanceof Error) {
                console.error('Error rollbacking release:', err);
                showMessage({
                    type: 'error',
                    content: err.message,
                    key: 'rollback',
                    duration: 4,
                    afterClose: () => navigate(0)
                });
            } else {
                console.error('An unexpected error occurred:', err);
                showMessage({
                    type: 'error',
                    content: 'Unexpected error',
                    key: 'rollback',
                    duration: 4,
                    afterClose: () => navigate(0)
                });
            }
        } finally {
            setConfirmLoading(false);
            setOpen(false);
        }
    };

    return (
        <>
            {contextHolder}
            <Modal
                title={release ? `Rollback ${release.name} in ${release.namespace}` : 'Rollback'}
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