export interface TempJwtPayload {
    user_id: string;
    email: string;
    purpose: string;
    exp: number;
    iat: number;
}
