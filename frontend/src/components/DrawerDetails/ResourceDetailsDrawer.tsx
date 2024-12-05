import React, { useEffect, useState } from "react";
import { Drawer, Spin, Collapse, Typography, Tooltip } from "antd";
import { HelmDataSourceItem, ResourceDataSourceItem } from "../../types";
import { getResource } from "../../api/k8s/getResource";
import styles from "./ResourceDetailsDrawer.module.css";
import { fetchRelease } from "../../api";

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

    const fetchResourceData = async (
        resourceType: string,
        resourceName: string,
        namespace: string
    ) => {
        setFetching(true);
        try {
            if (resourceType !== "Helm") {
                const details = await getResource(resourceType, resourceName, namespace);
                setResourceDetails(details);
                console.log(details);
            } else {
                const details = await fetchRelease(resourceName, namespace);
                setResourceDetails(details);
                console.log(details);
            }
        } catch (error: unknown) {
            if (error instanceof Error) {
                console.error("Error fetching resource details:", error.message);
            } else {
                console.error("Unknown error:", error);
            }
            setResourceDetails(null);
        }
        setFetching(false);
    };

    useEffect(() => {
        if (visible && record) {
            if ("namespace" in record && "name" in record) {
                const namespace = record.namespace as string;
                const resourceName = record.name as string;
                fetchResourceData(resourceType, resourceName, namespace);
            }
        }
    }, [visible, record]);

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
                        <Tooltip title={value?.toString() || "—"}>
                            <div className={styles.value}>{value?.toString() || "—"}</div>
                        </Tooltip>
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
            ) : resourceDetails ? (
                <div className={styles.detailsContainer}>
                    {renderObject(resourceDetails)}
                </div>
            ) : (
                <Paragraph>No data available</Paragraph>
            )}
        </Drawer>
    );
};


export default ResourceDetailsDrawer;
