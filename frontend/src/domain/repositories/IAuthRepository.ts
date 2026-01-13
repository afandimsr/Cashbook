import type { User } from '../entities/User';

export interface IAuthRepository {
    login(username: string, password: string): Promise<{ token: string; user: User }>;
    logout(): Promise<void>;
    getUser(): Promise<User | null>;
}
