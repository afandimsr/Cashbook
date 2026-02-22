import { create } from 'zustand';
import type { User } from '../domain/entities/User';
import type { TempUser } from '../domain/entities/TempUser';
import { AuthRepository } from '../infrastructure/auth/AuthRepository';
import { LoginUseCase } from '../application/auth/LoginUseCase';
import { tokenStorage } from '../infrastructure/storage/tokenStorage';
import { mapJwtToUser } from '../application/auth/Mapper/AuthMapper';
import { safeDecodeJwt } from '../infrastructure/auth/jwt.service';
import { safeDecodeTempJwt } from '../infrastructure/auth/jwt.service';
import { mapTempJwtToUser } from '../application/auth/Mapper/AuthMapper';

interface AuthState {
    user: User | null;
    token: string | null;
    isAuthenticated: boolean;
    isInitialized: boolean;
    isLoading: boolean;
    error: string | null;
    // 2FA state
    requires2FA: boolean;
    tempToken: string | null;
    tempUser: TempUser | null;
    login: (username: string, password: string) => Promise<void>;
    verify2FA: (code: string) => Promise<void>;
    verifyBackupCode: (code: string) => Promise<void>;
    logout: () => void;
    initializeAuth: () => Promise<void>;
    handleOAuthSuccess: (token: string) => Promise<void>;
    clear2FAState: () => void;
}

// Dependency Injection (simple version)
const authRepository = new AuthRepository();
const loginUseCase = new LoginUseCase(authRepository);

export const useAuthStore = create<AuthState>((set, get) => ({
    user: null,
    token: null,
    isAuthenticated: false,
    isInitialized: false,
    isLoading: false,
    error: null,
    requires2FA: false,
    tempToken: null,
    tempUser: null,

    initializeAuth: async () => {
        const token = tokenStorage.getToken();
        if (token) {
            try {
                const user = await authRepository.getUser();
                if (user) {
                    set({ user, token, isAuthenticated: true, isInitialized: true });
                } else {
                    tokenStorage.clearToken();
                    set({ user: null, token: null, isAuthenticated: false, isInitialized: true });
                }
            } catch (err) {
                tokenStorage.clearToken();
                set({ user: null, token: null, isAuthenticated: false, isInitialized: true });
            }
        } else {
            set({ isInitialized: true });
        }
    },

    login: async (username, password) => {
        set({ isLoading: true, error: null, requires2FA: false, tempToken: null });
        try {
            const result = await loginUseCase.execute(username, password);
            if (result.requires_2fa && result.temp_token) {
                // 2FA is required — store temp token and signal UI
                const payload = safeDecodeTempJwt(result.temp_token);
                if (!payload) {
                    throw new Error('Login failed: invalid token');
                }
                const tempUser = mapTempJwtToUser(payload);

                set({
                    tempUser: tempUser,
                    requires2FA: true,
                    tempToken: result.temp_token,
                    isLoading: false,

                });
                return;
            }

            // Normal login — decode token to get user
            if (result.token) {
                const payload = safeDecodeJwt(result.token);
                if (!payload) {
                    throw new Error('Login failed: invalid token');
                }
                const user = mapJwtToUser(payload);
                set({ user, token: result.token, isAuthenticated: true, isLoading: false });
            }
        } catch (err: any) {
            set({ error: err.message || 'Login failed', isLoading: false });
            throw err;
        }
    },

    verify2FA: async (code: string) => {
        const { tempToken } = get();
        if (!tempToken) {
            set({ error: 'No 2FA session found' });
            return;
        }

        set({ isLoading: true, error: null });
        try {
            const { user, token } = await authRepository.verify2FA(tempToken, code);
            set({
                user,
                token,
                isAuthenticated: true,
                isLoading: false,
                requires2FA: false,
                tempToken: null,
            });
        } catch (err: any) {
            set({ error: err.message || '2FA verification failed', isLoading: false });
            throw err;
        }
    },

    verifyBackupCode: async (code: string) => {
        const { tempToken } = get();
        if (!tempToken) {
            set({ error: 'No 2FA session found' });
            return;
        }

        set({ isLoading: true, error: null });
        try {
            const { user, token } = await authRepository.verifyBackupCode(tempToken, code);
            set({
                user,
                token,
                isAuthenticated: true,
                isLoading: false,
                requires2FA: false,
                tempToken: null,
            });
        } catch (err: any) {
            set({ error: err.message || 'Backup code verification failed', isLoading: false });
            throw err;
        }
    },

    logout: () => {
        authRepository.logout();
        set({ user: null, token: null, isAuthenticated: false, requires2FA: false, tempToken: null });
    },

    handleOAuthSuccess: async (token: string) => {
        set({ isLoading: true, error: null });
        try {
            tokenStorage.setToken(token);
            const user = await authRepository.getUser();
            if (user) {
                set({ user, token, isAuthenticated: true, isLoading: false });
            } else {
                throw new Error('Failed to fetch user data after OAuth');
            }
        } catch (err: any) {
            tokenStorage.clearToken();
            set({ error: err.message || 'OAuth login failed', isLoading: false });
            throw err;
        }
    },

    clear2FAState: () => {
        set({ requires2FA: false, tempToken: null, error: null });
    },
}));

// Initialize auth state on app load
useAuthStore.getState().initializeAuth();
