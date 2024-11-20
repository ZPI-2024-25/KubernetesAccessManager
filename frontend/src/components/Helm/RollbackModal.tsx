import {InputNumber, message, Modal, Space} from 'antd';
import {useEffect, useState} from "react";
import {HelmRelease} from "../../types";
import {rollbackRelease} from "../../api";

const RollbackModal = ({open, setOpen, release}: {
    open: boolean,
    setOpen: (open: boolean) => void,
    release: HelmRelease | undefined
}) => {
    const [confirmLoading, setConfirmLoading] = useState(false);
    const [revision, setRevision] = useState(0);

    useEffect(() => {
        const newRevision = parseInt(release?.revision || "0");
        setRevision(newRevision);
    }, [release]);

    const handleOk = async () => {
        setConfirmLoading(true);

        try {
            const result = await rollbackRelease(revision, release?.name || "", release?.namespace || "");

            if ('history_list' in result) {
                message.success('Rollback successful.');
            } else if ('status' in result && 'message' in result) {
                message.loading(`Rollback will continue in the background.`);
            } else {
                message.error('Rollback failed.');
            }
        } catch (error) {
            console.error('Error during rollback:', error);
            message.error('Rollback error.');
        } finally {
            setConfirmLoading(false);
            setOpen(false);
        }
    };

    return (
        <>
            <Modal
                title="Title"
                open={open}
                onOk={release ? handleOk : undefined}
                confirmLoading={confirmLoading}
                onCancel={() => setOpen(false)}
            >
                <Space>
                    <InputNumber min={1} value={revision} onChange={(value) => setRevision(value ? value : 0)} />
                </Space>
                {confirmLoading ? <p>The modal will be closed after two seconds </p> : <p>Rollback</p>}
            </Modal>
        </>
    );
};

export default RollbackModal;