import type { IUserRepository } from '../../domain/repositories/IUserRepository';

export class UnlinkTelegramUseCase {
    private userRepository: IUserRepository;

    constructor(userRepository: IUserRepository) {
        this.userRepository = userRepository;
    }

    async execute(id: number): Promise<void> {
        if (!id) throw new Error('User ID is required');

        return this.userRepository.unlinkTelegram(id);
    }
}
