import {ResourceDataSourceItem} from "../types";

export function extractCRDname(crd: ResourceDataSourceItem) {
    const resource = (crd.resource as string).toLowerCase() + "s.";
    const group = crd.group as string;

    return resource + group;
}