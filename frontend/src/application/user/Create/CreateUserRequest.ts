export interface CreateUserRequest {
    id: string;
    name: string;
    email: string;
    password: string;
    role: 'ADMIN' | 'USER' | null;
    isActive: boolean;
}
