import type { IUserRepository } from '../../domain/repositories/IUserRepository';
import type { User } from '../../domain/entities/User';

export class UpdateUserUseCase {
    private userRepository: IUserRepository;
    constructor(userRepository: IUserRepository) {
        this.userRepository = userRepository;
    }

    async execute(id: number, user: Partial<User>): Promise<User> {
        return this.userRepository.updateUser(id, user);
    }
}
