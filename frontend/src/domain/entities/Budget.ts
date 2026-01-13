export interface Budget {
    id: number;
    user_id: number;
    category_id: number;
    amount: number;
    month: number;
    year: number;
}

export interface BudgetStatus extends Budget {
    spent: number;
    remaining: number;
    percentage: number;
}
