import type { User } from "../../../domain/entities/User";

export function toUserApi(user: Partial<User>) {
    const { isActive, ...rest } = user;

    return {
        ...rest,
        ...(isActive !== undefined && { is_active: isActive }),
    };
}

export function fromUserApi(data: any): User {
    return {
        id: data.id,
        name: data.name,
        email: data.email,
        isActive: data.is_active,
        roles: data.roles,
        username: data.username,
        role: data.role,
    };
}
