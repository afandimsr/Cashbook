import type { IUserRepository } from '../../domain/repositories/IUserRepository';
import type { TelegramLinkRow } from '../../domain/entities/Telegram';

export class ListTelegramLinksUseCase {
    private userRepository: IUserRepository;

    constructor(userRepository: IUserRepository) {
        this.userRepository = userRepository;
    }

    async execute(): Promise<TelegramLinkRow[]> {
        return this.userRepository.listTelegramLinks();
    }
}
