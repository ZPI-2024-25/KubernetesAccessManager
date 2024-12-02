import {helmColumns} from "../consts/HelmColumns.ts";
import {HelmDataSourceItem, HelmReleaseList} from "../types";
import {useEffect, useState} from "react";
import {fetchReleases} from "../api";
import {message} from "antd";
import {useAuth} from "../components/AuthProvider/AuthProvider.tsx";
import { hasPermissionInAnyNamespace } from "../functions/authorization.ts";

export const useListReleases = () => {
    const [dataSource, setDataSource] = useState<HelmDataSourceItem[]>([]);
    const { userStatus } = useAuth();
    useEffect(() => {
        if (userStatus && !hasPermissionInAnyNamespace(userStatus, "helm", "l")) {
            setDataSource([]);
            return;
        }
        const fetchData = async () => {
            try {
                const response: HelmReleaseList = await fetchReleases('');

                const dynamicDataSource: HelmDataSourceItem[] = response.map((resource, index) => ({
                    key: index,
                    ...resource,
                }));
                setDataSource(dynamicDataSource);
            } catch (error) {
                console.error('Error fetching releases:', error);
                message.error('Failed to fetch releases.', 2);
            }
        };

        fetchData();
    }, []);

    return {helmColumns, dataSource, setDataSource};
}