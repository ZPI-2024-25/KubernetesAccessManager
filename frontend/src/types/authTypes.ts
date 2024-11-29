export interface UserStatus {
    permissions: {
        [namespace: string]: {
            [resource: string]: string[];
        };
    };
    user: {
        exp: number;
        preferred_username: string;
        email: string;
    };
}