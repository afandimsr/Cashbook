import { create } from 'zustand';
import { apiClient } from '../infrastructure/apiClient';
import type { SplitBill, SplitSummary, CreateSplitRequest } from '../domain/entities/SharedExpense';

interface SplitState {
  splits: SplitBill[];
  summary: SplitSummary | null;
  isLoading: boolean;
  error: string | null;
  
  fetchSplits: () => Promise<void>;
  fetchSummary: () => Promise<void>;
  fetchSplitByID: (id: string) => Promise<SplitBill>;
  createSplit: (request: CreateSplitRequest) => Promise<void>;
  settleParticipant: (billID: string, participantID: string) => Promise<void>;
}

export const useSplitStore = create<SplitState>((set, get) => ({
  splits: [],
  summary: null,
  isLoading: false,
  error: null,

  fetchSplits: async () => {
    set({ isLoading: true, error: null });
    try {
      const splits = await apiClient.get<SplitBill[]>('/splits');
      set({ splits: splits || [], isLoading: false });
    } catch (error: any) {
      set({ error: error.message, isLoading: false });
    }
  },

  fetchSummary: async () => {
    try {
      const summary = await apiClient.get<SplitSummary>('/splits/summary');
      set({ summary });
    } catch (err: any) {
      console.error('Failed to fetch split summary', err);
    }
  },

  fetchSplitByID: async (id: string) => {
    return apiClient.get<SplitBill>(`/splits/${id}`);
  },

  createSplit: async (request: CreateSplitRequest) => {
    set({ isLoading: true, error: null });
    try {
      await apiClient.post('/splits', request);
      await get().fetchSplits();
      await get().fetchSummary();
    } catch (err: any) {
      set({ error: err.message, isLoading: false });
      throw err;
    }
  },

  settleParticipant: async (billID: string, participantID: string) => {
    set({ isLoading: true, error: null });
    try {
      await apiClient.post(`/splits/${billID}/settle/${participantID}`, {});
      await get().fetchSplits();
      await get().fetchSummary();
    } catch (err: any) {
      set({ error: err.message, isLoading: false });
      throw err;
    }
  },
}));
