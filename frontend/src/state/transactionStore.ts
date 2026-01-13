import { create } from 'zustand';
import type { Transaction, DashboardSummary } from '../domain/entities/Transaction';
import { apiClient } from '../infrastructure/apiClient';

interface TransactionState {
    transactions: Transaction[];
    summary: DashboardSummary | null;
    isLoading: boolean;
    error: string | null;
    fetchTransactions: (page?: number, limit?: number, search?: string) => Promise<void>;
    fetchSummary: () => Promise<void>;
    addTransaction: (transaction: Omit<Transaction, 'id'>) => Promise<void>;
    updateTransaction: (id: number, transaction: Partial<Transaction>) => Promise<void>;
    deleteTransaction: (id: number) => Promise<void>;
}

export const useTransactionStore = create<TransactionState>((set, get) => ({
    transactions: [],
    summary: null,
    isLoading: false,
    error: null,

    fetchTransactions: async (page = 1, limit = 20, search = '') => {
        set({ isLoading: true, error: null });
        try {
            const data = await apiClient.get<Transaction[]>(`/transactions?page=${page}&limit=${limit}&q=${search}`);
            set({ transactions: data, isLoading: false });
        } catch (error: any) {
            set({ error: error.message, isLoading: false });
        }
    },

    fetchSummary: async () => {
        try {
            const summary = await apiClient.get<DashboardSummary>('/transactions/summary');
            set({ summary });
        } catch (err: any) {
            console.error('Failed to fetch summary', err);
        }
    },

    addTransaction: async (transaction) => {
        set({ isLoading: true, error: null });
        try {
            await apiClient.post('/transactions', transaction);
            await get().fetchTransactions();
            await get().fetchSummary();
        } catch (err: any) {
            set({ error: err.message, isLoading: false });
            throw err;
        }
    },

    updateTransaction: async (id, transaction) => {
        set({ isLoading: true, error: null });
        try {
            await apiClient.put(`/transactions/${id}`, transaction);
            await get().fetchTransactions();
            await get().fetchSummary();
        } catch (err: any) {
            set({ error: err.message, isLoading: false });
            throw err;
        }
    },

    deleteTransaction: async (id) => {
        set({ isLoading: true, error: null });
        try {
            await apiClient.delete(`/transactions/${id}`);
            await get().fetchTransactions();
            await get().fetchSummary();
        } catch (err: any) {
            set({ error: err.message, isLoading: false });
            throw err;
        }
    },
}));
