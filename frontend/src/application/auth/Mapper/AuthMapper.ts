import type { JwtPayload } from '../../../infrastructure/auth/JwtPayload';
import type { User } from '../../../domain/entities/User';

export function mapJwtToUser(payload: JwtPayload): User {
    return {
        id: payload.user_id,
        email: payload.email,
        name: payload.name || '',
        // prefer explicit roles array, fallback to single role claim
        role: (payload.roles && payload.roles.length > 0) ? (payload.roles[0] as any) : (payload.role as any),
        username: payload.username || '',
        roles: payload.roles,
        isActive: payload.is_active || false,
    };
}
