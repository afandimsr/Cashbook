import { useEffect, useCallback } from 'react';
import { useBudgetStore } from '../../state/budgetStore';
import { useCategoryStore } from '../../state/categoryStore';
import { useReportStore } from '../../state/reportStore';

export const useBudgets = (month: number, year: number) => {
    const { budgets, isLoading: isBudgetLoading, error: budgetError, fetchBudgets, updateBudget } = useBudgetStore();
    const { categories, fetchCategories } = useCategoryStore();
    const { categorySpending, fetchCategorySpending } = useReportStore();

    const refreshData = useCallback(() => {
        fetchBudgets(month, year);
        fetchCategories();
        fetchCategorySpending(month, year);
    }, [fetchBudgets, fetchCategories, fetchCategorySpending, month, year]);

    useEffect(() => {
        refreshData();
    }, [refreshData]);

    const handleUpdateBudget = async (data: { category_id: number; amount: number; month: number; year: number }) => {
        await updateBudget(data);
    };

    const getBudgetAnalysis = useCallback(() => {
        return categories
            ?.filter(c => c.type === 'expense')
            .map(category => {
                const budget = budgets.find(b => b.category_id === category.id);
                const spending = categorySpending.find(s => s.category_id === category.id)?.total_amount || 0;
                const limit = budget?.amount || 0;
                const percentage = limit > 0 ? (spending / limit) * 100 : 0;
                const isOver = spending > limit && limit > 0;

                return {
                    category,
                    budget,
                    spending,
                    limit,
                    percentage,
                    isOver
                };
            });
    }, [categories, budgets, categorySpending]);

    return {
        budgets,
        categories: categories?.filter(c => c.type === 'expense'),
        budgetAnalysis: getBudgetAnalysis(),
        isLoading: isBudgetLoading,
        error: budgetError,
        handleUpdateBudget,
        refreshData
    };
};
