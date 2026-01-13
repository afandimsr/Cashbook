import { create } from 'zustand';
import type { Budget } from '../domain/entities/Budget';
import { apiClient } from '../infrastructure/apiClient';

interface BudgetState {
    budgets: Budget[];
    isLoading: boolean;
    error: string | null;
    fetchBudgets: (month: number, year: number) => Promise<void>;
    updateBudget: (budget: Partial<Budget>) => Promise<void>;
}

export const useBudgetStore = create<BudgetState>((set) => ({
    budgets: [],
    isLoading: false,
    error: null,

    fetchBudgets: async (month: number, year: number) => {
        set({ isLoading: true, error: null });
        try {
            const data = await apiClient.get<Budget[]>(`/budgets?month=${month}&year=${year}`);
            set({ budgets: data || [], isLoading: false });
        } catch (error: any) {
            set({ error: error.message, isLoading: false });
        }
    },

    updateBudget: async (budget: Partial<Budget>) => {
        set({ isLoading: true, error: null });
        try {
            await apiClient.post('/budgets', budget);
            // Refresh budgets for the month/year of the updated budget
            if (budget.month && budget.year) {
                const data = await apiClient.get<Budget[]>(`/budgets?month=${budget.month}&year=${budget.year}`);
                set({ budgets: data || [], isLoading: false });
            } else {
                set({ isLoading: false });
            }
        } catch (error: any) {
            set({ error: error.message, isLoading: false });
            throw error;
        }
    },
}));
