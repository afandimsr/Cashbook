import { config } from '../app/config';
import { tokenStorage } from './storage/tokenStorage';
import { useAuthStore } from '../state/authStore';

const BASE_URL = config.API_URL || 'http://localhost:8080/api/v1';

async function request<T>(path: string, options: RequestInit = {}): Promise<T> {
    const token = tokenStorage.getToken();
    const headers = new Headers(options.headers);

    if (token) {
        headers.set('Authorization', `Bearer ${token}`);
    }

    if (options.body && !(options.body instanceof FormData)) {
        headers.set('Content-Type', 'application/json');
    }

    const response = await fetch(`${BASE_URL}${path}`, {
        ...options,
        headers,
    });

    // Handle 401 Unauthorized globally
    if (response.status === 401) {
        useAuthStore.getState().logout();
        throw new Error('Session expired. Please login again.');
    }

    let data;
    try {
        data = await response.json();
    } catch (e) {
        if (!response.ok) {
            throw new Error('Something went wrong');
        }
    }

    if (!response.ok) {
        throw new Error(data?.message || 'Something went wrong');
    }

    return data.data; // Assuming backend returns { success: true, message: "...", data: ... }
}

export const apiClient = {
    get: <T>(path: string) => request<T>(path, { method: 'GET' }),
    post: <T>(path: string, body: any) => request<T>(path, { method: 'POST', body: JSON.stringify(body) }),
    put: <T>(path: string, body: any) => request<T>(path, { method: 'PUT', body: JSON.stringify(body) }),
    delete: <T>(path: string) => request<T>(path, { method: 'DELETE' }),
};
