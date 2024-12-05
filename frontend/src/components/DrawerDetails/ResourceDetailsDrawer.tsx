import React, { useEffect, useState } from "react";
import { Drawer, Spin, Collapse, Typography, Alert } from "antd";
import { HelmDataSourceItem, ResourceDataSourceItem } from "../../types";
import { getResource } from "../../api/k8s/getResource";
import styles from "./ResourceDetailsDrawer.module.css";
import { fetchRelease, fetchReleaseHistory } from "../../api";

const { Paragraph } = Typography;
const { Panel } = Collapse;

interface DrawerDetailsProps {
    visible: boolean;
    record: ResourceDataSourceItem | HelmDataSourceItem | object | null;
    onClose: () => void;
    loading: boolean;
    resourceType: string;
}

const ResourceDetailsDrawer: React.FC<DrawerDetailsProps> = ({
                                                                 visible,
                                                                 record,
                                                                 onClose,
                                                                 resourceType,
                                                             }) => {
    const [resourceDetails, setResourceDetails] = useState<any | null>(null);
    const [fetching, setFetching] = useState(false);
    const [errorMessage, setErrorMessage] = useState<string | null>(null); // State for error messages

    const fetchResourceData = async (
        resourceType: string,
        resourceName: string,
        namespace: string
    ) => {
        setFetching(true);
        setErrorMessage(null); // Reset error message before fetching
        try {
            if (resourceType !== "Helm") {
                const details = await getResource(resourceType, resourceName, namespace);
                setResourceDetails(details);
                console.log(details);
            } else {
                const details = await fetchRelease(resourceName, namespace);
                const history = await fetchReleaseHistory(resourceName, namespace);

                const combinedDetails = [details, history];

                setResourceDetails(combinedDetails);
                console.log(combinedDetails);
            }
        } catch (error: unknown) {
            if (error instanceof Error) {
                const message = `Error fetching resource details: ${error.message}`;
                console.error(message);
                setErrorMessage(message); // Set the error message
            } else {
                const message = "Unknown error occurred while fetching resource details.";
                console.error(message);
                setErrorMessage(message); // Set the error message
            }
            setResourceDetails(null);
        }
        setFetching(false);
    };

    useEffect(() => {
        if (visible && record) {
            if ("name" in record || "resource" in record) {
                const namespace = "resource" in record ? "default" : (record.namespace as string);
                const resourceName = "resource" in record ? (record.resource as string) : (record.name as string);
                fetchResourceData(resourceType, resourceName, namespace);
            }
        }
    }, [visible, record]);

    useEffect(() => {
        onClose();
    }, [resourceType]);

    const renderObject = (obj: any, parentKey = ""): React.ReactNode => {
        if (!obj || typeof obj !== "object") return null;

        return Object.entries(obj).map(([key, value]) => {
            const panelKey = parentKey ? `${parentKey}.${key}` : key;

            if (value && typeof value === "object") {
                return (
                    <Collapse key={panelKey} className={styles.collapse}>
                        <Panel header={key} key={panelKey}>
                            {renderObject(value, panelKey)}
                        </Panel>
                    </Collapse>
                );
            } else {
                return (
                    <div key={panelKey} className={styles.singleDetail}>
                        <div className={styles.key}>{key}:</div>
                        <div className={styles.value}>{value?.toString() || "—"}</div>
                    </div>
                );
            }
        });
    };

    const renderTitle = () => {
        if (record && "name" in record) {
            return `${resourceType}: ${record.name}`;
        }
        return `${resourceType}: Unknown`;
    };

    return (
        <Drawer
            title={renderTitle()}
            placement="right"
            width={600}
            onClose={onClose}
            open={visible}
            mask={false}
        >
            {fetching ? (
                <Spin
                    size="large"
                    style={{ display: "block", textAlign: "center", marginTop: "100px" }}
                />
            ) : errorMessage ? (
                <Alert
                    message="Error"
                    description={errorMessage}
                    type="error"
                    showIcon
                    style={{ marginBottom: "16px" }}
                />
            ) : resourceDetails ? (
                <div className={styles.detailsContainer}>
                    {resourceType === "Helm" ? (
                        <Collapse defaultActiveKey={['release']}>
                            <Panel header="Release" key={"release"}>
                                {renderObject(resourceDetails[0])}
                            </Panel>

                            <Panel header="History" key={"history"}>
                                {resourceDetails[1]?.map((historyItem: any, index: number) => (
                                    <Collapse key={index}>
                                        <Panel header={`Revision ${index + 1}`} key={index + 1}>
                                            {renderObject(historyItem)}
                                        </Panel>
                                    </Collapse>
                                ))}
                            </Panel>
                        </Collapse>
                    ) : (
                        renderObject(resourceDetails)
                    )}
                </div>
            ) : (
                <Paragraph>No data available</Paragraph>
            )}
        </Drawer>
    );
};

export default ResourceDetailsDrawer;
