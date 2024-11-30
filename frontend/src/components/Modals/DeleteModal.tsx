import {Button, Modal} from "antd";
import {ResourceModalProps} from "../../types";
import useShowMessage from "../../hooks/useShowMessage.ts";
import {deleteResource} from "../../api/k8s/deleteResource.ts";

const DeleteConfirmModal = ({open, setOpen, resourceType, resource, removeResource}: ResourceModalProps) => {
    const {showMessage, contextHolder} = useShowMessage();

    const handleOk = async () => {
        if (!resource) return;

        try {
            const result = await deleteResource(resourceType, resource.name as string, resource.namespace as string);

            if (result.code === 200) {
                if (removeResource) {
                    removeResource(resource);
                }
                showMessage({type: 'success', content: 'Deleted resource.', key: 'delete'});
            } else {
                showMessage({type: 'error', content: 'Delete failed.', key: 'delete'});
            }
        } catch (err) {
            console.error('Error during deleting:', err);
            showMessage({type: 'error', content: 'Delete error.', key: 'delete'});
        } finally {
            setOpen(false);
        }
    }

    return (
        <>
            {contextHolder}
            <Modal
                title={resource?.name ? `Delete
                ${resource.name} from
                ${resource.namespace}` : 'Delete'}
                open={open}
                onOk={handleOk}
                footer={
                    [
                        <Button key="back" onClick={() => setOpen(false)}>
                            Cancel
                        </Button>,
                        <Button key="submit" type="primary" danger onClick={handleOk}>
                            Delete
                        </Button>
                    ]
                }
            >
                <p>
                    Do you really want to delete the resource
                    <strong>{` ${resource?.name as string}`}</strong>
                    {resource?.namespace ? (
                        <>
                            {' '}
                            from <strong>{resource.namespace as string}</strong>
                        </>
                    ) : null}?
                </p>

            </Modal>
        </>

    );
};

export default DeleteConfirmModal;
