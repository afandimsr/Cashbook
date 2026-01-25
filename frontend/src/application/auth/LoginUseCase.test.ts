import { describe, it, expect } from 'vitest';
import { LoginUseCase } from './LoginUseCase';
import { AuthRepository } from '../../infrastructure/auth/AuthRepository';

describe('LoginUseCase', () => {
    const authRepo = new AuthRepository();
    const loginUseCase = new LoginUseCase(authRepo);

    it('should login successfully with valid credentials', async () => {
        const { user, token } = await loginUseCase.execute('admin', 'admin');
        expect(user.username).toBe('admin');
        expect(token).toBeDefined();
    });

    it('should fail with invalid credentials', async () => {
        await expect(loginUseCase.execute('wrong', 'pass')).rejects.toThrow('Invalid credentials');
    });
});
