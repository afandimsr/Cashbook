export interface TelegramLinkCodeResponse {
    code: string;
    expires_at: string;
}

export interface TelegramStatus {
    linked: boolean;
    telegram_username?: string;
    linked_at?: string;
}

export interface TelegramLinkRow {
    id: number;
    user_id: number;
    telegram_chat_id: number;
    telegram_username?: string;
    is_active: boolean;
    linked_at: string;
}
