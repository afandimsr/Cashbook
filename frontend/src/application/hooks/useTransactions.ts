import { useEffect, useCallback } from 'react';
import { useSearchParams } from 'react-router-dom';
import { useTransactionStore, type TransactionFilter } from '../../state/transactionStore';
import type { Transaction } from '../../domain/entities/Transaction';

export const useTransactions = () => {
    const {
        transactions,
        isLoading,
        error,
        fetchTransactions,
        addTransaction,
        deleteTransaction
    } = useTransactionStore();

    const [searchParams, setSearchParams] = useSearchParams();

    const filters: TransactionFilter = {
        q: searchParams.get('q') || '',
        category_id: searchParams.get('category_id') || '',
        type: searchParams.get('type') || '',
        start_date: searchParams.get('start_date') || '',
        end_date: searchParams.get('end_date') || '',
    };

    const refreshTransactions = useCallback((page = 1, limit = 20) => {
        return fetchTransactions(page, limit, filters);
    }, [fetchTransactions, searchParams]); // Depends on searchParams to trigger on change

    useEffect(() => {
        refreshTransactions();
    }, [refreshTransactions]);

    const setFilters = (newFilters: Partial<TransactionFilter>) => {
        const nextParams = new URLSearchParams(searchParams);
        Object.entries(newFilters).forEach(([key, value]) => {
            if (value) {
                nextParams.set(key, value.toString());
            } else {
                nextParams.delete(key);
            }
        });
        setSearchParams(nextParams);
    };

    const clearFilters = () => {
        setSearchParams({});
    };

    const handleAddTransaction = async (data: Omit<Transaction, 'id'>) => {
        await addTransaction(data);
    };

    const handleDeleteTransaction = async (id: number) => {
        await deleteTransaction(id);
    };

    return {
        transactions,
        isLoading,
        error,
        refreshTransactions,
        handleAddTransaction,
        handleDeleteTransaction,
        filters,
        setFilters,
        clearFilters
    };
};
