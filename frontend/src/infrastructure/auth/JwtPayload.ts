export interface JwtPayload {
    user_id: string;
    email: string;
    name?: string;
    username?: string;
    role?: string;
    roles?: string[];
    is_active?: boolean;
    exp: number;
    iat: number;
}
