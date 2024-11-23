import React from "react";
import { Modal } from "antd";

interface DeleteConfirmModalProps {
    visible: boolean;
    resourceName: string;
    namespace: string;
    onConfirm: () => void;
    onCancel: () => void;
}

const DeleteConfirmModal: React.FC<DeleteConfirmModalProps> = ({
                                                                   visible,
                                                                   resourceName,
                                                                   namespace,
                                                                   onConfirm,
                                                                   onCancel,
                                                               }) => {
    return (
        <Modal
            title="Are you sure you want to delete this resource?"
            visible={visible}
            onOk={onConfirm}
            onCancel={onCancel}
            okText="Ok"
            cancelText="Cancel"
        >
            <p>Do you really want to delete the resource <strong>{resourceName}</strong>?</p>
            <p>This action will delete the resource from the <strong>{namespace}</strong> namespace.</p>
        </Modal>
    );
};

export default DeleteConfirmModal;
