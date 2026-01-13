export type Frequency = 'daily' | 'weekly' | 'monthly';

export interface RecurringTransaction {
    id: number;
    user_id: number;
    category_id: number;
    amount: number;
    type: 'income' | 'expense';
    note: string;
    frequency: Frequency;
    start_date: string;
    last_processed?: string;
}
