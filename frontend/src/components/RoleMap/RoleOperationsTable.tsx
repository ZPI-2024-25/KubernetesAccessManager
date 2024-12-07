import { useState } from 'react';
import { Input, Button, Table, Tag, Select } from 'antd';
import { IoMdClose } from 'react-icons/io';
import { ColumnsType } from 'antd/es/table';
import { Role, RoleOperation } from '../../types';
import { operationsOptions, resourcesOptions } from '../../consts/roleOptions';

interface RoleOperationsTableProps {
    role: Role;
    mapType: 'role' | 'sub';
    operationType: 'permit' | 'deny';
    handleUpdateOperationField: (
        role: Role,
        operation: RoleOperation,
        mapType: 'role' | 'sub',
        field: keyof RoleOperation,
        value: string | string[]
    ) => void;
    handleRemoveOperation: (role: Role, operation: RoleOperation, mapType: 'role' | 'sub') => void;
}

const RoleOperationsTable: React.FC<RoleOperationsTableProps> = ({
                                                                     role,
                                                                     mapType,
                                                                     operationType,
                                                                     handleUpdateOperationField,
                                                                     handleRemoveOperation,
                                                                 }) => {
    const [localNamespaces, setLocalNamespaces] = useState<Record<string, string>>({});

    const getRowKey = (record: RoleOperation) => {
        return `${record.namespace || ''}-${record.resource || ''}`;
    };

    const generateColumns = (): ColumnsType<RoleOperation> => [
        {
            title: '',
            dataIndex: 'remove',
            key: 'remove',
            width: '5%',
            render: (_, record: RoleOperation) => (
                <Button
                    type="text"
                    danger
                    icon={<IoMdClose />}
                    onClick={() => handleRemoveOperation(role, record, mapType)}
                />
            ),
        },
        {
            title: 'Namespace',
            dataIndex: 'namespace',
            key: 'namespace',
            width: '25%',
            render: (namespace: string, record: RoleOperation) => {
                const rowKey = getRowKey(record);
                const localValue = localNamespaces[rowKey] ?? namespace;

                return (
                    <Input
                        value={localValue}
                        placeholder="Namespace"
                        onChange={(e) => {
                            setLocalNamespaces((prev) => ({
                                ...prev,
                                [rowKey]: e.target.value,
                            }));
                        }}
                        onBlur={() => {
                            handleUpdateOperationField(role, record, mapType, 'namespace', localValue);
                            setLocalNamespaces((prev) => {
                                const { [rowKey]: _, ...rest } = prev;
                                return rest;
                            });
                        }}
                    />
                );
            },
        },
        {
            title: 'Resource',
            dataIndex: 'resource',
            key: 'resource',
            width: '25%',
            render: (resource: string, record: RoleOperation) => (
                <Select
                    style={{ width: '100%' }}
                    value={resource}
                    options={resourcesOptions}
                    onChange={(val) => handleUpdateOperationField(role, record, mapType, 'resource', val)}
                    showSearch
                    optionFilterProp="label"
                />
            ),
        },
        {
            title: 'Operations',
            dataIndex: 'operations',
            key: 'operations',
            render: (operations: string[], record: RoleOperation) => (
                <Select
                    mode="tags"
                    style={{ width: '100%' }}
                    value={operations}
                    options={operationsOptions}
                    onChange={(vals: string[]) => {
                        if (vals.includes('*')) {
                            vals = ['*'];
                        } else if (vals.length === 5 && !vals.includes('*')){
                            vals = ['*'];
                        } else {
                            vals = vals.filter((v) => v !== '*');
                        }

                        const allowedValues = operationsOptions.map((o) => o.value);
                        vals = vals.filter((v) => allowedValues.includes(v));

                        handleUpdateOperationField(role, record, mapType, 'operations', vals);
                    }}
                    tagRender={({ label, closable, onClose }) => {
                        const onPreventMouseDown = (event: React.MouseEvent) => {
                            event.preventDefault();
                            event.stopPropagation();
                        };

                        return (
                            <Tag
                                color="green"
                                onMouseDown={onPreventMouseDown}
                                closable={closable}
                                onClose={onClose}
                                style={{ marginInlineEnd: 4 }}
                            >
                                {label}
                            </Tag>
                        );
                    }}
                />
            ),
        },
    ];

    return (
        <Table
            columns={generateColumns()}
            dataSource={role[operationType]}
            rowKey={(record) => getRowKey(record)}
            pagination={false}
            size="small"
        />
    );
};

export default RoleOperationsTable;
