import type { IUserRepository } from '../../domain/repositories/IUserRepository';
import type { TelegramLinkCodeResponse } from '../../domain/entities/Telegram';

export class GenerateTelegramLinkCodeUseCase {
    private userRepository: IUserRepository;

    constructor(userRepository: IUserRepository) {
        this.userRepository = userRepository;
    }

    async execute(id: number): Promise<TelegramLinkCodeResponse> {
        if (!id) throw new Error('User ID is required');

        return this.userRepository.generateTelegramLinkCode(id);
    }
}
