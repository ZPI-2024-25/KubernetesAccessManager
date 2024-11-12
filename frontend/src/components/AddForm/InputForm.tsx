import { Input } from 'antd';

const { TextArea } = Input;

export const InputForm = () => {
    return (
        <TextArea
            placeholder="Autosize height with minimum and maximum number of lines"
            autoSize={{ minRows: 2, maxRows: 6 }}
        />
    );
};