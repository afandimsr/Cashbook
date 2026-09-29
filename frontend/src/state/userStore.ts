import { create } from 'zustand';
import type { User } from '../domain/entities/User';
import type { TelegramLinkRow, TelegramStatus } from '../domain/entities/Telegram';
import { UserRepositoryImpl } from '../infrastructure/api/user.api';
import { GetUsersUseCase } from '../application/user/GetUsersUseCase/GetUsersUseCase';
import { CreateUserUseCase } from '../application/user/Create/CreateUserUseCase';
import { UpdateUserUseCase } from '../application/user/UpdateUserUseCase';
import { DeleteUserUseCase } from '../application/user/DeleteUserUseCase';
import { ResetPasswordUseCase } from '../application/user/ResetPasswordUseCase';
import { GetTelegramStatusUseCase } from '../application/user/GetTelegramStatusUseCase';
import { GenerateTelegramLinkCodeUseCase } from '../application/user/GenerateTelegramLinkCodeUseCase';
import { UnlinkTelegramUseCase } from '../application/user/UnlinkTelegramUseCase';
import { ListTelegramLinksUseCase } from '../application/user/ListTelegramLinksUseCase';
import type { CreateUserRequest } from '../application/user/Create/CreateUserRequest'

interface UserState {
    users: User[];
    isLoading: boolean;
    error: string | null;
    telegramStatus: TelegramStatus | null;
    telegramLinks: TelegramLinkRow[];
    fetchUsers: () => Promise<void>;
    addUser: (user: Omit<CreateUserRequest, 'id'>) => Promise<void>;
    editUser: (id: number, user: Partial<User>) => Promise<void>;
    removeUser: (id: number) => Promise<void>;
    resetPassword: (id: number, password: string) => Promise<void>;
    fetchTelegramStatus: (id: number) => Promise<void>;
    generateTelegramLinkCode: (id: number) => Promise<{ code: string; expires_at: string }>;
    unlinkTelegram: (id: number) => Promise<void>;
    fetchTelegramLinks: () => Promise<void>;
}

// Dependency Injection
const userRepo = new UserRepositoryImpl();
const getUsersUseCase = new GetUsersUseCase(userRepo);
const createUserUseCase = new CreateUserUseCase(userRepo);
const updateUserUseCase = new UpdateUserUseCase(userRepo);
const deleteUserUseCase = new DeleteUserUseCase(userRepo);
const resetPasswordUseCase = new ResetPasswordUseCase(userRepo);
const getTelegramStatusUseCase = new GetTelegramStatusUseCase(userRepo);
const generateTelegramLinkCodeUseCase = new GenerateTelegramLinkCodeUseCase(userRepo);
const unlinkTelegramUseCase = new UnlinkTelegramUseCase(userRepo);
const listTelegramLinksUseCase = new ListTelegramLinksUseCase(userRepo);

export const useUserStore = create<UserState>((set, get) => ({
    users: [],
    isLoading: false,
    error: null,
    telegramStatus: null,
    telegramLinks: [],

    fetchUsers: async () => {
        set({ isLoading: true, error: null });
        try {
            const users = await getUsersUseCase.execute();

            // map users to match User entity if needed
            const userMapper = users.map(user => {
                return {
                    ...user,
                    isActive: user.is_active
                };
            });

            set({ users: userMapper, isLoading: false });
        } catch (err: any) {
            set({ error: err.message, isLoading: false });
        }
    },

    addUser: async (userData) => {
        set({ isLoading: true, error: null });
        try {
            await createUserUseCase.execute(userData);
            await get().fetchUsers(); // Refresh list
        } catch (err: any) {
            set({ error: err.message, isLoading: false });
            throw err;
        }
    },

    editUser: async (id, userData) => {
        set({ isLoading: true, error: null });
        try {
            await updateUserUseCase.execute(id, userData);
            await get().fetchUsers();
        } catch (err: any) {
            set({ error: err.message, isLoading: false });
            throw err;
        }
    },

    removeUser: async (id) => {
        set({ isLoading: true, error: null });
        try {
            await deleteUserUseCase.execute(id);
            await get().fetchUsers();
        } catch (err: any) {
            set({ error: err.message, isLoading: false });
            throw err;
        }
    },
    
    resetPassword: async (id, password) => {
        set({ isLoading: true, error: null });
        try {
            await resetPasswordUseCase.execute(id, password);
            set({ isLoading: false });
        } catch (err: any) {
            set({ error: err.message, isLoading: false });
            throw err;
        }
    },

    fetchTelegramStatus: async (id) => {
        set({ error: null });
        try {
            const status = await getTelegramStatusUseCase.execute(id);
            set({ telegramStatus: status });
        } catch (err: any) {
            set({ error: err.message, telegramStatus: null });
            throw err;
        }
    },

    generateTelegramLinkCode: async (id) => {
        return generateTelegramLinkCodeUseCase.execute(id);
    },

    unlinkTelegram: async (id) => {
        await unlinkTelegramUseCase.execute(id);
        set({ telegramStatus: { linked: false } });
    },

    fetchTelegramLinks: async () => {
        set({ isLoading: true, error: null });
        try {
            const links = await listTelegramLinksUseCase.execute();
            set({ telegramLinks: links, isLoading: false });
        } catch (err: any) {
            set({ error: err.message, isLoading: false });
        }
    },
}));
