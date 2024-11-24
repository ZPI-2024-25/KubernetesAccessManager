import {HelmRelease} from ".";

export interface HelmModalProps {
    open: boolean;
    setOpen: (open: boolean) => void;
    release?: HelmRelease;
}