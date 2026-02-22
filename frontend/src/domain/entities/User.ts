export interface User {
    id: string;
    username: string;
    name: string;
    email: string;
    role: 'ADMIN' | 'USER' | null;
    roles?: string[];
    isActive: boolean;
    totpEnabled?: boolean;
}

export interface TwoFASetupResponse {
    secret: string;
    qr_code: string;
}

export interface LoginResponse {
    token?: string;
    requires_2fa?: boolean;
    temp_token?: string;
}

export interface MFASettings {
    id: number;
    enforce_2fa: boolean;
    updated_by: number;
    updated_at: string;
}
