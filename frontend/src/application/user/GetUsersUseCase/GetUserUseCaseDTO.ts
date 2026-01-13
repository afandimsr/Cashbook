export interface GetUserUseCaseDTO {
    id: string;
    name: string;
    email: string;
    username: string;
    role: 'ADMIN' | 'USER' | null;
    is_active: boolean;
}