import type { IUserRepository } from '../../../domain/repositories/IUserRepository';
import type { GetUserUseCaseDTO } from '../GetUsersUseCase/GetUserUseCaseDTO';

export class GetUsersUseCase {
    private userRepository: IUserRepository;

    constructor(userRepository: IUserRepository) {
        this.userRepository = userRepository;
    }

    async execute(): Promise<GetUserUseCaseDTO[]> {
        return this.userRepository.getUsers();
    }
}
