import type { UserRepository } from '../../domain/repositories/IUserRepository';
import type { User } from '../../domain/entities/User';

export class UpdateUserUseCase {
    private userRepository: UserRepository;

    constructor(userRepository: UserRepository) {
        this.userRepository = userRepository;
    }

    async execute(id: string, user: Partial<User>): Promise<User> {
        return this.userRepository.updateUser(id, user);
    }
}
