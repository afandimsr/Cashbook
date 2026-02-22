import type { IAuthRepository } from '../../domain/repositories/IAuthRepository';
import type { User, TwoFASetupResponse, LoginResponse } from '../../domain/entities/User';
import { tokenStorage } from '../storage/tokenStorage';
import { apiClient } from '../apiClient';
import { mapJwtToUser, mapTempJwtToUser } from '../../application/auth/Mapper/AuthMapper';
import { safeDecodeJwt, safeDecodeTempJwt } from './jwt.service';

export class AuthRepository implements IAuthRepository {
    async login(username: string, password: string): Promise<LoginResponse> {
        const response = await apiClient.post<LoginResponse>('/login', {
            email: username,
            password: password,
        });

        // If 2FA is required, return the challenge response
        if (response?.requires_2fa && response?.temp_token) {
            tokenStorage.setToken(response.temp_token);
            return {
                requires_2fa: true,
                temp_token: response.temp_token,
            };
        }

        // Normal login flow
        if (!response?.token) {
            throw new Error('Login failed: token not returned');
        }

        tokenStorage.setToken(response.token);
        return { token: response.token };
    }

    async verify2FA(tempToken: string, code: string): Promise<{ token: string; user: User }> {
        const response = await apiClient.post<LoginResponse>('/2fa/verify', {
            temp_token: tempToken,
            code: code,
        });

        if (!response?.token) {
            throw new Error('2FA verification failed: token not returned');
        }

        const payload = safeDecodeJwt(response.token);
        if (!payload) {
            throw new Error('2FA verification failed: invalid token');
        }
        const user = mapJwtToUser(payload);

        tokenStorage.setToken(response.token);
        return { token: response.token, user };
    }

    async verifyBackupCode(tempToken: string, code: string): Promise<{ token: string; user: User }> {
        const response = await apiClient.post<LoginResponse>('/2fa/backup/verify', {
            temp_token: tempToken,
            code: code,
        });

        if (!response?.token) {
            throw new Error('Backup code verification failed');
        }

        const payload = safeDecodeJwt(response.token);
        if (!payload) {
            throw new Error('Backup code verification failed: invalid token');
        }
        const user = mapJwtToUser(payload);

        tokenStorage.setToken(response.token);
        return { token: response.token, user };
    }

    async setup2FA(): Promise<TwoFASetupResponse> {
        return await apiClient.post<TwoFASetupResponse>('/2fa/setup', {});
    }

    async verifySetup2FA(code: string): Promise<void> {
        await apiClient.post('/2fa/setup/verify', { code });
    }

    async disable2FA(): Promise<void> {
        await apiClient.delete('/2fa/disable');
    }

    async generateBackupCodes(): Promise<string[]> {
        return await apiClient.post<string[]>('/2fa/backup-codes', {});
    }

    async logout(): Promise<void> {
        tokenStorage.clearToken();
    }

    async getUser(): Promise<User | null> {
        const token = tokenStorage.getToken();
        if (!token) return null;

        try {
            // check if tempToken is expired
            const tempToken = tokenStorage.getToken();
            if (tempToken) {
                const payload = safeDecodeTempJwt(tempToken);
                if (!payload) {
                    throw new Error('Login failed: invalid token');
                }


                const user = mapTempJwtToUser(payload);
                if (user.purpose === 'verify' || user.purpose === 'setup') {
                    return null;
                }
            }


            const payload = safeDecodeJwt(token);
            if (!payload) {
                throw new Error('Login failed: invalid token');
            }
            const user = mapJwtToUser(payload);
            return user;
        } catch (err) {
            return null;
        }
    }
}
