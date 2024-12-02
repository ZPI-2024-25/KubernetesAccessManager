import {useEffect, useState} from "react";
import {ResourceColumnType, ResourceDataSourceItem} from "../types";
import {fetchResources} from "../api";
import {formatAge} from "../functions/formatAge.ts";
import {useAuth} from "../components/AuthProvider/AuthProvider.tsx";
import { hasPermissionInAnyNamespace } from "../functions/authorization.ts";

export const useListResource = (resourcelabel: string ) => {
    const [columns, setColumns] = useState<ResourceColumnType[]>([]);
    const [dataSource, setDataSource] = useState<ResourceDataSourceItem[]>([]);
    const { userStatus } = useAuth();

    useEffect(() => {
        if (!resourcelabel) return;

        if (userStatus && !hasPermissionInAnyNamespace(userStatus, resourcelabel, "l")) {
            setColumns([]);
            setDataSource([]);
            return;
        }

        const fetchData = async () => {
            const response = await fetchResources(resourcelabel);

            const dynamicColumns: ResourceColumnType[] = response.columns.map((column) => ({
                title: column,
                dataIndex: column,
                key: column,
                width: 150,
                render: (text: React.ReactNode, record: ResourceDataSourceItem): React.ReactNode => {
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
        };

        fetchData();
    }, [resourcelabel]);

    return { columns, dataSource, setDataSource };
};