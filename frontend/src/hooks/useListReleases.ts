import {helmColumns} from "../consts/HelmColumns.ts";
import {HelmDataSourceItem} from "../types";
import {useEffect, useState} from "react";
import {fetchReleases} from "../api";
import {message} from "antd";

export const useListReleases = () => {
    const [dataSource, setDataSource] = useState<HelmDataSourceItem[]>([]);

    useEffect(() => {
        const fetchData = async () => {
            try {
                const response = await fetchReleases('');

                const dynamicDataSource: HelmDataSourceItem[] = response.map((resource, index) => ({
                    key: index,
                    ...resource,
                }));
                setDataSource(dynamicDataSource);
            } catch (error) {
                if (error instanceof Error) {
                    console.error('Error fetching releases:', error);
                    message.error(error.message, 4);
                } else {
                    message.error('An unexpected error occurred.');
                }
            }
        };

        fetchData();
    }, []);

    return {helmColumns, dataSource, setDataSource};
}