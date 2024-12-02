export interface Permissions {
    [namespace: string]: {
        [resource: string]: string[];
    };
}

export interface UserStatus {
    permissions: Permissions;
    user: {
        exp: number;
        preferred_username: string;
        email: string;
    };
}

export type Operation = "c" | "r" | "u" | "d" | "l";