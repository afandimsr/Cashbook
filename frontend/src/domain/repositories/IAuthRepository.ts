import type { User, TwoFASetupResponse, LoginResponse } from '../entities/User';

export interface IAuthRepository {
    login(username: string, password: string): Promise<LoginResponse>;
    logout(): Promise<void>;
    getUser(): Promise<User | null>;
    verify2FA(tempToken: string, code: string): Promise<{ token: string; user: User }>;
    verifyBackupCode(tempToken: string, code: string): Promise<{ token: string; user: User }>;
    setup2FA(): Promise<TwoFASetupResponse>;
    verifySetup2FA(code: string): Promise<void>;
    disable2FA(): Promise<void>;
    generateBackupCodes(): Promise<string[]>;
}
