import type { IUserRepository } from '../../domain/repositories/IUserRepository';

export class DeleteUserUseCase {
    private userRepository: IUserRepository;

    constructor(userRepository: IUserRepository) {
        this.userRepository = userRepository;
    }

    async execute(id: number): Promise<void> {
        return this.userRepository.deleteUser(id);
    }
}
