export interface CategorySpending {
    category_id: number;
    category_name: string;
    total_amount: number;
    color: string;
}

export interface SpendingTrend {
    month: string;
    income: number;
    expense: number;
}
