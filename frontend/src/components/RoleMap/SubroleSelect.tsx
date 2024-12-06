import React from "react";
import {Select, Tag} from "antd";
import {Role} from "../../types";

interface SubroleOption {
    label: string;
    value: string;
}

interface SubroleSelectProps {
    role: Role;
    mapType: "role" | "sub";
    roleMap: Role[];
    setRoleMap: (val: Role[]) => void;
    subroleMap: Role[];
    setSubroleMap: (val: Role[]) => void;
    subroles: SubroleOption[];
}

const SubroleSelect: React.FC<SubroleSelectProps> = ({
                                                         role,
                                                         mapType,
                                                         roleMap,
                                                         setRoleMap,
                                                         subroleMap,
                                                         setSubroleMap,
                                                         subroles
                                                     }) => {
    const handleChange = (value: string[]) => {
        if (mapType === "role") {
            const roleIndex = roleMap.findIndex(r => r.name === role.name);
            if (roleIndex === -1) return;

            const updatedRoleMap = [...roleMap];
            const updatedRole = {...role};
            updatedRole.subroles = value;
            updatedRoleMap[roleIndex] = updatedRole;

            setRoleMap(updatedRoleMap);
        } else {
            const roleIndex = subroleMap.findIndex(r => r.name === role.name);
            if (roleIndex === -1) return;

            const updatedSubRoleMap = [...subroleMap];
            const updatedRole = {...role};
            updatedRole.subroles = value;
            updatedSubRoleMap[roleIndex] = updatedRole;

            setSubroleMap(updatedSubRoleMap);
        }
    };

    // eslint-disable-next-line
    const tagRender = ({label, value, closable, onClose}: any) => {
        const isValid = subroles.some((opt) => opt.value === value);

        const onPreventMouseDown = (event: React.MouseEvent) => {
            event.preventDefault();
            event.stopPropagation();
        }

        return (
            <Tag
                color={isValid ? "blue" : "red"}
                onMouseDown={onPreventMouseDown}
                closable={closable}
                onClose={onClose}
                style={{marginInlineEnd: 4}}
            >
                {label}
            </Tag>
        );
    };

    return (
        <Select
            mode="tags"
            allowClear
            style={{width: "100%"}}
            placeholder="Select subroles"
            value={role.subroles ?? []}
            onChange={handleChange}
            options={subroles}
            tagRender={tagRender}
        />
    );
};

export default SubroleSelect;
