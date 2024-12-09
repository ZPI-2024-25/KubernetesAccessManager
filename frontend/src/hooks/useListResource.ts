import {ReactNode, useEffect, useState} from "react";
import {ResourceColumnType, ResourceDataSourceItem} from "../types";
import {fetchResources} from "../api";
import {formatAge} from "../functions/formatAge.ts";
import {message} from "antd";
import {useAuth} from "../components/AuthProvider/AuthProvider.tsx";
import { hasPermissionInAnyNamespace } from "../functions/authorization.ts";

export const useListResource = (resourcelabel: string, namespace: string ) => {
    const [columns, setColumns] = useState<ResourceColumnType[]>([]);
    const [dataSource, setDataSource] = useState<ResourceDataSourceItem[]>([]);
    const [wasSuccessful, setWasSuccessful] = useState(false);
    const [namespaces, setNamespaces] = useState<string[]>([]);
    const { permissions } = useAuth();

    const extractNamespaces = () => {
        const namespaces = new Set<string>();
        dataSource.forEach((record) => namespaces.add(record.namespace ? record.namespace as string : ''));
        return Array.from(namespaces).filter((namespace) => namespace !== '');
    }

    useEffect(() => {
        if (!resourcelabel) return;

        if (permissions && !hasPermissionInAnyNamespace(permissions, resourcelabel, "l")) {
            setColumns([]);
            setDataSource([]);
            return;
        }

        const fetchData = async () => {
            try {
                const response = await fetchResources(resourcelabel, namespace);

                const dynamicColumns = response.columns.map((column) => ({
                    title: column,
                    dataIndex: column,
                    key: column,
                    width: 150,
                    render: (text: ReactNode, record: ResourceDataSourceItem): ReactNode => {
                        if (column.toLowerCase() == ('age')) {
                            return formatAge(record[column] as string);
                        }
                        return text;
                    },
                }));

                setColumns(dynamicColumns);
                setDataSource(
                    response.resource_list.map((resource, index) => ({
                        key: index,
                        ...resource,
                    }))
                );

                setWasSuccessful(true);
            } catch (error) {
                if (error instanceof Error) {
                    console.error('Error fetching resources:', error);
                    message.error(error.message, 4);
                } else {
                    message.error('An unexpected error occurred.');
                }
                setWasSuccessful(false);
            }
        };

        fetchData();
    }, [resourcelabel, namespace]);

    useEffect(() => {
        setNamespaces(extractNamespaces());
    }, [resourcelabel]);

    return { columns, dataSource, namespaces, setDataSource, wasSuccessful };
};