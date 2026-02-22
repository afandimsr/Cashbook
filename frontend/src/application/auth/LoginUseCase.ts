import type { IAuthRepository } from '../../domain/repositories/IAuthRepository';
import type { LoginResponse } from '../../domain/entities/User';

export class LoginUseCase {
    private authRepository: IAuthRepository;

    constructor(authRepository: IAuthRepository) {
        this.authRepository = authRepository;
    }

    async execute(username: string, password: string): Promise<LoginResponse> {
        if (!username || !password) {
            throw new Error('Username and password are required');
        }
        return this.authRepository.login(username, password);
    }
}
