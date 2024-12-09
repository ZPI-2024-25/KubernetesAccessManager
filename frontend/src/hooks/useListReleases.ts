import {helmColumns} from "../consts/HelmColumns.ts";
import {HelmDataSourceItem} from "../types";
import {useEffect, useState} from "react";
import {fetchReleases} from "../api";
import {message} from "antd";
import {useAuth} from "../components/AuthProvider/AuthProvider.tsx";
import { hasPermissionInAnyNamespace } from "../functions/authorization.ts";
import { helmResourceLabel } from "../consts/MenuItem.tsx";

export const useListReleases = (namespace: string) => {
    const [dataSource, setDataSource] = useState<HelmDataSourceItem[]>([]);
    const [wasSuccessful, setWasSuccessful] = useState(false);
    const [namespaces, setNamespaces] = useState<string[]>([]);
    const { permissions } = useAuth();

    const extractNamespaces = () => {
        const namespaces = new Set<string>();
        dataSource.forEach((record) => namespaces.add(record.namespace ? record.namespace as string : ''));
        return Array.from(namespaces).filter((namespace) => namespace !== '');
    }

    useEffect(() => {
        if (permissions && !hasPermissionInAnyNamespace(permissions, helmResourceLabel, "l")) {
            setDataSource([]);
            return;
        }
        const fetchData = async () => {
            try {
                const response = await fetchReleases(namespace);

                const dynamicDataSource: HelmDataSourceItem[] = response.map((resource, index) => ({
                    key: index,
                    ...resource,
                }));
                setDataSource(dynamicDataSource);

                setWasSuccessful(true);
            } catch (error) {
                if (error instanceof Error) {
                    console.error('Error fetching releases:', error);
                    message.error(error.message, 4);
                } else {
                    message.error('An unexpected error occurred.');
                }
                setWasSuccessful(false);
            }
        };

        fetchData();
    }, [namespace]);

    useEffect(() => {
        setNamespaces(extractNamespaces());
    }, []);

    return {helmColumns, dataSource, setDataSource, namespaces, wasSuccessful};
}