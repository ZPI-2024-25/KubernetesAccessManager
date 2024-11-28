import {HelmDataSourceItem, ResourceDataSourceItem} from ".";

export interface ResourceModalProps {
    open: boolean;
    setOpen: (open: boolean) => void;
    resourceType: string;
    resource?: ResourceDataSourceItem;
    removeResource?: (record: ResourceDataSourceItem) => void;
}

export interface HelmModalProps {
    open: boolean;
    setOpen: (open: boolean) => void;
    release?: HelmDataSourceItem;
    removeRelease?: (release: HelmDataSourceItem) => void;
}