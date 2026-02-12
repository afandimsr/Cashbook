import type { IUserRepository } from '../../domain/repositories/IUserRepository';

export class ResetPasswordUseCase {
    private userRepository: IUserRepository;

    constructor(userRepository: IUserRepository) {
        this.userRepository = userRepository;
    }

    async execute(id: string, password: string): Promise<void> {
        if (!id) throw new Error('User ID is required');
        if (!password) throw new Error('Password is required');

        return this.userRepository.resetPassword(id, password);
    }
}
