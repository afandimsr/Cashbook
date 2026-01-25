import { create } from 'zustand';
import type { CategorySpending } from '../domain/entities/Report';
import { apiClient } from '../infrastructure/apiClient';

interface ReportState {
    categorySpending: CategorySpending[];
    isLoading: boolean;
    error: string | null;
    fetchCategorySpending: (month: number, year: number) => Promise<void>;
}

export const useReportStore = create<ReportState>((set) => ({
    categorySpending: [],
    isLoading: false,
    error: null,

    fetchCategorySpending: async (month: number, year: number) => {
        set({ isLoading: true, error: null });
        try {
            const data = await apiClient.get<CategorySpending[]>(`/reports/spending?month=${month}&year=${year}`);
            set({ categorySpending: data || [], isLoading: false });
        } catch (error: any) {
            set({ error: error.message, isLoading: false });
        }
    },
}));
