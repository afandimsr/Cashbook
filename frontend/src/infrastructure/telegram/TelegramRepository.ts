import type { TelegramLinkCodeResponse, TelegramStatus } from '../../domain/entities/Telegram';
import { apiClient } from '../apiClient';

export class TelegramRepository {
    async generateLinkCode(): Promise<TelegramLinkCodeResponse> {
        return await apiClient.post<TelegramLinkCodeResponse>('/telegram/link-code', {});
    }

    async getStatus(): Promise<TelegramStatus> {
        return await apiClient.get<TelegramStatus>('/telegram/status');
    }

    async unlink(): Promise<void> {
        await apiClient.delete('/telegram/link');
    }
}
