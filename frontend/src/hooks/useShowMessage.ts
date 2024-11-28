import { message } from 'antd';

type MessageType = 'success' | 'error' | 'loading';

interface ShowMessageProps {
    type: MessageType;
    content: string;
    duration?: number;
    key?: string;
    afterClose?: () => void;
}

const useShowMessage = () => {
    const [messageApi, contextHolder] = message.useMessage();

    const showMessage = ({
                             type,
                             content,
                             duration = 2,
                             key,
                             afterClose,
                         }: ShowMessageProps) => {
        messageApi.open({
            type,
            content,
            duration,
            key,
            onClose: () => {
                if (afterClose) {
                    afterClose();
                }
            },
        });
    };

    return { showMessage, contextHolder };
};

export default useShowMessage;
