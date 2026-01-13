import type { User } from '../entities/User';
import type { GetUserUseCaseDTO } from '../../application/user/GetUsersUseCase/GetUserUseCaseDTO';
import type { CreateUserRequest } from '../../application/user/Create/CreateUserRequest';

export interface IUserRepository {
    getUsers(): Promise<GetUserUseCaseDTO[]>;
    createUser(user: Omit<CreateUserRequest, 'id'>): Promise<User>;
    updateUser(id: string, user: Partial<User>): Promise<User>;
    deleteUser(id: string): Promise<void>;
}
