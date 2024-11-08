import {MenuProps} from "antd";
import React from "react";

export interface MenuItem {
    key: string;
    label:string;
    resourcelabel:string;
    icon?: React.ReactNode;
    children?: MenuItem[];
}

export type MenuItemType = Required<MenuProps>['items'][number];