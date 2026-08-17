export interface GetUserUseCaseDTO {
    id: number;
    name: string;
    email: string;
    username: string;
    role: 'ADMIN' | 'USER' | null;
    is_active: boolean;
}