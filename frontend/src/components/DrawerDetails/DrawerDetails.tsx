import React from 'react';
import { Drawer, Spin, Typography } from 'antd';

const { Title, Paragraph } = Typography;

interface DrawerDetailsProps {
    visible: boolean;
    record: object | null;
    onClose: () => void;
    loading: boolean;
}

const DrawerDetails: React.FC<DrawerDetailsProps> = ({ visible, record, onClose, loading }) => {
    return (
        <Drawer
            title="Details"
            placement="right"
            width={600}
            onClose={onClose}
            open={visible}
        >
            {loading ? (
                <Spin size="large" style={{ display: "block", textAlign: "center", marginTop: "100px" }} />
            ) : record ? (
                Array.isArray(record) ? (
                    // Отображаем список для истории Helm-релиза
                    <div>
                        <Title level={4}>Release History</Title>
                        {record.map((item, index) => (
                            <div key={index} style={{ marginBottom: "16px" }}>
                                <strong>Version:</strong> {item.version} <br />
                                <strong>Status:</strong> {item.status} <br />
                                <strong>Date:</strong> {item.date}
                            </div>
                        ))}
                    </div>
                ) : (
                    <div>
                        <Title level={4}>Resource Details</Title>
                        {Object.entries(record).map(([key, value]) => (
                            <div key={key} style={{ marginBottom: "12px" }}>
                                <strong>{key}:</strong> {value?.toString() || 'N/A'}
                            </div>
                        ))}
                    </div>
                )
            ) : (
                <Paragraph>No data available</Paragraph>
            )}
        </Drawer>
    );
};

export default DrawerDetails;
