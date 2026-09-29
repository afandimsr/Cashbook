import type { IUserRepository } from '../../domain/repositories/IUserRepository';
import type { TelegramStatus } from '../../domain/entities/Telegram';

export class GetTelegramStatusUseCase {
    private userRepository: IUserRepository;

    constructor(userRepository: IUserRepository) {
        this.userRepository = userRepository;
    }

    async execute(id: number): Promise<TelegramStatus> {
        if (!id) throw new Error('User ID is required');

        return this.userRepository.getTelegramStatus(id);
    }
}
