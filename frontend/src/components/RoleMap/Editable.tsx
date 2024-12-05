import {Form, FormInstance, Input, InputRef} from "antd";
import {createContext, useContext, useEffect, useRef, useState} from "react";
import {RoleOperation} from "../../types";

const EditableContext = createContext<FormInstance | null>(null);


// eslint-disable-next-line @typescript-eslint/no-unused-vars
export const EditableRow = ({index, ...props}: { index: number }) => {
    const [form] = Form.useForm();
    return (
        <Form form={form} component={false}>
            <EditableContext.Provider value={form}>
                <tr {...props} />
            </EditableContext.Provider>
        </Form>
    );
}

interface EditableCellProps {
    title: React.ReactNode;
    editable: boolean;
    children: React.ReactNode;
    dataIndex: keyof RoleOperation;
    record: RoleOperation;
    handleSave: (record: RoleOperation) => void;
}

export const EditableCell = ({
                                 title,
                                 editable,
                                 children,
                                 dataIndex,
                                 record,
                                 handleSave,
                                 ...restProps
                             }: EditableCellProps) => {
    const [editing, setEditing] = useState(false);
    const inputRef = useRef<InputRef>(null);
    const form = useContext(EditableContext)!;

    useEffect(() => {
        if (editing) {
            inputRef.current?.focus();
        }
    }, [editing]);

    const toggleEdit = () => {
        setEditing(!editing);
        form.setFieldsValue({[dataIndex]: record[dataIndex]});
    };

    const save = async () => {
        try {
            const values = await form.validateFields();

            toggleEdit();
            handleSave({...record, ...values});
        } catch (errInfo) {
            console.log('Save failed:', errInfo);
        }
    };

    let childNode = children;

    if (editable) {
        childNode = editing ? (
            <Form.Item
                style={{margin: 0}}
                name={dataIndex}
                rules={[{required: true, message: `${title} is required.`}]}
            >
                <Input ref={inputRef} onPressEnter={save} onBlur={save}/>
            </Form.Item>
        ) : (
            <div
                className="editable-cell-value-wrap"
                style={{paddingInlineEnd: 24}}
                onClick={toggleEdit}
            >
                {children}
            </div>
        );
    }

    return <td {...restProps}>{childNode}</td>;
};