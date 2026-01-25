import { create } from 'zustand';
import { apiClient } from '../infrastructure/apiClient';
import type { Category } from '../domain/entities/Category';

interface CategoryState {
    categories: Category[];
    isLoading: boolean;
    error: string | null;
    fetchCategories: () => Promise<void>;
    addCategory: (category: Omit<Category, 'id'>) => Promise<void>;
    updateCategory: (id: number, category: Partial<Category>) => Promise<void>;
    deleteCategory: (id: number) => Promise<void>;
}

export const useCategoryStore = create<CategoryState>((set, get) => ({
    categories: [],
    isLoading: false,
    error: null,

    fetchCategories: async () => {
        set({ isLoading: true, error: null });
        try {
            const categories = await apiClient.get<Category[]>('/categories');
            set({ categories, isLoading: false });
        } catch (err: any) {
            set({ error: err.message, isLoading: false });
        }
    },

    addCategory: async (category) => {
        set({ isLoading: true, error: null });
        try {
            await apiClient.post('/categories', category);
            await get().fetchCategories();
        } catch (err: any) {
            set({ error: err.message, isLoading: false });
            throw err;
        }
    },

    updateCategory: async (id, category) => {
        set({ isLoading: true, error: null });
        try {
            await apiClient.put(`/categories/${id}`, category);
            await get().fetchCategories();
        } catch (err: any) {
            set({ error: err.message, isLoading: false });
            throw err;
        }
    },

    deleteCategory: async (id) => {
        set({ isLoading: true, error: null });
        try {
            await apiClient.delete(`/categories/${id}`);
            await get().fetchCategories();
        } catch (err: any) {
            set({ error: err.message, isLoading: false });
            throw err;
        }
    },
}));
