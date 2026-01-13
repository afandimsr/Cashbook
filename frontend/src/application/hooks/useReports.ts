import { useEffect, useCallback } from 'react';
import { useReportStore } from '../../state/reportStore';

export const useReports = (month: number, year: number) => {
    const { categorySpending, isLoading, error, fetchCategorySpending } = useReportStore();

    useEffect(() => {
        fetchCategorySpending(month, year);
    }, [fetchCategorySpending, month, year]);

    const getChartData = useCallback(() => {
        return categorySpending.map(s => ({
            name: s.category_name,
            value: s.total_amount,
            color: s.color
        }));
    }, [categorySpending]);

    const getSortedSpending = useCallback(() => {
        return [...categorySpending].sort((a, b) => b.total_amount - a.total_amount);
    }, [categorySpending]);

    return {
        categorySpending,
        chartData: getChartData(),
        sortedSpending: getSortedSpending(),
        isLoading,
        error,
        refreshReports: () => fetchCategorySpending(month, year)
    };
};
