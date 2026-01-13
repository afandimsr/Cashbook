import { create } from 'zustand';
import { apiClient } from '../infrastructure/apiClient';
import type { RecurringTransaction } from '../domain/entities/RecurringTransaction';

interface RecurringState {
    recurringTransactions: RecurringTransaction[];
    isLoading: boolean;
    error: string | null;
    fetchRecurring: () => Promise<void>;
    createRecurring: (rt: Partial<RecurringTransaction>) => Promise<void>;
    deleteRecurring: (id: number) => Promise<void>;
    processDue: () => Promise<void>;
}

export const useRecurringStore = create<RecurringState>((set) => ({
    recurringTransactions: [],
    isLoading: false,
    error: null,

    fetchRecurring: async () => {
        set({ isLoading: true, error: null });
        try {
            const data = await apiClient.get<RecurringTransaction[]>('/recurring');
            set({ recurringTransactions: data, isLoading: false });
        } catch (error: any) {
            set({ error: error.message, isLoading: false });
        }
    },

    createRecurring: async (recurring) => {
        set({ isLoading: true, error: null });
        try {
            const data = await apiClient.post<RecurringTransaction>('/recurring', recurring);
            set(state => ({
                recurringTransactions: [...state.recurringTransactions, data],
                isLoading: false
            }));
        } catch (error: any) {
            set({ error: error.message, isLoading: false });
            throw error;
        }
    },

    deleteRecurring: async (id) => {
        set({ isLoading: true, error: null });
        try {
            await apiClient.delete(`/recurring/${id}`);
            set(state => ({
                recurringTransactions: state.recurringTransactions.filter(rt => rt.id !== id),
                isLoading: false
            }));
        } catch (error: any) {
            set({ error: error.message, isLoading: false });
            throw error;
        }
    },

    processDue: async () => {
        set({ isLoading: true, error: null });
        try {
            await apiClient.post('/recurring/process', {});
            set({ isLoading: false });
        } catch (error: any) {
            set({ error: error.message, isLoading: false });
            throw error;
        }
    }
}));
