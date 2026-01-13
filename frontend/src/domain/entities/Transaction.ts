export interface Transaction {
    id?: number;
    user_id?: number;
    category_id: number;
    amount: number;
    note: string;
    date: string; // ISO string
    type: 'income' | 'expense';
}

export interface DashboardSummary {
    total_income: number;
    total_expense: number;
    balance: number;
}
