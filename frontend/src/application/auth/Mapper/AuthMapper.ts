import type { JwtPayload } from '../../../infrastructure/auth/JwtPayload';
import type { TempJwtPayload } from '../../../infrastructure/auth/TempJwtPayload';
import type { User } from '../../../domain/entities/User';
import type { TempUser } from '../../../domain/entities/TempUser';

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

export function mapTempJwtToUser(payload: TempJwtPayload): TempUser {
    return {
        user_id: payload.user_id,
        email: payload.email,
        purpose: payload.purpose,
    };
}
