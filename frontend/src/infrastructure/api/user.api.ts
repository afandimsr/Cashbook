import type { IUserRepository } from '../../domain/repositories/IUserRepository';
import type { User } from '../../domain/entities/User';
import { apiClient } from '../apiClient';

import type { GetUserUseCaseDTO } from '../../application/user/GetUsersUseCase/GetUserUseCaseDTO';
import type { CreateUserRequest } from '../../application/user/Create/CreateUserRequest';
import { toUserApi } from './mappers/user.mapper';

// const MOCK_USERS: User[] = [
//     { id: '1', username: 'admin', name: 'admin', email: 'admin@example.com', role: 'ADMIN', isActive: true },
//     { id: '2', username: 'user1', name: 'user1', email: 'user1@example.com', role: 'USER', isActive: true },
//     { id: '3', username: 'karyawan', name: 'karyawan', email: 'karyawan@example.com', role: 'USER', isActive: false },
// ];

export class UserRepositoryImpl implements IUserRepository {
    // private users: User[] = [...MOCK_USERS];

    async getUsers(): Promise<GetUserUseCaseDTO[]> {
        // await new Promise((resolve) => setTimeout(resolve, 500));
        const response = await apiClient.get<GetUserUseCaseDTO[]>('/users');
        return response;
    }

    async createUser(userData: Omit<CreateUserRequest, 'id'>): Promise<User> {
        try {
            const created = await apiClient.post<User>('/users', userData);
            return created;
        } catch (error: unknown) {
            const message = error instanceof Error ? error.message : String(error);
            throw new Error(`Failed to create user: ${message}`);
        }
    }

    async updateUser(id: string, userData: Partial<User>): Promise<User> {
        await new Promise((resolve) => setTimeout(resolve, 500));

        const payload = toUserApi(userData);

        try {
            const created = await apiClient.put<User>('/users/' + id, payload);
            return created;
        } catch (error: unknown) {
            const message = error instanceof Error ? error.message : String(error);
            throw new Error(`Failed to update user: ${message}`);
        }
    }

    async deleteUser(id: string): Promise<void> {
        await new Promise((resolve) => setTimeout(resolve, 500));

        try {
            await apiClient.delete('/users/' + id);
        } catch (error: unknown) {
            const message = error instanceof Error ? error.message : String(error);
            throw new Error(`Failed to delete user: ${message}`);
        }
    }
}
