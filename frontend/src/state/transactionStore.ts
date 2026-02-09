import { create } from 'zustand';
import type { Transaction, DashboardSummary } from '../domain/entities/Transaction';
import { apiClient } from '../infrastructure/apiClient';

export interface TransactionFilter {
    q?: string;
    category_id?: number | string;
    type?: string;
    start_date?: string;
    end_date?: string;
}

export interface PaginatedTransactions {
    transactions: Transaction[];
    total: number;
    total_amount: number;
}

interface TransactionState {
    transactions: Transaction[];
    summary: DashboardSummary | null;
    total: number;
    totalAmount: number;
    isLoading: boolean;
    error: string | null;
    fetchTransactions: (page?: number, limit?: number, filters?: TransactionFilter) => Promise<void>;
    fetchSummary: () => Promise<void>;
    addTransaction: (transaction: Omit<Transaction, 'id'>) => Promise<void>;
    updateTransaction: (id: number, transaction: Partial<Transaction>) => Promise<void>;
    deleteTransaction: (id: number) => Promise<void>;
}

export const useTransactionStore = create<TransactionState>((set, get) => ({
    transactions: [],
    summary: null,
    total: 0,
    totalAmount: 0,
    isLoading: false,
    error: null,

    fetchTransactions: async (page = 1, limit = 20, filters = {}) => {
        set({ isLoading: true, error: null });
        try {
            const params = new URLSearchParams({
                page: page.toString(),
                limit: limit.toString(),
                ...Object.fromEntries(
                    Object.entries(filters).filter(([_, v]) => v !== undefined && v !== '')
                )
            });
            const data = await apiClient.get<PaginatedTransactions>(`/transactions?${params.toString()}`);
            set({
                transactions: data.transactions || [],
                total: data.total || 0,
                totalAmount: data.total_amount || 0,
                isLoading: false
            });
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
