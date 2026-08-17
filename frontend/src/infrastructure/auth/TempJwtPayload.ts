export interface TempJwtPayload {
    user_id: number;
    email: string;
    purpose: string;
    exp: number;
    iat: number;
}
