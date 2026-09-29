import type { User } from '../entities/User';
import type { GetUserUseCaseDTO } from '../../application/user/GetUsersUseCase/GetUserUseCaseDTO';
import type { CreateUserRequest } from '../../application/user/Create/CreateUserRequest';
import type { TelegramLinkCodeResponse, TelegramLinkRow, TelegramStatus } from '../entities/Telegram';

export interface IUserRepository {
    getUsers(): Promise<GetUserUseCaseDTO[]>;
    createUser(user: Omit<CreateUserRequest, 'id'>): Promise<User>;
    updateUser(id: number, user: Partial<User>): Promise<User>;
    deleteUser(id: number): Promise<void>;
    resetPassword(id: number, password: string): Promise<void>;
    getTelegramStatus(id: number): Promise<TelegramStatus>;
    generateTelegramLinkCode(id: number): Promise<TelegramLinkCodeResponse>;
    unlinkTelegram(id: number): Promise<void>;
    listTelegramLinks(): Promise<TelegramLinkRow[]>;
}
