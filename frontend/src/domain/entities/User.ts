export interface User {
    id: string;
    username: string;
    name: string;
    email: string;
    role: 'ADMIN' | 'USER' | null;
    roles?: string[];
    isActive: boolean;
}
