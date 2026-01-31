export interface CreateUserRequest {
    id: string;
    name: string;
    email: string;
    password: string;
    roles: string[];
    isActive: boolean;
}
